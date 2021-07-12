package alibaba_ram_user

import (
	"fmt"
	"leapp_daemon/domain/domain_alibaba"
	"leapp_daemon/infrastructure/http/http_error"
	"sync"
)

var alibabaRamUserSessionsFacadeSingleton *AlibabaRamUserSessionsFacade
var alibabaRamUserSessionsFacadeLock sync.Mutex
var alibabaRamUserSessionsLock sync.Mutex

type AlibabaRamUserSessionsObserver interface {
	UpdateAlibabaRamUserSessions([]AlibabaRamUserSession, []AlibabaRamUserSession) error
}

type AlibabaRamUserSessionsFacade struct {
	alibabaRamUserSessions []AlibabaRamUserSession
	observers              []AlibabaRamUserSessionsObserver
}

func NewAlibabaRamUserSessionsFacade() *AlibabaRamUserSessionsFacade {
	return &AlibabaRamUserSessionsFacade{
		alibabaRamUserSessions: make([]AlibabaRamUserSession, 0),
	}
}

func GetAlibabaRamUserSessionsFacade() *AlibabaRamUserSessionsFacade {
	alibabaRamUserSessionsFacadeLock.Lock()
	defer alibabaRamUserSessionsFacadeLock.Unlock()

	if alibabaRamUserSessionsFacadeSingleton == nil {
		alibabaRamUserSessionsFacadeSingleton = &AlibabaRamUserSessionsFacade{
			alibabaRamUserSessions: make([]AlibabaRamUserSession, 0),
		}
	}

	return alibabaRamUserSessionsFacadeSingleton
}

func (fac *AlibabaRamUserSessionsFacade) Subscribe(observer AlibabaRamUserSessionsObserver) {
	fac.observers = append(fac.observers, observer)
}

func (fac *AlibabaRamUserSessionsFacade) GetSessions() []AlibabaRamUserSession {
	return fac.alibabaRamUserSessions
}

func (fac *AlibabaRamUserSessionsFacade) SetSessions(sessions []AlibabaRamUserSession) error {
	fac.alibabaRamUserSessions = sessions

	fac.updateState(sessions)
	return nil
}

func (facade *AlibabaRamUserSessionsFacade) EditSession(sessionId string, sessionName string, region string, namedProfileId string) error {
	alibabaRamUserSessionsLock.Lock()
	defer alibabaRamUserSessionsLock.Unlock()

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
	sessionToEdit.Region = region
	sessionToEdit.NamedProfileId = namedProfileId
	return facade.replaceSession(sessionToEdit)
}

func (facade *AlibabaRamUserSessionsFacade) AddSession(newSession AlibabaRamUserSession) error {
	alibabaRamUserSessionsLock.Lock()
	defer alibabaRamUserSessionsLock.Unlock()

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

func (fac *AlibabaRamUserSessionsFacade) RemoveSession(id string) error {
	alibabaRamUserSessionsLock.Lock()
	defer alibabaRamUserSessionsLock.Unlock()

	currentSessions := fac.GetSessions()
	newSessions := make([]AlibabaRamUserSession, 0)

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

func (fac *AlibabaRamUserSessionsFacade) GetSessionById(id string) (*AlibabaRamUserSession, error) {
	for _, alibabaRamUserSession := range fac.GetSessions() {
		if alibabaRamUserSession.Id == id {
			return &alibabaRamUserSession, nil
		}
	}
	return nil, http_error.NewNotFoundError(fmt.Errorf("session with id %s not found", id))
}

func (facade *AlibabaRamUserSessionsFacade) StartingSession(sessionId string) error {
	return facade.setSessionStatus(sessionId, domain_alibaba.Pending)
}

func (facade *AlibabaRamUserSessionsFacade) StartSession(sessionId string) error {
	return facade.setSessionStatus(sessionId, domain_alibaba.Active)
}

func (facade *AlibabaRamUserSessionsFacade) StopSession(sessionId string) error {
	return facade.setSessionStatus(sessionId, domain_alibaba.NotActive)
}

/*func (fac *AlibabaRamUserSessionsFacade) SetSessionStatusToPending(id string) error {
	alibabaRamUserSessionsLock.Lock()
	defer alibabaRamUserSessionsLock.Unlock()

	alibabaRamUserSession, err := fac.GetSessionById(id)
	if err != nil {
		return err
	}

	if !(alibabaRamUserSession.Status == domain_alibaba.NotActive) {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("session with id " + id + "cannot be started because it's in pending or active state"))
	}

	oldAlibabaRamUserSessions := fac.GetSessions()
	newAlibabaRamUserSessions := make([]AlibabaRamUserSession, 0)

	for i := range oldAlibabaRamUserSessions {
		newAlibabaRamUserSession := oldAlibabaRamUserSessions[i]
		newAlibabaRamUserSessions = append(newAlibabaRamUserSessions, newAlibabaRamUserSession)
	}

	for i, session := range newAlibabaRamUserSessions {
		if session.Id == id {
			newAlibabaRamUserSessions[i].Status = domain_alibaba.Pending
		}
	}

	err = fac.updateState(newAlibabaRamUserSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *AlibabaRamUserSessionsFacade) SetSessionStatusToActive(id string) error {
	alibabaRamUserSessionsLock.Lock()
	defer alibabaRamUserSessionsLock.Unlock()

	alibabaRamUserSession, err := fac.GetSessionById(id)
	if err != nil {
		return err
	}

	if !(alibabaRamUserSession.Status == domain_alibaba.Pending) {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("Alibaba Ram User session with id " + id + "cannot be started because it's not in pending state"))
	}

	oldAlibabaRamUserSessions := fac.GetSessions()
	newAlibabaRamUserSessions := make([]AlibabaRamUserSession, 0)

	for i := range oldAlibabaRamUserSessions {
		newAlibabaRamUserSession := oldAlibabaRamUserSessions[i]
		newAlibabaRamUserSessions = append(newAlibabaRamUserSessions, newAlibabaRamUserSession)
	}

	for i, session := range newAlibabaRamUserSessions {
		if session.Id == id {
			newAlibabaRamUserSessions[i].Status = domain_alibaba.Active
		}
	}

	err = fac.updateState(newAlibabaRamUserSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *AlibabaRamUserSessionsFacade) SetSessionStatusToInactive(id string) error {
	alibabaRamUserSessionsLock.Lock()
	defer alibabaRamUserSessionsLock.Unlock()

	alibabaRamUserSession, err := fac.GetSessionById(id)
	if err != nil {
		return err
	}
	if alibabaRamUserSession.Status != domain_alibaba.Active {
		fmt.Println(alibabaRamUserSession.Status)
		return http_error.NewUnprocessableEntityError(fmt.Errorf("Alibaba Ram User session with id " + id + "cannot be started because it's not in active state"))
	}

	oldAlibabaRamUserSessions := fac.GetSessions()
	newAlibabaRamUserSessions := make([]AlibabaRamUserSession, 0)

	for i := range oldAlibabaRamUserSessions {
		newAlibabaRamUserSession := oldAlibabaRamUserSessions[i]
		newAlibabaRamUserSessions = append(newAlibabaRamUserSessions, newAlibabaRamUserSession)
	}

	for i, session := range newAlibabaRamUserSessions {
		if session.Id == id {
			newAlibabaRamUserSessions[i].Status = domain_alibaba.NotActive
		}
	}

	err = fac.updateState(newAlibabaRamUserSessions)
	if err != nil {
		return err
	}

	return nil
}*/

func (facade *AlibabaRamUserSessionsFacade) setSessionStatus(id string, status domain_alibaba.AlibabaSessionStatus) error {
	alibabaRamUserSessionsLock.Lock()
	defer alibabaRamUserSessionsLock.Unlock()

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

func (facade *AlibabaRamUserSessionsFacade) replaceSession(newSession *AlibabaRamUserSession) error {
	newSessions := make([]AlibabaRamUserSession, 0)
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

func (fac *AlibabaRamUserSessionsFacade) updateState(newSessions []AlibabaRamUserSession) {
	oldSessions := fac.GetSessions()
	fac.alibabaRamUserSessions = newSessions

	for _, observer := range fac.observers {
		observer.UpdateAlibabaRamUserSessions(oldSessions, newSessions)
	}
}
