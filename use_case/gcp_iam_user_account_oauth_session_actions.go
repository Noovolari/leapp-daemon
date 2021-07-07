package use_case

import (
	"leapp_daemon/domain/domain_gcp"
	"leapp_daemon/domain/domain_gcp/gcp_iam_user_account_oauth"
	"leapp_daemon/infrastructure/http/http_error"
)

type GcpIamUserAccountOauthSessionActions struct {
	GcpApi                              GcpApi
	Environment                         Environment
	Keychain                            Keychain
	GcpIamUserAccountOauthSessionFacade GcpIamUserAccountOauthSessionsFacade
}

func (actions *GcpIamUserAccountOauthSessionActions) GetSession(sessionId string) (gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession, error) {
	return actions.GcpIamUserAccountOauthSessionFacade.GetSessionById(sessionId)
}

func (actions *GcpIamUserAccountOauthSessionActions) GetOAuthUrl() (string, error) {
	return actions.GcpApi.GetOauthUrl()
}

func (actions *GcpIamUserAccountOauthSessionActions) CreateSession(name string, accountId string, projectName string, oauthCode string) error {

	newSessionId := actions.Environment.GenerateUuid()
	credentialsLabel := newSessionId + "-gcp-iam-user-account-oauth-session-credentials"

	gcpSession := gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{
		Id:               newSessionId,
		Name:             name,
		AccountId:        accountId,
		ProjectName:      projectName,
		CredentialsLabel: credentialsLabel,
		Status:           domain_gcp.NotActive,
		StartTime:        "",
		LastStopTime:     "",
	}

	token, err := actions.GcpApi.GetOauthToken(oauthCode)
	if err != nil {
		return err
	}

	credentials := actions.GcpApi.GetCredentials(token)

	err = actions.Keychain.SetSecret(credentials, credentialsLabel)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return actions.GcpIamUserAccountOauthSessionFacade.AddSession(gcpSession)
}

func (actions *GcpIamUserAccountOauthSessionActions) StartSession(sessionId string) error {

	facade := actions.GcpIamUserAccountOauthSessionFacade
	currentTime := actions.Environment.GetTime()

	for _, currentSession := range facade.GetSessions() {
		if currentSession.Status != domain_gcp.NotActive && currentSession.Id != sessionId {
			err := facade.StopSession(currentSession.Id, currentTime)
			if err != nil {
				return err
			}
		}
	}
	return facade.StartSession(sessionId, currentTime)
}

func (actions *GcpIamUserAccountOauthSessionActions) StopSession(sessionId string) error {
	return actions.GcpIamUserAccountOauthSessionFacade.StopSession(sessionId, actions.Environment.GetTime())
}

func (actions *GcpIamUserAccountOauthSessionActions) DeleteSession(sessionId string) error {
	facade := actions.GcpIamUserAccountOauthSessionFacade

	sessionToDelete, err := facade.GetSessionById(sessionId)
	if err != nil {
		return err
	}

	_ = actions.Keychain.DeleteSecret(sessionToDelete.CredentialsLabel)
	return facade.RemoveSession(sessionId)
}

func (actions *GcpIamUserAccountOauthSessionActions) EditSession(sessionId string, name string, projectName string) error {
	sessionsFacade := actions.GcpIamUserAccountOauthSessionFacade

	return sessionsFacade.EditSession(sessionId, name, projectName)
}
