package mock

import (
	"errors"
	"fmt"
	"leapp_daemon/domain/domain_alibaba/alibaba_ram_role_chained"
	"leapp_daemon/infrastructure/http/http_error"
)

type AlibabaRamRoleChainedSessionsFacadeMock struct {
	calls                               []string
	ExpErrorOnGetSessionById            bool
	ExpErrorOnAddSession                bool
	ExpErrorOnRemoveSession             bool
	ExpErrorOnEditSession               bool
	ExpErrorOnStartSession              bool
	ExpErrorOnStartingSession           bool
	ExpErrorOnStopSession               bool
	ExpErrorOnSetSessionTokenExpiration bool
	ExpGetSessionById                   alibaba_ram_role_chained.AlibabaRamRoleChainedSession
	ExpGetSessions                      []alibaba_ram_role_chained.AlibabaRamRoleChainedSession
}

func NewAlibabaRamRoleChainedSessionsFacadeMock() AlibabaRamRoleChainedSessionsFacadeMock {
	return AlibabaRamRoleChainedSessionsFacadeMock{calls: []string{}, ExpGetSessions: []alibaba_ram_role_chained.AlibabaRamRoleChainedSession{}}
}

func (facade *AlibabaRamRoleChainedSessionsFacadeMock) GetCalls() []string {
	return facade.calls
}

func (facade *AlibabaRamRoleChainedSessionsFacadeMock) GetSessions() []alibaba_ram_role_chained.AlibabaRamRoleChainedSession {
	facade.calls = append(facade.calls, "GetSessions()")
	return facade.ExpGetSessions

}

func (facade *AlibabaRamRoleChainedSessionsFacadeMock) GetSessionById(sessionId string) (alibaba_ram_role_chained.AlibabaRamRoleChainedSession, error) {
	facade.calls = append(facade.calls, fmt.Sprintf("GetSessionById(%v)", sessionId))
	if facade.ExpErrorOnGetSessionById {
		return alibaba_ram_role_chained.AlibabaRamRoleChainedSession{}, http_error.NewNotFoundError(errors.New("session not found"))
	}
	return facade.ExpGetSessionById, nil
}

func (facade *AlibabaRamRoleChainedSessionsFacadeMock) AddSession(session alibaba_ram_role_chained.AlibabaRamRoleChainedSession) error {
	facade.calls = append(facade.calls, fmt.Sprintf("AddSession(%v)", session.Name))
	if facade.ExpErrorOnAddSession {
		return http_error.NewConflictError(errors.New("session already exist"))
	}
	return nil
}

func (facade *AlibabaRamRoleChainedSessionsFacadeMock) RemoveSession(sessionId string) error {
	facade.calls = append(facade.calls, fmt.Sprintf("RemoveSession(%v)", sessionId))
	if facade.ExpErrorOnRemoveSession {
		return http_error.NewNotFoundError(errors.New("session not found"))
	}
	return nil
}

func (facade *AlibabaRamRoleChainedSessionsFacadeMock) EditSession(sessionId string, sessionName string, region string,
	accountNumber string, userName string, mfaDevice string, namedProfileId string) error {
	facade.calls = append(facade.calls, fmt.Sprintf("EditSession(%v, %v, %v, %v, %v, %v, %v)",
		sessionId, sessionName, region, accountNumber, userName, mfaDevice, namedProfileId))
	if facade.ExpErrorOnEditSession {
		return http_error.NewConflictError(errors.New("unable to edit session, collision detected"))
	}

	return nil
}

func (facade *AlibabaRamRoleChainedSessionsFacadeMock) StartSession(sessionId string, startTime string) error {
	facade.calls = append(facade.calls, fmt.Sprintf("StartSession(%v, %v)", sessionId, startTime))
	if facade.ExpErrorOnStartSession {
		return http_error.NewInternalServerError(errors.New("unable to start the session"))
	}
	return nil
}

func (facade *AlibabaRamRoleChainedSessionsFacadeMock) StartingSession(sessionId string) error {
	facade.calls = append(facade.calls, fmt.Sprintf("StartingSession(%v)", sessionId))
	if facade.ExpErrorOnStartingSession {
		return http_error.NewInternalServerError(errors.New("starting session failed"))
	}
	return nil
}

func (facade *AlibabaRamRoleChainedSessionsFacadeMock) StopSession(sessionId string, stopTime string) error {
	facade.calls = append(facade.calls, fmt.Sprintf("StopSession(%v, %v)", sessionId, stopTime))
	if facade.ExpErrorOnStopSession {
		return http_error.NewInternalServerError(errors.New("unable to stop the session"))
	}
	return nil
}

func (facade *AlibabaRamRoleChainedSessionsFacadeMock) SetSessionTokenExpiration(sessionId string, sessionTokenExpiration string) error {
	facade.calls = append(facade.calls, fmt.Sprintf("SetSessionTokenExpiration(%v, %v)", sessionId, sessionTokenExpiration))
	if facade.ExpErrorOnSetSessionTokenExpiration {
		return http_error.NewInternalServerError(errors.New("unable to set token expiration"))
	}
	return nil
}
