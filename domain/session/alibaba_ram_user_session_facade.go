package session

import (
	"fmt"
	"leapp_daemon/infrastructure/http/http_error"
	"sync"
)

var alibabaRamUserSessionsFacadeSingleton *AlibabaRamUserSessionsFacade
var alibabaRamUserSessionsFacadeLock sync.Mutex
var alibabaRamUserSessionsLock sync.Mutex

type AlibabaRamUserSessionsObserver interface {
	UpdateAlibabaRamUserSessions(oldAlibabaRamUserSessions []AlibabaRamUserSession, newAlibabaRamUserSessions []AlibabaRamUserSession) error
}

type AlibabaRamUserSessionsFacade struct {
	alibabaRamUserSessions []AlibabaRamUserSession
	observers              []AlibabaRamUserSessionsObserver
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

func (fac *AlibabaRamUserSessionsFacade) SetSessions(newAlibabaRamUserSessions []AlibabaRamUserSession) error {
	fac.alibabaRamUserSessions = newAlibabaRamUserSessions

	err := fac.updateState(newAlibabaRamUserSessions)
	if err != nil {
		return err
	}
	return nil
}

func (fac *AlibabaRamUserSessionsFacade) UpdateSession(newSession AlibabaRamUserSession) error {
	allSessions := fac.GetSessions()
	for i, alibabaRamUserSession := range allSessions {
		if alibabaRamUserSession.Id == newSession.Id {
			allSessions[i] = newSession
		}
	}
	err := fac.SetSessions(allSessions)
	return err
}

func (fac *AlibabaRamUserSessionsFacade) AddSession(alibabaRamUserSession AlibabaRamUserSession) error {
	alibabaRamUserSessionsLock.Lock()
	defer alibabaRamUserSessionsLock.Unlock()

	oldAlibabaRamUserSessions := fac.GetSessions()
	newAlibabaRamUserSessions := make([]AlibabaRamUserSession, 0)

	for i := range oldAlibabaRamUserSessions {
		newAlibabaRamUserSession := oldAlibabaRamUserSessions[i]
		newAlibabaRamUserSessionAccount := *oldAlibabaRamUserSessions[i].Account
		newAlibabaRamUserSession.Account = &newAlibabaRamUserSessionAccount
		newAlibabaRamUserSessions = append(newAlibabaRamUserSessions, newAlibabaRamUserSession)
	}

	for _, sess := range newAlibabaRamUserSessions {
		if alibabaRamUserSession.Id == sess.Id {
			return http_error.NewConflictError(fmt.Errorf("a AlibabaRamUserSession with id " + alibabaRamUserSession.Id +
				" is already present"))
		}

		if alibabaRamUserSession.Alias == sess.Alias {
			return http_error.NewUnprocessableEntityError(fmt.Errorf("a session with the same alias " +
				"is already present"))
		}
	}

	newAlibabaRamUserSessions = append(newAlibabaRamUserSessions, alibabaRamUserSession)

	err := fac.updateState(newAlibabaRamUserSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *AlibabaRamUserSessionsFacade) RemoveSession(id string) error {
	alibabaRamUserSessionsLock.Lock()
	defer alibabaRamUserSessionsLock.Unlock()

	oldAlibabaRamUserSessions := fac.GetSessions()
	newAlibabaRamUserSessions := make([]AlibabaRamUserSession, 0)

	for i := range oldAlibabaRamUserSessions {
		newAlibabaRamUserSession := oldAlibabaRamUserSessions[i]
		newAlibabaRamUserSessionAccount := *oldAlibabaRamUserSessions[i].Account
		newAlibabaRamUserSession.Account = &newAlibabaRamUserSessionAccount
		newAlibabaRamUserSessions = append(newAlibabaRamUserSessions, newAlibabaRamUserSession)
	}

	for i, sess := range newAlibabaRamUserSessions {
		if sess.Id == id {
			newAlibabaRamUserSessions = append(newAlibabaRamUserSessions[:i], newAlibabaRamUserSessions[i+1:]...)
			break
		}
	}

	if len(fac.GetSessions()) == len(newAlibabaRamUserSessions) {
		return http_error.NewNotFoundError(fmt.Errorf("alibaba Ram User session with id %s not found", id))
	}

	err := fac.updateState(newAlibabaRamUserSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *AlibabaRamUserSessionsFacade) GetSessionById(id string) (*AlibabaRamUserSession, error) {
	for _, alibabaRamUserSession := range fac.GetSessions() {
		if alibabaRamUserSession.Id == id {
			return &alibabaRamUserSession, nil
		}
	}
	return nil, http_error.NewNotFoundError(fmt.Errorf("alibaba Ram User session with id %s not found", id))
}

func (fac *AlibabaRamUserSessionsFacade) SetStatusToPending(id string) error {
	alibabaRamUserSessionsLock.Lock()
	defer alibabaRamUserSessionsLock.Unlock()

	alibabaRamUserSession, err := fac.GetSessionById(id)
	if err != nil {
		return err
	}

	if !(alibabaRamUserSession.Status == NotActive) {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("Alibaba Ram User session with id " + id + "cannot be started because it's in pending or active state"))
	}

	oldAlibabaRamUserSessions := fac.GetSessions()
	newAlibabaRamUserSessions := make([]AlibabaRamUserSession, 0)

	for i := range oldAlibabaRamUserSessions {
		newAlibabaRamUserSession := oldAlibabaRamUserSessions[i]
		newAlibabaRamUserSessionAccount := *oldAlibabaRamUserSessions[i].Account
		newAlibabaRamUserSession.Account = &newAlibabaRamUserSessionAccount
		newAlibabaRamUserSessions = append(newAlibabaRamUserSessions, newAlibabaRamUserSession)
	}

	for i, session := range newAlibabaRamUserSessions {
		if session.Id == id {
			newAlibabaRamUserSessions[i].Status = Pending
		}
	}

	err = fac.updateState(newAlibabaRamUserSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *AlibabaRamUserSessionsFacade) SetStatusToActive(id string) error {
	alibabaRamUserSessionsLock.Lock()
	defer alibabaRamUserSessionsLock.Unlock()

	alibabaRamUserSession, err := fac.GetSessionById(id)
	if err != nil {
		return err
	}

	if !(alibabaRamUserSession.Status == Pending) {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("Alibaba Ram User session with id " + id + "cannot be started because it's not in pending state"))
	}

	oldAlibabaRamUserSessions := fac.GetSessions()
	newAlibabaRamUserSessions := make([]AlibabaRamUserSession, 0)

	for i := range oldAlibabaRamUserSessions {
		newAlibabaRamUserSession := oldAlibabaRamUserSessions[i]
		newAlibabaRamUserSessionAccount := *oldAlibabaRamUserSessions[i].Account
		newAlibabaRamUserSession.Account = &newAlibabaRamUserSessionAccount
		newAlibabaRamUserSessions = append(newAlibabaRamUserSessions, newAlibabaRamUserSession)
	}

	for i, session := range newAlibabaRamUserSessions {
		if session.Id == id {
			newAlibabaRamUserSessions[i].Status = Active
		}
	}

	err = fac.updateState(newAlibabaRamUserSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *AlibabaRamUserSessionsFacade) SetStatusToInactive(id string) error {
	alibabaRamUserSessionsLock.Lock()
	defer alibabaRamUserSessionsLock.Unlock()

	alibabaRamUserSession, err := fac.GetSessionById(id)
	if err != nil {
		return err
	}
	if alibabaRamUserSession.Status != Active {
		fmt.Println(alibabaRamUserSession.Status)
		return http_error.NewUnprocessableEntityError(fmt.Errorf("Alibaba Ram User session with id " + id + "cannot be started because it's not in active state"))
	}

	oldAlibabaRamUserSessions := fac.GetSessions()
	newAlibabaRamUserSessions := make([]AlibabaRamUserSession, 0)

	for i := range oldAlibabaRamUserSessions {
		newAlibabaRamUserSession := oldAlibabaRamUserSessions[i]
		newAlibabaRamUserSessionAccount := *oldAlibabaRamUserSessions[i].Account
		newAlibabaRamUserSession.Account = &newAlibabaRamUserSessionAccount
		newAlibabaRamUserSessions = append(newAlibabaRamUserSessions, newAlibabaRamUserSession)
	}

	for i, session := range newAlibabaRamUserSessions {
		if session.Id == id {
			newAlibabaRamUserSessions[i].Status = NotActive
		}
	}

	err = fac.updateState(newAlibabaRamUserSessions)
	if err != nil {
		return err
	}

	return nil
}

func (fac *AlibabaRamUserSessionsFacade) updateState(newState []AlibabaRamUserSession) error {
	oldAlibabaRamUserSessions := fac.GetSessions()
	fac.alibabaRamUserSessions = newState

	for _, observer := range fac.observers {
		err := observer.UpdateAlibabaRamUserSessions(oldAlibabaRamUserSessions, newState)
		if err != nil {
			return err
		}
	}

	return nil
}
