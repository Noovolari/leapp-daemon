package session

import (
	"fmt"
	"leapp_daemon/infrastructure/http/http_error"
	"sync"
)

var federatedAlibabaSessionsFacadeSingleton *AlibabaRamRoleFederatedSessionsFacade
var federatedAlibabaSessionsFacadeLock sync.Mutex
var federatedAlibabaSessionsLock sync.Mutex

type AlibabaRamRoleFederatedSessionsObserver interface {
	UpdateAlibabaRamRoleFederatedSessions(oldAlibabaRamRoleFederatedSessions []AlibabaRamRoleFederatedSession, newAlibabaRamRoleFederatedSessions []AlibabaRamRoleFederatedSession) error
}

type AlibabaRamRoleFederatedSessionsFacade struct {
	federatedAlibabaSessions []AlibabaRamRoleFederatedSession
	observers                []AlibabaRamRoleFederatedSessionsObserver
}

func GetAlibabaRamRoleFederatedSessionsFacade() *AlibabaRamRoleFederatedSessionsFacade {
	federatedAlibabaSessionsFacadeLock.Lock()
	defer federatedAlibabaSessionsFacadeLock.Unlock()

	if federatedAlibabaSessionsFacadeSingleton == nil {
		federatedAlibabaSessionsFacadeSingleton = &AlibabaRamRoleFederatedSessionsFacade{
			federatedAlibabaSessions: make([]AlibabaRamRoleFederatedSession, 0),
		}
	}

	return federatedAlibabaSessionsFacadeSingleton
}

func (fac *AlibabaRamRoleFederatedSessionsFacade) Subscribe(observer AlibabaRamRoleFederatedSessionsObserver) {
	fac.observers = append(fac.observers, observer)
}

func (fac *AlibabaRamRoleFederatedSessionsFacade) GetSessions() []AlibabaRamRoleFederatedSession {
	return fac.federatedAlibabaSessions
}

func (fac *AlibabaRamRoleFederatedSessionsFacade) SetSessions(newAlibabaRamRoleFederatedSessions []AlibabaRamRoleFederatedSession) error {
	fac.federatedAlibabaSessions = newAlibabaRamRoleFederatedSessions

	err := fac.updateState(newAlibabaRamRoleFederatedSessions)
	if err != nil {
		return err
	}
	return nil
}

func (fac *AlibabaRamRoleFederatedSessionsFacade) UpdateSession(newSession AlibabaRamRoleFederatedSession) error {
	allSessions := fac.GetSessions()
	for i, federatedAlibabaSession := range allSessions {
		if federatedAlibabaSession.Id == newSession.Id {
			allSessions[i] = newSession
		}
	}
	err := fac.SetSessions(allSessions)
	return err
}

func (fac *AlibabaRamRoleFederatedSessionsFacade) AddSession(federatedAlibabaSession AlibabaRamRoleFederatedSession) error {
	federatedAlibabaSessionsLock.Lock()
	defer federatedAlibabaSessionsLock.Unlock()

	oldAlibabaRamRoleFederatedSessions := fac.GetSessions()
	newAlibabaRamRoleFederatedSessions := make([]AlibabaRamRoleFederatedSession, 0)

	for i := range oldAlibabaRamRoleFederatedSessions {
		newAlibabaRamRoleFederatedSession := oldAlibabaRamRoleFederatedSessions[i]
		newAlibabaRamRoleFederatedSessionAccount := *oldAlibabaRamRoleFederatedSessions[i].Account
		newAlibabaRamRoleFederatedSession.Account = &newAlibabaRamRoleFederatedSessionAccount
		newAlibabaRamRoleFederatedSessions = append(newAlibabaRamRoleFederatedSessions, newAlibabaRamRoleFederatedSession)
	}

	for _, sess := range newAlibabaRamRoleFederatedSessions {
		if federatedAlibabaSession.Id == sess.Id {
			return http_error.NewConflictError(fmt.Errorf("a AlibabaRamRoleFederatedSession with id " + federatedAlibabaSession.Id +
				" is already present"))
		}

		/*if federatedAlibabaSession.Alias == sess.Alias {
			return http_error.NewUnprocessableEntityError(fmt.Errorf("a session with the same alias " +
				"is already present"))
		}*/
	}

	newAlibabaRamRoleFederatedSessions = append(newAlibabaRamRoleFederatedSessions, federatedAlibabaSession)

	err := fac.updateState(newAlibabaRamRoleFederatedSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *AlibabaRamRoleFederatedSessionsFacade) RemoveSession(id string) error {
	federatedAlibabaSessionsLock.Lock()
	defer federatedAlibabaSessionsLock.Unlock()

	oldAlibabaRamRoleFederatedSessions := fac.GetSessions()
	newAlibabaRamRoleFederatedSessions := make([]AlibabaRamRoleFederatedSession, 0)

	for i := range oldAlibabaRamRoleFederatedSessions {
		newAlibabaRamRoleFederatedSession := oldAlibabaRamRoleFederatedSessions[i]
		newAlibabaRamRoleFederatedSessionAccount := *oldAlibabaRamRoleFederatedSessions[i].Account
		newAlibabaRamRoleFederatedSession.Account = &newAlibabaRamRoleFederatedSessionAccount
		newAlibabaRamRoleFederatedSessions = append(newAlibabaRamRoleFederatedSessions, newAlibabaRamRoleFederatedSession)
	}

	for i, sess := range newAlibabaRamRoleFederatedSessions {
		if sess.Id == id {
			newAlibabaRamRoleFederatedSessions = append(newAlibabaRamRoleFederatedSessions[:i], newAlibabaRamRoleFederatedSessions[i+1:]...)
			break
		}
	}

	if len(fac.GetSessions()) == len(newAlibabaRamRoleFederatedSessions) {
		return http_error.NewNotFoundError(fmt.Errorf("federated Alibaba session with id %s not found", id))
	}

	err := fac.updateState(newAlibabaRamRoleFederatedSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *AlibabaRamRoleFederatedSessionsFacade) GetSessionById(id string) (*AlibabaRamRoleFederatedSession, error) {
	for _, federatedAlibabaSession := range fac.GetSessions() {
		if federatedAlibabaSession.Id == id {
			return &federatedAlibabaSession, nil
		}
	}
	return nil, http_error.NewNotFoundError(fmt.Errorf("federated Alibaba session with id %s not found", id))
}

func (fac *AlibabaRamRoleFederatedSessionsFacade) SetSessionById(newSession AlibabaRamRoleFederatedSession) {
	allSessions := fac.GetSessions()
	for i, federatedAlibabaSession := range allSessions {
		if federatedAlibabaSession.Id == newSession.Id {
			allSessions[i] = newSession
		}
	}
	fac.SetSessions(allSessions)
}

func (fac *AlibabaRamRoleFederatedSessionsFacade) SetStatusToPending(id string) error {
	federatedAlibabaSessionsLock.Lock()
	defer federatedAlibabaSessionsLock.Unlock()

	federatedAlibabaSession, err := fac.GetSessionById(id)
	if err != nil {
		return err
	}

	if !(federatedAlibabaSession.Status == NotActive) {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("federated Alibaba session with id " + id + "cannot be started because it's in pending or active state"))
	}

	oldAlibabaRamRoleFederatedSessions := fac.GetSessions()
	newAlibabaRamRoleFederatedSessions := make([]AlibabaRamRoleFederatedSession, 0)

	for i := range oldAlibabaRamRoleFederatedSessions {
		newAlibabaRamRoleFederatedSession := oldAlibabaRamRoleFederatedSessions[i]
		newAlibabaRamRoleFederatedSessionAccount := *oldAlibabaRamRoleFederatedSessions[i].Account
		newAlibabaRamRoleFederatedSession.Account = &newAlibabaRamRoleFederatedSessionAccount
		newAlibabaRamRoleFederatedSessions = append(newAlibabaRamRoleFederatedSessions, newAlibabaRamRoleFederatedSession)
	}

	for i, session := range newAlibabaRamRoleFederatedSessions {
		if session.Id == id {
			newAlibabaRamRoleFederatedSessions[i].Status = Pending
		}
	}

	err = fac.updateState(newAlibabaRamRoleFederatedSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *AlibabaRamRoleFederatedSessionsFacade) SetStatusToActive(id string) error {
	federatedAlibabaSessionsLock.Lock()
	defer federatedAlibabaSessionsLock.Unlock()

	federatedAlibabaSession, err := fac.GetSessionById(id)
	if err != nil {
		return err
	}

	if !(federatedAlibabaSession.Status == Pending) {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("federated Alibaba session with id " + id + "cannot be started because it's not in pending state"))
	}

	oldAlibabaRamRoleFederatedSessions := fac.GetSessions()
	newAlibabaRamRoleFederatedSessions := make([]AlibabaRamRoleFederatedSession, 0)

	for i := range oldAlibabaRamRoleFederatedSessions {
		newAlibabaRamRoleFederatedSession := oldAlibabaRamRoleFederatedSessions[i]
		newAlibabaRamRoleFederatedSessionAccount := *oldAlibabaRamRoleFederatedSessions[i].Account
		newAlibabaRamRoleFederatedSession.Account = &newAlibabaRamRoleFederatedSessionAccount
		newAlibabaRamRoleFederatedSessions = append(newAlibabaRamRoleFederatedSessions, newAlibabaRamRoleFederatedSession)
	}

	for i, session := range newAlibabaRamRoleFederatedSessions {
		if session.Id == id {
			newAlibabaRamRoleFederatedSessions[i].Status = Active
		}
	}

	err = fac.updateState(newAlibabaRamRoleFederatedSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *AlibabaRamRoleFederatedSessionsFacade) SetStatusToInactive(id string) error {
	federatedAlibabaSessionsLock.Lock()
	defer federatedAlibabaSessionsLock.Unlock()

	federatedAlibabaSession, err := fac.GetSessionById(id)
	if err != nil {
		return err
	}
	if federatedAlibabaSession.Status != Active {
		fmt.Println(federatedAlibabaSession.Status)
		return http_error.NewUnprocessableEntityError(fmt.Errorf("federated Alibaba session with id " + id + "cannot be stopped because it's not in active state"))
	}

	oldAlibabaRamRoleFederatedSessions := fac.GetSessions()
	newAlibabaRamRoleFederatedSessions := make([]AlibabaRamRoleFederatedSession, 0)

	for i := range oldAlibabaRamRoleFederatedSessions {
		newAlibabaRamRoleFederatedSession := oldAlibabaRamRoleFederatedSessions[i]
		newAlibabaRamRoleFederatedSessionAccount := *oldAlibabaRamRoleFederatedSessions[i].Account
		newAlibabaRamRoleFederatedSession.Account = &newAlibabaRamRoleFederatedSessionAccount
		newAlibabaRamRoleFederatedSessions = append(newAlibabaRamRoleFederatedSessions, newAlibabaRamRoleFederatedSession)
	}

	for i, session := range newAlibabaRamRoleFederatedSessions {
		if session.Id == id {
			newAlibabaRamRoleFederatedSessions[i].Status = NotActive
		}
	}

	err = fac.updateState(newAlibabaRamRoleFederatedSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *AlibabaRamRoleFederatedSessionsFacade) updateState(newState []AlibabaRamRoleFederatedSession) error {
	oldAlibabaRamRoleFederatedSessions := fac.GetSessions()
	fac.federatedAlibabaSessions = newState

	for _, observer := range fac.observers {
		err := observer.UpdateAlibabaRamRoleFederatedSessions(oldAlibabaRamRoleFederatedSessions, newState)
		if err != nil {
			return err
		}
	}

	return nil
}

/*
func CreateAlibabaRamRoleFederatedSession(sessionContainer Container, name string, accountNumber string, roleName string, roleArn string, idpArn string,
	region string, ssoUrl string, profile string) error {

	sessions, err := sessionContainer.GetAlibabaRamRoleFederatedSessions()
	if err != nil {
		return err
	}

	for _, session := range sessions {
		account := session.Account
		if account.AccountNumber == accountNumber && account.Role.Name == roleName {
			err = http_error2.NewUnprocessableEntityError(fmt.Errorf("an account with the same account number and " +
				"role name is already present"))
			return err
		}
	}

	role := AlibabaRamRoleFederatedRole{
		Name: roleName,
		Arn:  roleArn,
	}

	federatedAlibabaAccount := AlibabaRamRoleFederatedAccount{
		AccountNumber: accountNumber,
		Name:          name,
		Role:          &role,
		IdpArn:        idpArn,
		Region:        region,
		SsoUrl:        ssoUrl,
	}

	uuidString := uuid.New().String()
	uuidString = strings.Replace(uuidString, "-", "", -1)

	namedProfileId, err := named_profile.CreateNamedProfile(sessionContainer, profile)
	if err != nil {
		return err
	}

	session := AlibabaRamRoleFederatedSession{
		Id:        uuidString,
		Status:    NotActive,
		StartTime: "",
		Account:   &federatedAlibabaAccount,
		Profile:   namedProfileId,
	}

	err = sessionContainer.SetAlibabaRamRoleFederatedSessions(append(sessions, &session))
	if err != nil { return err }

	return nil
}

func GetAlibabaRamRoleFederatedSession(sessionContainer Container, id string) (*AlibabaRamRoleFederatedSession, error) {
	sessions, err := sessionContainer.GetAlibabaRamRoleFederatedSessions()
	if err != nil {
		return nil, err
	}

	for index, _ := range sessions {
		if sessions[index].Id == id {
			return sessions[index], nil
		}
	}

	return nil, http_error2.NewNotFoundError(fmt.Errorf("No session found with id:" + id))
}

func ListAlibabaRamRoleFederatedSession(sessionContainer Container, query string) ([]*AlibabaRamRoleFederatedSession, error) {
	sessions, err := sessionContainer.GetAlibabaRamRoleFederatedSessions()
	if err != nil {
		return nil, err
	}

	filteredList := make([]*AlibabaRamRoleFederatedSession, 0)

	if query == "" {
		return append(filteredList, sessions...), nil
	} else {
		for _, session := range sessions {
			if  strings.Contains(session.Id, query) ||
				strings.Contains(session.Profile, query) ||
				strings.Contains(session.Account.Name, query) ||
				strings.Contains(session.Account.IdpArn, query) ||
				strings.Contains(session.Account.SsoUrl, query) ||
				strings.Contains(session.Account.Region, query) ||
				strings.Contains(session.Account.AccountNumber, query) ||
				strings.Contains(session.Account.Role.Name, query) ||
				strings.Contains(session.Account.Role.Arn, query) {

				filteredList = append(filteredList, session)
			}
		}

		return filteredList, nil
	}
}

func UpdateAlibabaRamRoleFederatedSession(sessionContainer Container, id string, name string, accountNumber string, roleName string, roleArn string, idpArn string,
	region string, ssoUrl string, profile string) error {

	sessions, err := sessionContainer.GetAlibabaRamRoleFederatedSessions()
	if err != nil { return err }

	found := false
	for index := range sessions {
		if sessions[index].Id == id {
			namedProfileId, err := named_profile.EditNamedProfile(sessionContainer, sessions[index].Profile, profile)
			if err != nil { return err }

			sessions[index].Profile = namedProfileId
			sessions[index].Account = &AlibabaRamRoleFederatedAccount{
				AccountNumber: accountNumber,
				Name:          name,
				Region:        region,
				IdpArn: 	   idpArn,
				SsoUrl:        ssoUrl,
			}

			sessions[index].Account.Role = &AlibabaRamRoleFederatedRole{
				Name: roleName,
				Arn:  roleArn,
			}

			found = true
		}
	}

	if found == false {
		err = http_error2.NewNotFoundError(fmt.Errorf("federated AWS session with id " + id + " not found"))
		return err
	}

	err = sessionContainer.SetAlibabaRamRoleFederatedSessions(sessions)
	if err != nil { return err }

	return nil
}

func DeleteAlibabaRamRoleFederatedSession(sessionContainer Container, id string) error {
	sessions, err := sessionContainer.GetAlibabaRamRoleFederatedSessions()
	if err != nil {
		return err
	}

	found := false
	for index := range sessions {
		if sessions[index].Id == id {
			sessions = append(sessions[:index], sessions[index+1:]...)
			found = true
			break
		}
	}

	if found == false {
		err = http_error2.NewNotFoundError(fmt.Errorf("federated AWS session with id " + id + " not found"))
		return err
	}

	err = sessionContainer.SetAlibabaRamRoleFederatedSessions(sessions)
	if err != nil {
		return err
	}

	return nil
}

func StartAlibabaRamRoleFederatedSession(sessionContainer Container, id string) error {
	sess, err := GetAlibabaRamRoleFederatedSession(sessionContainer, id)
	if err != nil {
		return err
	}

	println("Rotating session with id", sess.Id)
	err = sess.Rotate(nil)
	if err != nil { return err }

	return nil
}

func StopAlibabaRamRoleFederatedSession(sessionContainer Container, id string) error {
	sess, err := GetAlibabaRamRoleFederatedSession(sessionContainer, id)
	if err != nil {
		return err
	}

	sess.Status = NotActive
	return nil
}
*/
