package session

import (
	"fmt"
	"leapp_daemon/infrastructure/http/http_error"
	"sync"
	"time"
)

var awsIamUserSessionsLock sync.Mutex

type AwsIamUserSessionsObserver interface {
	UpdateAwsIamUserSessions(oldSessions []AwsIamUserSession, newSessions []AwsIamUserSession) error
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

func (fac *AwsIamUserSessionsFacade) Subscribe(observer AwsIamUserSessionsObserver) {
	fac.observers = append(fac.observers, observer)
}

func (fac *AwsIamUserSessionsFacade) GetSessions() []AwsIamUserSession {
	return fac.awsIamUserSessions
}

func (fac *AwsIamUserSessionsFacade) SetSessions(sessions []AwsIamUserSession) {
	fac.awsIamUserSessions = sessions
}

func (fac *AwsIamUserSessionsFacade) AddSession(session AwsIamUserSession) error {
	awsIamUserSessionsLock.Lock()
	defer awsIamUserSessionsLock.Unlock()

	oldSessions := fac.GetSessions()
	newSessions := make([]AwsIamUserSession, 0)

	for i := range oldSessions {
		newSession := oldSessions[i]
		newSessionAccount := *oldSessions[i].Account
		newSession.Account = &newSessionAccount
		newSessions = append(newSessions, newSession)
	}

	for _, sess := range newSessions {
		if session.Id == sess.Id {
			return http_error.NewConflictError(fmt.Errorf("a AwsIamUserSession with id " + session.Id +
				" is already present"))
		}

		if session.Alias == sess.Alias {
			return http_error.NewUnprocessableEntityError(fmt.Errorf("a session with the same alias " +
				"is already present"))
		}
	}

	newSessions = append(newSessions, session)

	err := fac.updateState(newSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *AwsIamUserSessionsFacade) RemoveSession(id string) error {
	awsIamUserSessionsLock.Lock()
	defer awsIamUserSessionsLock.Unlock()

	oldSessions := fac.GetSessions()
	newSessions := make([]AwsIamUserSession, 0)

	for i := range oldSessions {
		newSession := oldSessions[i]
		newSessionAccount := *oldSessions[i].Account
		newSession.Account = &newSessionAccount
		newSessions = append(newSessions, newSession)
	}

	for i, sess := range newSessions {
		if sess.Id == id {
			newSessions = append(newSessions[:i], newSessions[i+1:]...)
			break
		}
	}

	if len(fac.GetSessions()) == len(newSessions) {
		return http_error.NewNotFoundError(fmt.Errorf("aws iam user session with id %s not found", id))
	}

	err := fac.updateState(newSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *AwsIamUserSessionsFacade) GetSessionById(id string) (*AwsIamUserSession, error) {
	for _, session := range fac.GetSessions() {
		if session.Id == id {
			return &session, nil
		}
	}
	return nil, http_error.NewNotFoundError(fmt.Errorf("aws iam user session with id %s not found", id))
}

func (fac *AwsIamUserSessionsFacade) SetSessionStatusToPending(id string) error {
	awsIamUserSessionsLock.Lock()
	defer awsIamUserSessionsLock.Unlock()

	session, err := fac.GetSessionById(id)
	if err != nil {
		return err
	}

	if !(session.Status == NotActive) {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("aws iam user session with id " + id + "cannot be started because it's in pending or active state"))
	}

	oldSessions := fac.GetSessions()
	newSessions := make([]AwsIamUserSession, 0)

	for i := range oldSessions {
		newSession := oldSessions[i]
		newSessionAccount := *oldSessions[i].Account
		newSession.Account = &newSessionAccount
		newSessions = append(newSessions, newSession)
	}

	for i, session := range newSessions {
		if session.Id == id {
			newSessions[i].Status = Pending
		}
	}

	err = fac.updateState(newSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *AwsIamUserSessionsFacade) SetSessionStatusToActive(id string) error {
	awsIamUserSessionsLock.Lock()
	defer awsIamUserSessionsLock.Unlock()

	session, err := fac.GetSessionById(id)
	if err != nil {
		return err
	}

	if !(session.Status == Pending) {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("aws iam user session with id " + id + "cannot be started because it's not in pending state"))
	}

	oldSessions := fac.GetSessions()
	newSessions := make([]AwsIamUserSession, 0)

	for i := range oldSessions {
		newSession := oldSessions[i]
		newSessionAccount := *oldSessions[i].Account
		newSession.Account = &newSessionAccount
		newSessions = append(newSessions, newSession)
	}

	for i, session := range newSessions {
		if session.Id == id {
			newSessions[i].Status = Active
			newSessions[i].StartTime = time.Now().Format(time.RFC3339)
		}
	}

	err = fac.updateState(newSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *AwsIamUserSessionsFacade) SetSessionTokenExpiration(sessionId string, sessionTokenExpiration time.Time) error {
	awsIamUserSessionsLock.Lock()
	defer awsIamUserSessionsLock.Unlock()

	oldSessions := fac.GetSessions()
	newSessions := make([]AwsIamUserSession, 0)

	for i := range oldSessions {
		newSession := oldSessions[i]
		newSessionAccount := *oldSessions[i].Account
		newSession.Account = &newSessionAccount
		newSessions = append(newSessions, newSession)
	}

	for i, session := range newSessions {
		if session.Id == sessionId {
			newSessions[i].Account.SessionTokenExpiration = sessionTokenExpiration.Format(time.RFC3339)

			err := fac.updateState(newSessions)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return http_error.NewNotFoundError(fmt.Errorf("aws iam user session with sessionId %s not found", sessionId))
}

func (fac *AwsIamUserSessionsFacade) updateState(newState []AwsIamUserSession) error {
	oldSessions := fac.GetSessions()
	fac.SetSessions(newState)

	for _, observer := range fac.observers {
		err := observer.UpdateAwsIamUserSessions(oldSessions, newState)
		if err != nil {
			return err
		}
	}

	return nil
}
