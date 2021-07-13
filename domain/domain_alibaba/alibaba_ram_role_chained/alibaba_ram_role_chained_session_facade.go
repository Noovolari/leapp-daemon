package alibaba_ram_role_chained

import (
	"fmt"
	"leapp_daemon/domain/domain_alibaba"
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
	observers                     []AlibabaRamRoleChainedSessionsObserver
}

func NewAlibabaRamRoleChainedSessionsFacade() *AlibabaRamRoleChainedSessionsFacade {
	return &AlibabaRamRoleChainedSessionsFacade{
		alibabaRamRoleChainedSessions: make([]AlibabaRamRoleChainedSession, 0),
	}
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

func (fac *AlibabaRamRoleChainedSessionsFacade) SetSessions(newSessions []AlibabaRamRoleChainedSession) error {
	fac.alibabaRamRoleChainedSessions = newSessions

	err := fac.updateState(newSessions)
	if err != nil {
		return err
	}
	return nil
}

func (facade *AlibabaRamRoleChainedSessionsFacade) AddSession(newSession AlibabaRamRoleChainedSession) error {
	alibabaRamRoleChainedSessionsLock.Lock()
	defer alibabaRamRoleChainedSessionsLock.Unlock()

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

func (facade *AlibabaRamRoleChainedSessionsFacade) EditSession(sessionId string, sessionName string, roleName string, accountNumber string, roleArn string,
	region string, parentId string, parentType string, namedProfileId string) error {
	alibabaRamRoleChainedSessionsLock.Lock()
	defer alibabaRamRoleChainedSessionsLock.Unlock()

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
	sessionToEdit.AccountNumber = accountNumber
	sessionToEdit.RoleArn = roleArn
	sessionToEdit.Region = region
	sessionToEdit.ParentId = parentId
	sessionToEdit.ParentType = parentType
	sessionToEdit.NamedProfileId = namedProfileId
	return facade.replaceSession(sessionToEdit)
}

func (fac *AlibabaRamRoleChainedSessionsFacade) RemoveSession(id string) error {
	alibabaRamRoleChainedSessionsLock.Lock()
	defer alibabaRamRoleChainedSessionsLock.Unlock()

	currentSessions := fac.GetSessions()
	newSessions := make([]AlibabaRamRoleChainedSession, 0)

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

func (fac *AlibabaRamRoleChainedSessionsFacade) GetSessionById(id string) (*AlibabaRamRoleChainedSession, error) {
	for _, alibabaRamRoleChainedSession := range fac.GetSessions() {
		if alibabaRamRoleChainedSession.Id == id {
			return &alibabaRamRoleChainedSession, nil
		}
	}
	return nil, http_error.NewNotFoundError(fmt.Errorf("session with id %s not found", id))
}

func (fac *AlibabaRamRoleChainedSessionsFacade) SetSessionById(newSession *AlibabaRamRoleChainedSession) {
	currentSessions := fac.GetSessions()
	for i, session := range currentSessions {
		if session.Id == newSession.Id {
			currentSessions[i] = *newSession
		}
	}
	fac.SetSessions(currentSessions)
}

func (facade *AlibabaRamRoleChainedSessionsFacade) StartingSession(sessionId string) error {
	return facade.setSessionStatus(sessionId, domain_alibaba.Pending)
}

func (facade *AlibabaRamRoleChainedSessionsFacade) StartSession(sessionId string) error {
	return facade.setSessionStatus(sessionId, domain_alibaba.Active)
}

func (facade *AlibabaRamRoleChainedSessionsFacade) StopSession(sessionId string) error {
	return facade.setSessionStatus(sessionId, domain_alibaba.NotActive)
}
/*func (fac *AlibabaRamRoleChainedSessionsFacade) SetSessionStatusToPending(id string) error {

	alibabaRamRoleChainedSessionsLock.Lock()
	defer alibabaRamRoleChainedSessionsLock.Unlock()

	alibabaRamRoleChainedSession, err := fac.GetSessionById(id)
	if err != nil {
		return err
	}

	if !(alibabaRamRoleChainedSession.Status == domain_alibaba.NotActive) {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("session with id " + id + "cannot be started because it's in pending or active state"))
	}

	oldAlibabaRamRoleChainedSessions := fac.GetSessions()
	newAlibabaRamRoleChainedSessions := make([]AlibabaRamRoleChainedSession, 0)

	for i := range oldAlibabaRamRoleChainedSessions {
		newAlibabaRamRoleChainedSession := oldAlibabaRamRoleChainedSessions[i]
		newAlibabaRamRoleChainedSessions = append(newAlibabaRamRoleChainedSessions, newAlibabaRamRoleChainedSession)
	}

	for i, session := range newAlibabaRamRoleChainedSessions {
		if session.Id == id {
			newAlibabaRamRoleChainedSessions[i].Status = domain_alibaba.Pending
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

	if !(alibabaRamRoleChainedSession.Status == domain_alibaba.Pending) {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("session with id " + id + "cannot be started because it's not in pending state"))
	}

	oldAlibabaRamRoleChainedSessions := fac.GetSessions()
	newAlibabaRamRoleChainedSessions := make([]AlibabaRamRoleChainedSession, 0)

	for i := range oldAlibabaRamRoleChainedSessions {
		newAlibabaRamRoleChainedSession := oldAlibabaRamRoleChainedSessions[i]
		newAlibabaRamRoleChainedSessions = append(newAlibabaRamRoleChainedSessions, newAlibabaRamRoleChainedSession)
	}

	for i, session := range newAlibabaRamRoleChainedSessions {
		if session.Id == id {
			newAlibabaRamRoleChainedSessions[i].Status = domain_alibaba.Active
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
	if alibabaRamRoleChainedSession.Status != domain_alibaba.Active {
		fmt.Println(alibabaRamRoleChainedSession.Status)
		return http_error.NewUnprocessableEntityError(fmt.Errorf("session with id " + id + "cannot be stopped because it's not in active state"))
	}

	oldAlibabaRamRoleChainedSessions := fac.GetSessions()
	newAlibabaRamRoleChainedSessions := make([]AlibabaRamRoleChainedSession, 0)

	for i := range oldAlibabaRamRoleChainedSessions {
		newAlibabaRamRoleChainedSession := oldAlibabaRamRoleChainedSessions[i]
		newAlibabaRamRoleChainedSessions = append(newAlibabaRamRoleChainedSessions, newAlibabaRamRoleChainedSession)
	}

	for i, session := range newAlibabaRamRoleChainedSessions {
		if session.Id == id {
			newAlibabaRamRoleChainedSessions[i].Status = domain_alibaba.NotActive
		}
	}

	err = fac.updateState(newAlibabaRamRoleChainedSessions)
	if err != nil {
		return err
	}

	return nil
}*/

func (facade *AlibabaRamRoleChainedSessionsFacade) setSessionStatus(id string, status domain_alibaba.AlibabaSessionStatus) error {
	alibabaRamRoleChainedSessionsLock.Lock()
	defer alibabaRamRoleChainedSessionsLock.Unlock()

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

func (facade *AlibabaRamRoleChainedSessionsFacade) replaceSession(newSession *AlibabaRamRoleChainedSession) error {
	newSessions := make([]AlibabaRamRoleChainedSession, 0)
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
