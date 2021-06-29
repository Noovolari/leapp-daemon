package session

import (
	"fmt"
	"leapp_daemon/infrastructure/http/http_error"
	"sync"
)

var awsIamUserSessionsLock sync.Mutex

type AwsIamUserSessionsObserver interface {
	UpdateAwsIamUserSessions(oldSessions []AwsIamUserSession, newSessions []AwsIamUserSession)
}

type AwsIamUserSessionsFacade struct {
	awsIamUserSessions []AwsIamUserSession
	observers          []AwsIamUserSessionsObserver
}

func NewAwsIamUserSessionsFacade() *AwsIamUserSessionsFacade {
	return &AwsIamUserSessionsFacade{
		awsIamUserSessions: make([]AwsIamUserSession, 0),
	}
}

func (facade *AwsIamUserSessionsFacade) Subscribe(observer AwsIamUserSessionsObserver) {
	facade.observers = append(facade.observers, observer)
}

func (facade *AwsIamUserSessionsFacade) GetSessions() []AwsIamUserSession {
	return facade.awsIamUserSessions
}

func (facade *AwsIamUserSessionsFacade) GetSessionById(sessionId string) (AwsIamUserSession, error) {
	for _, session := range facade.GetSessions() {
		if session.Id == sessionId {
			return session, nil
		}
	}

	return AwsIamUserSession{}, http_error.NewNotFoundError(fmt.Errorf("session with id %s not found", sessionId))
}

func (facade *AwsIamUserSessionsFacade) SetSessions(sessions []AwsIamUserSession) {
	awsIamUserSessionsLock.Lock()
	defer awsIamUserSessionsLock.Unlock()

	facade.awsIamUserSessions = sessions
}

func (facade *AwsIamUserSessionsFacade) AddSession(newSession AwsIamUserSession) error {
	awsIamUserSessionsLock.Lock()
	defer awsIamUserSessionsLock.Unlock()

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

func (facade *AwsIamUserSessionsFacade) RemoveSession(sessionId string) error {
	awsIamUserSessionsLock.Lock()
	defer awsIamUserSessionsLock.Unlock()

	currentSessions := facade.GetSessions()
	newSessions := make([]AwsIamUserSession, 0)

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

func (facade *AwsIamUserSessionsFacade) EditSession(sessionId string, sessionName string, region string,
	accessKeyIdLabel string, secretKeyLabel string, sessionTokenLabel string, mfaDevice string,
	sessionTokenExpiration string, namedProfileId string) error {
	awsIamUserSessionsLock.Lock()
	defer awsIamUserSessionsLock.Unlock()

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
	sessionToEdit.AccessKeyIdLabel = accessKeyIdLabel
	sessionToEdit.SecretKeyLabel = secretKeyLabel
	sessionToEdit.SessionTokenLabel = sessionTokenLabel
	sessionToEdit.MfaDevice = mfaDevice
	sessionToEdit.SessionTokenExpiration = sessionTokenExpiration
	sessionToEdit.NamedProfileId = namedProfileId
	return facade.replaceSession(sessionId, sessionToEdit)
}

func (facade *AwsIamUserSessionsFacade) SetSessionTokenExpiration(sessionId string, sessionTokenExpiration string) error {
	awsIamUserSessionsLock.Lock()
	defer awsIamUserSessionsLock.Unlock()

	sessionToEdit, err := facade.GetSessionById(sessionId)
	if err != nil {
		return err
	}

	sessionToEdit.SessionTokenExpiration = sessionTokenExpiration
	return facade.replaceSession(sessionId, sessionToEdit)
}

func (facade *AwsIamUserSessionsFacade) StartingSession(sessionId string) error {
	return facade.setSessionStatus(sessionId, Pending, "", "")
}

func (facade *AwsIamUserSessionsFacade) StartSession(sessionId string, startTime string) error {
	return facade.setSessionStatus(sessionId, Active, startTime, "")
}

func (facade *AwsIamUserSessionsFacade) StopSession(sessionId string, stopTime string) error {
	return facade.setSessionStatus(sessionId, NotActive, "", stopTime)
}

func (facade *AwsIamUserSessionsFacade) setSessionStatus(sessionId string, status Status, startTime string, lastStopTime string) error {
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

func (facade *AwsIamUserSessionsFacade) replaceSession(sessionId string, newSession AwsIamUserSession) error {
	newSessions := make([]AwsIamUserSession, 0)
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

func (facade *AwsIamUserSessionsFacade) updateState(newSessions []AwsIamUserSession) {
	oldSessions := facade.GetSessions()
	facade.awsIamUserSessions = newSessions

	for _, observer := range facade.observers {
		observer.UpdateAwsIamUserSessions(oldSessions, newSessions)
	}
}
