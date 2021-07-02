package session

import (
	"fmt"
	"leapp_daemon/infrastructure/http/http_error"
	"sync"
)

var alibabaRamRoleChainedSessionsFacadeSingleton *AlibabaRamRoleChainedSessionsFacade
var alibabaRamRoleChainedSessionsFacadeLock sync.Mutex
var alibabaRamRoleChainedSessionsLock sync.Mutex

type AlibabaRamRoleChainedSessionsObserver interface {
	UpdateAlibabaRamRoleChainedSessions(oldAlibabaRamRoleChainedSessions []AlibabaRamRoleChainedSession, newAlibabaRamRoleChainedSessions []AlibabaRamRoleChainedSession) error
}

type AlibabaRamRoleChainedSessionsFacade struct {
	alibabaRamRoleChainedSessions []AlibabaRamRoleChainedSession
	observers              []AlibabaRamRoleChainedSessionsObserver
}

func GetAlibabaRamRoleChainedSessionsFacade() *AlibabaRamRoleChainedSessionsFacade {
	alibabaRamRoleChainedSessionsFacadeLock.Lock()
	defer alibabaRamRoleChainedSessionsFacadeLock.Unlock()

	if alibabaRamRoleChainedSessionsFacadeSingleton == nil {
		alibabaRamRoleChainedSessionsFacadeSingleton = &AlibabaRamRoleChainedSessionsFacade{
			alibabaRamRoleChainedSessions: make([]AlibabaRamRoleChainedSession, 0),
		}
	}

	return alibabaRamRoleChainedSessionsFacadeSingleton
}

func (fac *AlibabaRamRoleChainedSessionsFacade) Subscribe(observer AlibabaRamRoleChainedSessionsObserver) {
	fac.observers = append(fac.observers, observer)
}

func (fac *AlibabaRamRoleChainedSessionsFacade) GetSessions() []AlibabaRamRoleChainedSession {
	return fac.alibabaRamRoleChainedSessions
}

func (fac *AlibabaRamRoleChainedSessionsFacade) SetSessions(alibabaRamRoleChainedSessions []AlibabaRamRoleChainedSession) error {
	oldAlibabaRamRoleChainedSessions := fac.GetSessions()
	fac.alibabaRamRoleChainedSessions = alibabaRamRoleChainedSessions

	for _, observer := range fac.observers {
		err := observer.UpdateAlibabaRamRoleChainedSessions(oldAlibabaRamRoleChainedSessions, alibabaRamRoleChainedSessions)
		if err != nil {
			return err
		}
	}
	return nil
}

func (fac *AlibabaRamRoleChainedSessionsFacade) AddSession(alibabaRamRoleChainedSession AlibabaRamRoleChainedSession) error {
	alibabaRamRoleChainedSessionsLock.Lock()
	defer alibabaRamRoleChainedSessionsLock.Unlock()

	oldAlibabaRamRoleChainedSessions := fac.GetSessions()
	newAlibabaRamRoleChainedSessions := make([]AlibabaRamRoleChainedSession, 0)

	for i := range oldAlibabaRamRoleChainedSessions {
		newAlibabaRamRoleChainedSession := oldAlibabaRamRoleChainedSessions[i]
		newAlibabaRamRoleChainedSessionAccount := *oldAlibabaRamRoleChainedSessions[i].Account
		newAlibabaRamRoleChainedSession.Account = &newAlibabaRamRoleChainedSessionAccount
		newAlibabaRamRoleChainedSessions = append(newAlibabaRamRoleChainedSessions, newAlibabaRamRoleChainedSession)
	}

	for _, sess := range newAlibabaRamRoleChainedSessions {
		if alibabaRamRoleChainedSession.Id == sess.Id {
			return http_error.NewConflictError(fmt.Errorf("a AlibabaRamRoleChainedSession with id " + alibabaRamRoleChainedSession.Id +
				" is already present"))
		}

		/*if alibabaRamRoleChainedSession.Alias == sess.Alias {
			return http_error.NewUnprocessableEntityError(fmt.Errorf("a session with the same alias " +
				"is already present"))
		}*/
	}

	newAlibabaRamRoleChainedSessions = append(newAlibabaRamRoleChainedSessions, alibabaRamRoleChainedSession)

	err := fac.updateState(newAlibabaRamRoleChainedSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *AlibabaRamRoleChainedSessionsFacade) RemoveSession(id string) error {
	alibabaRamRoleChainedSessionsLock.Lock()
	defer alibabaRamRoleChainedSessionsLock.Unlock()

	oldAlibabaRamRoleChainedSessions := fac.GetSessions()
	newAlibabaRamRoleChainedSessions := make([]AlibabaRamRoleChainedSession, 0)

	for i := range oldAlibabaRamRoleChainedSessions {
		newAlibabaRamRoleChainedSession := oldAlibabaRamRoleChainedSessions[i]
		newAlibabaRamRoleChainedSessionAccount := *oldAlibabaRamRoleChainedSessions[i].Account
		newAlibabaRamRoleChainedSession.Account = &newAlibabaRamRoleChainedSessionAccount
		newAlibabaRamRoleChainedSessions = append(newAlibabaRamRoleChainedSessions, newAlibabaRamRoleChainedSession)
	}

	for i, sess := range newAlibabaRamRoleChainedSessions {
		if sess.Id == id {
			newAlibabaRamRoleChainedSessions = append(newAlibabaRamRoleChainedSessions[:i], newAlibabaRamRoleChainedSessions[i+1:]...)
			break
		}
	}

	if len(fac.GetSessions()) == len(newAlibabaRamRoleChainedSessions) {
		return http_error.NewNotFoundError(fmt.Errorf("trusted Alibaba session with id %s not found", id))
	}

	err := fac.updateState(newAlibabaRamRoleChainedSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *AlibabaRamRoleChainedSessionsFacade) GetSessionById(id string) (*AlibabaRamRoleChainedSession, error) {
	for _, alibabaRamRoleChainedSession := range fac.GetSessions() {
		if alibabaRamRoleChainedSession.Id == id {
			return &alibabaRamRoleChainedSession, nil
		}
	}
	return nil, http_error.NewNotFoundError(fmt.Errorf("trusted Alibaba session with id %s not found", id))
}

func (fac *AlibabaRamRoleChainedSessionsFacade) SetSessionById(newSession AlibabaRamRoleChainedSession) {
	allSessions := fac.GetSessions()
	for i, alibabaRamRoleChainedSession := range allSessions {
		if alibabaRamRoleChainedSession.Id == newSession.Id {
			allSessions[i] = newSession
		}
	}
	fac.SetSessions(allSessions)
}

func (fac *AlibabaRamRoleChainedSessionsFacade) SetSessionStatusToPending(id string) error {

	alibabaRamRoleChainedSessionsLock.Lock()
	defer alibabaRamRoleChainedSessionsLock.Unlock()

	alibabaRamRoleChainedSession, err := fac.GetSessionById(id)
	if err != nil {
		return err
	}

	if !(alibabaRamRoleChainedSession.Status == NotActive) {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("trusted Alibaba session with id " + id + "cannot be started because it's in pending or active state"))
	}

	oldAlibabaRamRoleChainedSessions := fac.GetSessions()
	newAlibabaRamRoleChainedSessions := make([]AlibabaRamRoleChainedSession, 0)

	for i := range oldAlibabaRamRoleChainedSessions {
		newAlibabaRamRoleChainedSession := oldAlibabaRamRoleChainedSessions[i]
		newAlibabaRamRoleChainedSessionAccount := *oldAlibabaRamRoleChainedSessions[i].Account
		newAlibabaRamRoleChainedSession.Account = &newAlibabaRamRoleChainedSessionAccount
		newAlibabaRamRoleChainedSessions = append(newAlibabaRamRoleChainedSessions, newAlibabaRamRoleChainedSession)
	}

	for i, session := range newAlibabaRamRoleChainedSessions {
		if session.Id == id {
			newAlibabaRamRoleChainedSessions[i].Status = Pending
		}
	}

	err = fac.updateState(newAlibabaRamRoleChainedSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *AlibabaRamRoleChainedSessionsFacade) SetSessionStatusToActive(id string) error {
	alibabaRamRoleChainedSessionsLock.Lock()
	defer alibabaRamRoleChainedSessionsLock.Unlock()

	alibabaRamRoleChainedSession, err := fac.GetSessionById(id)
	if err != nil {
		return err
	}

	if !(alibabaRamRoleChainedSession.Status == Pending) {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("trusted Alibaba session with id " + id + "cannot be started because it's not in pending state"))
	}

	oldAlibabaRamRoleChainedSessions := fac.GetSessions()
	newAlibabaRamRoleChainedSessions := make([]AlibabaRamRoleChainedSession, 0)

	for i := range oldAlibabaRamRoleChainedSessions {
		newAlibabaRamRoleChainedSession := oldAlibabaRamRoleChainedSessions[i]
		newAlibabaRamRoleChainedSessionAccount := *oldAlibabaRamRoleChainedSessions[i].Account
		newAlibabaRamRoleChainedSession.Account = &newAlibabaRamRoleChainedSessionAccount
		newAlibabaRamRoleChainedSessions = append(newAlibabaRamRoleChainedSessions, newAlibabaRamRoleChainedSession)
	}

	for i, session := range newAlibabaRamRoleChainedSessions {
		if session.Id == id {
			newAlibabaRamRoleChainedSessions[i].Status = Active
		}
	}

	err = fac.updateState(newAlibabaRamRoleChainedSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *AlibabaRamRoleChainedSessionsFacade) SetSessionStatusToInactive(id string) error {
	alibabaRamRoleChainedSessionsLock.Lock()
	defer alibabaRamRoleChainedSessionsLock.Unlock()

	alibabaRamRoleChainedSession, err := fac.GetSessionById(id)
	if err != nil {
		return err
	}
	if alibabaRamRoleChainedSession.Status != Active {
		fmt.Println(alibabaRamRoleChainedSession.Status)
		return http_error.NewUnprocessableEntityError(fmt.Errorf("trusted Alibaba session with id " + id + "cannot be stopped because it's not in active state"))
	}

	oldAlibabaRamRoleChainedSessions := fac.GetSessions()
	newAlibabaRamRoleChainedSessions := make([]AlibabaRamRoleChainedSession, 0)

	for i := range oldAlibabaRamRoleChainedSessions {
		newAlibabaRamRoleChainedSession := oldAlibabaRamRoleChainedSessions[i]
		newAlibabaRamRoleChainedSessionAccount := *oldAlibabaRamRoleChainedSessions[i].Account
		newAlibabaRamRoleChainedSession.Account = &newAlibabaRamRoleChainedSessionAccount
		newAlibabaRamRoleChainedSessions = append(newAlibabaRamRoleChainedSessions, newAlibabaRamRoleChainedSession)
	}

	for i, session := range newAlibabaRamRoleChainedSessions {
		if session.Id == id {
			newAlibabaRamRoleChainedSessions[i].Status = NotActive
		}
	}

	err = fac.updateState(newAlibabaRamRoleChainedSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *AlibabaRamRoleChainedSessionsFacade) updateState(newState []AlibabaRamRoleChainedSession) error {
	oldAlibabaRamRoleChainedSessions := fac.GetSessions()
	fac.alibabaRamRoleChainedSessions = newState

	for _, observer := range fac.observers {
		err := observer.UpdateAlibabaRamRoleChainedSessions(oldAlibabaRamRoleChainedSessions, newState)
		if err != nil {
			return err
		}
	}

	return nil
}

/*
func CreateAlibabaRamRoleChainedSession(sessionContainer Container, name string, accountNumber string, roleName string, roleArn string, idpArn string,
	region string, ssoUrl string, profile string) error {

	sessions, err := sessionContainer.GetAlibabaRamRoleChainedSessions()
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

	role := AlibabaRamRoleChainedRole{
		Name: roleName,
		Arn:  roleArn,
	}

	alibabaRamRoleChainedAccount := AlibabaRamRoleChainedAccount{
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

	session := AlibabaRamRoleChainedSession{
		Id:        uuidString,
		Status:    NotActive,
		StartTime: "",
		Account:   &alibabaRamRoleChainedAccount,
		Profile:   namedProfileId,
	}

	err = sessionContainer.SetAlibabaRamRoleChainedSessions(append(sessions, &session))
	if err != nil { return err }

	return nil
}

func GetAlibabaRamRoleChainedSession(sessionContainer Container, id string) (*AlibabaRamRoleChainedSession, error) {
	sessions, err := sessionContainer.GetAlibabaRamRoleChainedSessions()
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

func ListAlibabaRamRoleChainedSession(sessionContainer Container, query string) ([]*AlibabaRamRoleChainedSession, error) {
	sessions, err := sessionContainer.GetAlibabaRamRoleChainedSessions()
	if err != nil {
		return nil, err
	}

	filteredList := make([]*AlibabaRamRoleChainedSession, 0)

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

func UpdateAlibabaRamRoleChainedSession(sessionContainer Container, id string, name string, accountNumber string, roleName string, roleArn string, idpArn string,
	region string, ssoUrl string, profile string) error {

	sessions, err := sessionContainer.GetAlibabaRamRoleChainedSessions()
	if err != nil { return err }

	found := false
	for index := range sessions {
		if sessions[index].Id == id {
			namedProfileId, err := named_profile.EditNamedProfile(sessionContainer, sessions[index].Profile, profile)
			if err != nil { return err }

			sessions[index].Profile = namedProfileId
			sessions[index].Account = &AlibabaRamRoleChainedAccount{
				AccountNumber: accountNumber,
				Name:          name,
				Region:        region,
				IdpArn: 	   idpArn,
				SsoUrl:        ssoUrl,
			}

			sessions[index].Account.Role = &AlibabaRamRoleChainedRole{
				Name: roleName,
				Arn:  roleArn,
			}

			found = true
		}
	}

	if found == false {
		err = http_error2.NewNotFoundError(fmt.Errorf("trusted AWS session with id " + id + " not found"))
		return err
	}

	err = sessionContainer.SetAlibabaRamRoleChainedSessions(sessions)
	if err != nil { return err }

	return nil
}

func DeleteAlibabaRamRoleChainedSession(sessionContainer Container, id string) error {
	sessions, err := sessionContainer.GetAlibabaRamRoleChainedSessions()
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
		err = http_error2.NewNotFoundError(fmt.Errorf("trusted AWS session with id " + id + " not found"))
		return err
	}

	err = sessionContainer.SetAlibabaRamRoleChainedSessions(sessions)
	if err != nil {
		return err
	}

	return nil
}

func StartAlibabaRamRoleChainedSession(sessionContainer Container, id string) error {
	sess, err := GetAlibabaRamRoleChainedSession(sessionContainer, id)
	if err != nil {
		return err
	}

	println("Rotating session with id", sess.Id)
	err = sess.Rotate(nil)
	if err != nil { return err }

	return nil
}

func StopAlibabaRamRoleChainedSession(sessionContainer Container, id string) error {
	sess, err := GetAlibabaRamRoleChainedSession(sessionContainer, id)
	if err != nil {
		return err
	}

	sess.Status = NotActive
	return nil
}
*/
