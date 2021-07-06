package gcp_iam_user_account_oauth

import (
	"fmt"
	"leapp_daemon/domain/domain_gcp"
	"leapp_daemon/infrastructure/http/http_error"
	"sync"
)

var gcpIamUserAccountOauthSessionsLock sync.Mutex

type GcpIamUserAccountOauthSessionsObserver interface {
	UpdateGcpIamUserAccountOauthSessions(oldSessions []GcpIamUserAccountOauthSession, newSessions []GcpIamUserAccountOauthSession)
}

type GcpIamUserAccountOauthSessionsFacade struct {
	gcpIamUserAccountOauthSessions []GcpIamUserAccountOauthSession
	observers                      []GcpIamUserAccountOauthSessionsObserver
}

func NewGcpIamUserAccountOauthSessionsFacade() *GcpIamUserAccountOauthSessionsFacade {
	return &GcpIamUserAccountOauthSessionsFacade{
		gcpIamUserAccountOauthSessions: make([]GcpIamUserAccountOauthSession, 0),
	}
}

func (facade *GcpIamUserAccountOauthSessionsFacade) Subscribe(observer GcpIamUserAccountOauthSessionsObserver) {
	facade.observers = append(facade.observers, observer)
}

func (facade *GcpIamUserAccountOauthSessionsFacade) GetSessions() []GcpIamUserAccountOauthSession {
	return facade.gcpIamUserAccountOauthSessions
}

func (facade *GcpIamUserAccountOauthSessionsFacade) GetSessionById(sessionId string) (GcpIamUserAccountOauthSession, error) {
	for _, session := range facade.GetSessions() {
		if session.Id == sessionId {
			return session, nil
		}
	}

	return GcpIamUserAccountOauthSession{}, http_error.NewNotFoundError(fmt.Errorf("session with id %s not found", sessionId))
}

func (facade *GcpIamUserAccountOauthSessionsFacade) SetSessions(sessions []GcpIamUserAccountOauthSession) {
	gcpIamUserAccountOauthSessionsLock.Lock()
	defer gcpIamUserAccountOauthSessionsLock.Unlock()

	facade.gcpIamUserAccountOauthSessions = sessions
}

func (facade *GcpIamUserAccountOauthSessionsFacade) AddSession(newSession GcpIamUserAccountOauthSession) error {
	gcpIamUserAccountOauthSessionsLock.Lock()
	defer gcpIamUserAccountOauthSessionsLock.Unlock()

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

func (facade *GcpIamUserAccountOauthSessionsFacade) RemoveSession(sessionId string) error {
	gcpIamUserAccountOauthSessionsLock.Lock()
	defer gcpIamUserAccountOauthSessionsLock.Unlock()

	currentSessions := facade.GetSessions()
	newSessions := make([]GcpIamUserAccountOauthSession, 0)

	for _, session := range currentSessions {
		if session.Id != sessionId {
			newSessions = append(newSessions, session)
		}
	}

	if len(currentSessions) == len(newSessions) {
		return http_error.NewNotFoundError(fmt.Errorf("session with id %s not found", sessionId))
	}

	facade.updateState(newSessions)
	return nil
}

func (facade *GcpIamUserAccountOauthSessionsFacade) EditSession(sessionId string, sessionName string, projectName string, namedConfigurationId string) error {
	gcpIamUserAccountOauthSessionsLock.Lock()
	defer gcpIamUserAccountOauthSessionsLock.Unlock()

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
	sessionToEdit.ProjectName = projectName
	sessionToEdit.NamedConfigurationId = namedConfigurationId

	return facade.replaceSession(sessionId, sessionToEdit)
}

func (facade *GcpIamUserAccountOauthSessionsFacade) StartSession(sessionId string, startTime string) error {
	return facade.setSessionStatus(sessionId, domain_gcp.Active, startTime, "")
}

func (facade *GcpIamUserAccountOauthSessionsFacade) StopSession(sessionId string, stopTime string) error {
	return facade.setSessionStatus(sessionId, domain_gcp.NotActive, "", stopTime)
}

func (facade *GcpIamUserAccountOauthSessionsFacade) setSessionStatus(sessionId string, status domain_gcp.GcpSessionStatus, startTime string, lastStopTime string) error {
	gcpIamUserAccountOauthSessionsLock.Lock()
	defer gcpIamUserAccountOauthSessionsLock.Unlock()

	sessionToUpdate, err := facade.GetSessionById(sessionId)
	if err != nil {
		return err
	}

	sessionToUpdate.Status = status
	if startTime != "" {
		sessionToUpdate.StartTime = startTime
		sessionToUpdate.LastStopTime = ""
	}
	if lastStopTime != "" {
		sessionToUpdate.StartTime = ""
		sessionToUpdate.LastStopTime = lastStopTime
	}
	return facade.replaceSession(sessionId, sessionToUpdate)
}

func (facade *GcpIamUserAccountOauthSessionsFacade) replaceSession(sessionId string, newSession GcpIamUserAccountOauthSession) error {
	newSessions := make([]GcpIamUserAccountOauthSession, 0)
	for _, session := range facade.GetSessions() {
		if session.Id == sessionId {
			newSessions = append(newSessions, newSession)
		} else {
			newSessions = append(newSessions, session)
		}
	}

	facade.updateState(newSessions)
	return nil
}

func (facade *GcpIamUserAccountOauthSessionsFacade) updateState(newSessions []GcpIamUserAccountOauthSession) {
	oldSessions := facade.GetSessions()
	facade.gcpIamUserAccountOauthSessions = newSessions

	for _, observer := range facade.observers {
		observer.UpdateGcpIamUserAccountOauthSessions(oldSessions, newSessions)
	}
}
