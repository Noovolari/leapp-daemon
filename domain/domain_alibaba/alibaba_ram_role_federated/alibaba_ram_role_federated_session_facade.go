package alibaba_ram_role_federated

import (
	"fmt"
	"leapp_daemon/domain/domain_alibaba"
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
	alibabaRamRoleFederatedSessions []AlibabaRamRoleFederatedSession
	observers                       []AlibabaRamRoleFederatedSessionsObserver
}

func NewAlibabaRamRoleFederatedSessionsFacade() *AlibabaRamRoleFederatedSessionsFacade {
	return &AlibabaRamRoleFederatedSessionsFacade{
		alibabaRamRoleFederatedSessions: make([]AlibabaRamRoleFederatedSession, 0),
	}
}

func GetAlibabaRamRoleFederatedSessionsFacade() *AlibabaRamRoleFederatedSessionsFacade {
	federatedAlibabaSessionsFacadeLock.Lock()
	defer federatedAlibabaSessionsFacadeLock.Unlock()

	if federatedAlibabaSessionsFacadeSingleton == nil {
		federatedAlibabaSessionsFacadeSingleton = &AlibabaRamRoleFederatedSessionsFacade{
			alibabaRamRoleFederatedSessions: make([]AlibabaRamRoleFederatedSession, 0),
		}
	}

	return federatedAlibabaSessionsFacadeSingleton
}

func (fac *AlibabaRamRoleFederatedSessionsFacade) Subscribe(observer AlibabaRamRoleFederatedSessionsObserver) {
	fac.observers = append(fac.observers, observer)
}

func (fac *AlibabaRamRoleFederatedSessionsFacade) GetSessions() []AlibabaRamRoleFederatedSession {
	return fac.alibabaRamRoleFederatedSessions
}

func (fac *AlibabaRamRoleFederatedSessionsFacade) SetSessions(newSessions []AlibabaRamRoleFederatedSession) error {
	fac.alibabaRamRoleFederatedSessions = newSessions

	err := fac.updateState(newSessions)
	if err != nil {
		return err
	}
	return nil
}

func (facade *AlibabaRamRoleFederatedSessionsFacade) EditSession(sessionId string, sessionName string, roleName string, roleArn string, idpArn string,
	region string, ssoUrl string, namedProfileId string) error {
	federatedAlibabaSessionsLock.Lock()
	defer federatedAlibabaSessionsLock.Unlock()

	sessionToEdit, err := facade.GetSessionById(sessionId)
	if err != nil {
		return err
	}

	currentSessions := facade.GetSessions()
	for _, sess := range currentSessions {

		if sess.Id != sessionId && sess.Name == sessionName {
			return http_error.NewConflictError(fmt.Errorf("a session named %v is already present", sess.Name))
		}
	}

	sessionToEdit.Name = sessionName
	sessionToEdit.RoleName = roleName
	sessionToEdit.RoleArn = roleArn
	sessionToEdit.IdpArn = idpArn
	sessionToEdit.Region = region
	sessionToEdit.SsoUrl = ssoUrl
	sessionToEdit.NamedProfileId = namedProfileId
	return facade.replaceSession(sessionToEdit)
}

func (facade *AlibabaRamRoleFederatedSessionsFacade) AddSession(newSession AlibabaRamRoleFederatedSession) error {
	federatedAlibabaSessionsLock.Lock()
	defer federatedAlibabaSessionsLock.Unlock()

	currentSessions := facade.GetSessions()

	for _, sess := range currentSessions {
		if newSession.Id == sess.Id {
			return http_error.NewConflictError(fmt.Errorf("a session with id %v is already present", newSession.Id))
		}

		if newSession.Name == sess.Name {
			return http_error.NewConflictError(fmt.Errorf("a session named %v is already present", sess.Name))
		}
	}

	newSessions := append(currentSessions, newSession)

	facade.updateState(newSessions)
	return nil

}

func (fac *AlibabaRamRoleFederatedSessionsFacade) RemoveSession(id string) error {
	federatedAlibabaSessionsLock.Lock()
	defer federatedAlibabaSessionsLock.Unlock()

	currentSessions := fac.GetSessions()
	newSessions := make([]AlibabaRamRoleFederatedSession, 0)

	for _, session := range currentSessions {
		if session.Id != id {
			newSessions = append(newSessions, session)
		}
	}

	if len(fac.GetSessions()) == len(newSessions) {
		return http_error.NewNotFoundError(fmt.Errorf("session with id %s not found", id))
	}

	fac.updateState(newSessions)
	return nil
}

func (fac *AlibabaRamRoleFederatedSessionsFacade) GetSessionById(id string) (*AlibabaRamRoleFederatedSession, error) {
	for _, federatedAlibabaSession := range fac.GetSessions() {
		if federatedAlibabaSession.Id == id {
			return &federatedAlibabaSession, nil
		}
	}
	return nil, http_error.NewNotFoundError(fmt.Errorf("session with id %s not found", id))
}

func (fac *AlibabaRamRoleFederatedSessionsFacade) SetSessionById(newSession AlibabaRamRoleFederatedSession) {
	currentSessions := fac.GetSessions()
	for i, session := range currentSessions {
		if session.Id == newSession.Id {
			currentSessions[i] = newSession
		}
	}
	fac.SetSessions(currentSessions)
}

func (facade *AlibabaRamRoleFederatedSessionsFacade) StartingSession(sessionId string) error {
	return facade.setSessionStatus(sessionId, domain_alibaba.Pending)
}

func (facade *AlibabaRamRoleFederatedSessionsFacade) StartSession(sessionId string) error {
	return facade.setSessionStatus(sessionId, domain_alibaba.Active)
}

func (facade *AlibabaRamRoleFederatedSessionsFacade) StopSession(sessionId string) error {
	return facade.setSessionStatus(sessionId, domain_alibaba.NotActive)
}

/*func (fac *AlibabaRamRoleFederatedSessionsFacade) SetSessionStatusToPending(id string) error {
	federatedAlibabaSessionsLock.Lock()
	defer federatedAlibabaSessionsLock.Unlock()

	federatedAlibabaSession, err := fac.GetSessionById(id)
	if err != nil {
		return err
	}

	if !(federatedAlibabaSession.Status == domain_alibaba.NotActive) {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("session with id " + id + "cannot be started because it's in pending or active state"))
	}

	oldAlibabaRamRoleFederatedSessions := fac.GetSessions()
	newAlibabaRamRoleFederatedSessions := make([]AlibabaRamRoleFederatedSession, 0)

	for i := range oldAlibabaRamRoleFederatedSessions {
		newAlibabaRamRoleFederatedSession := oldAlibabaRamRoleFederatedSessions[i]
		newAlibabaRamRoleFederatedSessions = append(newAlibabaRamRoleFederatedSessions, newAlibabaRamRoleFederatedSession)
	}

	for i, session := range newAlibabaRamRoleFederatedSessions {
		if session.Id == id {
			newAlibabaRamRoleFederatedSessions[i].Status = domain_alibaba.Pending
		}
	}

	err = fac.updateState(newAlibabaRamRoleFederatedSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *AlibabaRamRoleFederatedSessionsFacade) SetSessionStatusToActive(id string) error {
	federatedAlibabaSessionsLock.Lock()
	defer federatedAlibabaSessionsLock.Unlock()

	federatedAlibabaSession, err := fac.GetSessionById(id)
	if err != nil {
		return err
	}

	if !(federatedAlibabaSession.Status == domain_alibaba.Pending) {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("session with id " + id + "cannot be started because it's not in pending state"))
	}

	oldAlibabaRamRoleFederatedSessions := fac.GetSessions()
	newAlibabaRamRoleFederatedSessions := make([]AlibabaRamRoleFederatedSession, 0)

	for i := range oldAlibabaRamRoleFederatedSessions {
		newAlibabaRamRoleFederatedSession := oldAlibabaRamRoleFederatedSessions[i]
		newAlibabaRamRoleFederatedSessions = append(newAlibabaRamRoleFederatedSessions, newAlibabaRamRoleFederatedSession)
	}

	for i, session := range newAlibabaRamRoleFederatedSessions {
		if session.Id == id {
			newAlibabaRamRoleFederatedSessions[i].Status = domain_alibaba.Active
		}
	}

	err = fac.updateState(newAlibabaRamRoleFederatedSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *AlibabaRamRoleFederatedSessionsFacade) SetSessionStatusToInactive(id string) error {
	federatedAlibabaSessionsLock.Lock()
	defer federatedAlibabaSessionsLock.Unlock()

	federatedAlibabaSession, err := fac.GetSessionById(id)
	if err != nil {
		return err
	}
	if federatedAlibabaSession.Status != domain_alibaba.Active {
		fmt.Println(federatedAlibabaSession.Status)
		return http_error.NewUnprocessableEntityError(fmt.Errorf("session with id " + id + "cannot be stopped because it's not in active state"))
	}

	oldAlibabaRamRoleFederatedSessions := fac.GetSessions()
	newAlibabaRamRoleFederatedSessions := make([]AlibabaRamRoleFederatedSession, 0)

	for i := range oldAlibabaRamRoleFederatedSessions {
		newAlibabaRamRoleFederatedSession := oldAlibabaRamRoleFederatedSessions[i]
		newAlibabaRamRoleFederatedSessions = append(newAlibabaRamRoleFederatedSessions, newAlibabaRamRoleFederatedSession)
	}

	for i, session := range newAlibabaRamRoleFederatedSessions {
		if session.Id == id {
			newAlibabaRamRoleFederatedSessions[i].Status = domain_alibaba.NotActive
		}
	}

	err = fac.updateState(newAlibabaRamRoleFederatedSessions)
	if err != nil {
		return err
	}

	return nil
}*/

func (facade *AlibabaRamRoleFederatedSessionsFacade) setSessionStatus(id string, status domain_alibaba.AlibabaSessionStatus) error {
	federatedAlibabaSessionsLock.Lock()
	defer federatedAlibabaSessionsLock.Unlock()

	sessionToUpdate, err := facade.GetSessionById(id)
	if err != nil {
		return err
	}

	sessionToUpdate.Status = status
	/*if startTime != "" {
		sessionToUpdate.StartTime = startTime
		sessionToUpdate.LastStopTime = ""
	}
	if lastStopTime != "" {
		sessionToUpdate.StartTime = ""
		sessionToUpdate.LastStopTime = lastStopTime
	}*/
	return facade.replaceSession(sessionToUpdate)
}

func (facade *AlibabaRamRoleFederatedSessionsFacade) replaceSession(newSession *AlibabaRamRoleFederatedSession) error {
	newSessions := make([]AlibabaRamRoleFederatedSession, 0)
	for _, session := range facade.GetSessions() {
		if session.Id == newSession.Id {
			newSessions = append(newSessions, *newSession)
		} else {
			newSessions = append(newSessions, session)
		}
	}

	facade.updateState(newSessions)
	return nil
}

func (fac *AlibabaRamRoleFederatedSessionsFacade) updateState(newState []AlibabaRamRoleFederatedSession) error {
	oldAlibabaRamRoleFederatedSessions := fac.GetSessions()
	fac.alibabaRamRoleFederatedSessions = newState

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
