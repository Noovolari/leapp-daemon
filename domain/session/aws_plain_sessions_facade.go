package session

import (
	"fmt"
	"leapp_daemon/infrastructure/http/http_error"
	"sync"
	"time"
)

var awsPlainSessionsLock sync.Mutex

type AwsPlainSessionsObserver interface {
	UpdateAwsPlainSessions(oldSessions []AwsPlainSession, newSessions []AwsPlainSession) error
}

type AwsPlainSessionsFacade struct {
	awsPlainSessions []AwsPlainSession
	observers        []AwsPlainSessionsObserver
}

func NewAwsPlainSessionsFacade() *AwsPlainSessionsFacade {
	return &AwsPlainSessionsFacade{
		awsPlainSessions: make([]AwsPlainSession, 0),
	}
}

func (fac *AwsPlainSessionsFacade) Subscribe(observer AwsPlainSessionsObserver) {
	fac.observers = append(fac.observers, observer)
}

func (fac *AwsPlainSessionsFacade) GetSessions() []AwsPlainSession {
	return fac.awsPlainSessions
}

func (fac *AwsPlainSessionsFacade) SetSessions(sessions []AwsPlainSession) {
	fac.awsPlainSessions = sessions
}

func (fac *AwsPlainSessionsFacade) AddSession(session AwsPlainSession) error {
	awsPlainSessionsLock.Lock()
	defer awsPlainSessionsLock.Unlock()

	oldSessions := fac.GetSessions()
	newSessions := make([]AwsPlainSession, 0)

	for i := range oldSessions {
		newSession := oldSessions[i]
		newSessionAccount := *oldSessions[i].Account
		newSession.Account = &newSessionAccount
		newSessions = append(newSessions, newSession)
	}

	for _, sess := range newSessions {
		if session.Id == sess.Id {
			return http_error.NewConflictError(fmt.Errorf("a AwsPlainSession with id " + session.Id +
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

func (fac *AwsPlainSessionsFacade) RemoveSession(id string) error {
	awsPlainSessionsLock.Lock()
	defer awsPlainSessionsLock.Unlock()

	oldSessions := fac.GetSessions()
	newSessions := make([]AwsPlainSession, 0)

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
		return http_error.NewNotFoundError(fmt.Errorf("aws plain session with id %s not found", id))
	}

	err := fac.updateState(newSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *AwsPlainSessionsFacade) GetSessionById(id string) (*AwsPlainSession, error) {
	for _, session := range fac.GetSessions() {
		if session.Id == id {
			return &session, nil
		}
	}
	return nil, http_error.NewNotFoundError(fmt.Errorf("aws plain session with id %s not found", id))
}

func (fac *AwsPlainSessionsFacade) SetSessionStatusToPending(id string) error {
	awsPlainSessionsLock.Lock()
	defer awsPlainSessionsLock.Unlock()

	session, err := fac.GetSessionById(id)
	if err != nil {
		return err
	}

	if !(session.Status == NotActive) {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("aws plain session with id " + id + "cannot be started because it's in pending or active state"))
	}

	oldSessions := fac.GetSessions()
	newSessions := make([]AwsPlainSession, 0)

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

func (fac *AwsPlainSessionsFacade) SetSessionStatusToActive(id string) error {
	awsPlainSessionsLock.Lock()
	defer awsPlainSessionsLock.Unlock()

	session, err := fac.GetSessionById(id)
	if err != nil {
		return err
	}

	if !(session.Status == Pending) {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("aws plain session with id " + id + "cannot be started because it's not in pending state"))
	}

	oldSessions := fac.GetSessions()
	newSessions := make([]AwsPlainSession, 0)

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

func (fac *AwsPlainSessionsFacade) SetSessionTokenExpiration(sessionId string, sessionTokenExpiration time.Time) error {
	awsPlainSessionsLock.Lock()
	defer awsPlainSessionsLock.Unlock()

	oldSessions := fac.GetSessions()
	newSessions := make([]AwsPlainSession, 0)

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

	return http_error.NewNotFoundError(fmt.Errorf("aws plain session with sessionId %s not found", sessionId))
}

func (fac *AwsPlainSessionsFacade) updateState(newState []AwsPlainSession) error {
	oldSessions := fac.GetSessions()
	fac.SetSessions(newState)

	for _, observer := range fac.observers {
		err := observer.UpdateAwsPlainSessions(oldSessions, newState)
		if err != nil {
			return err
		}
	}

	return nil
}
