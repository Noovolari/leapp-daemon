package use_case

import (
	"golang.org/x/oauth2"
	"leapp_daemon/domain/domain_gcp"
	"leapp_daemon/domain/domain_gcp/gcp_iam_user_account_oauth"
	"leapp_daemon/test"
	"leapp_daemon/test/mock"
	"net/http"
	"reflect"
	"testing"
)

var (
	gcpApiMock                                       mock.GcpApiMock
	envMock                                          mock.EnvironmentMock
	gcpIamUserAccountOauthSessionActionsKeychainMock mock.KeychainMock
	gcpIamUserAccountOauthSessionFacadeMock          mock.GcpIamUserAccountOauthSessionsFacadeMock
	gcpIamUserAccountOauthSessionActions             *GcpIamUserAccountOauthSessionActions
)

func gcpIamUserAccountOauthSessionActionsSetup() {
	gcpApiMock = mock.NewGcpApiMock()
	envMock = mock.NewEnvironmentMock()
	gcpIamUserAccountOauthSessionActionsKeychainMock = mock.NewKeychainMock()
	gcpIamUserAccountOauthSessionFacadeMock = mock.NewGcpIamUserAccountOauthSessionsFacadeMock()
	gcpIamUserAccountOauthSessionActions = &GcpIamUserAccountOauthSessionActions{
		GcpApi:                              &gcpApiMock,
		Environment:                         &envMock,
		Keychain:                            &gcpIamUserAccountOauthSessionActionsKeychainMock,
		GcpIamUserAccountOauthSessionFacade: &gcpIamUserAccountOauthSessionFacadeMock,
	}
}

func gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t *testing.T, gcpApiMockCalls, envMockCalls, keychainMockCalls, facadeMockCalls []string) {
	if !reflect.DeepEqual(gcpApiMock.GetCalls(), gcpApiMockCalls) {
		t.Fatalf("gcpApiMock expectation violation.\nMock calls: %v", gcpApiMock.GetCalls())
	}
	if !reflect.DeepEqual(envMock.GetCalls(), envMockCalls) {
		t.Fatalf("envMock expectation violation.\nMock calls: %v", envMock.GetCalls())
	}
	if !reflect.DeepEqual(gcpIamUserAccountOauthSessionActionsKeychainMock.GetCalls(), keychainMockCalls) {
		t.Fatalf("keychainMock expectation violation.\nMock calls: %v", gcpIamUserAccountOauthSessionActionsKeychainMock.GetCalls())
	}
	if !reflect.DeepEqual(gcpIamUserAccountOauthSessionFacadeMock.GetCalls(), facadeMockCalls) {
		t.Fatalf("facadeMock expectation violation.\nMock calls: %v", gcpIamUserAccountOauthSessionFacadeMock.GetCalls())
	}
}

func TestGetSession(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()

	session := gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{Name: "test_session"}
	gcpIamUserAccountOauthSessionFacadeMock.ExpGetSessionById = session

	actualSession, err := gcpIamUserAccountOauthSessionActions.GetSession("ID")
	if err != nil && !reflect.DeepEqual(session, actualSession) {
		t.Fatalf("Returned unexpected session")
	}
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{}, []string{"GetSessionById(ID)"})
}

func TestGetSession_SessionFacadeReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	gcpIamUserAccountOauthSessionFacadeMock.ExpErrorOnGetSessionById = true

	_, err := gcpIamUserAccountOauthSessionActions.GetSession("ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session not found")
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{}, []string{"GetSessionById(ID)"})
}

func TestGetOAuthUrl(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	gcpApiMock.ExpOauthUrl = "url"

	actualOauthUrl, err := gcpIamUserAccountOauthSessionActions.GetOAuthUrl()
	if err != nil && !reflect.DeepEqual("url", actualOauthUrl) {
		t.Fatalf("Returned unexpected oauth url")
	}
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{"GetOauthUrl()"}, []string{}, []string{}, []string{})
}

func TestGetOAuthUrl_GcpApiReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	gcpApiMock.ExpErrorOnGetOauthUrl = true

	_, err := gcpIamUserAccountOauthSessionActions.GetOAuthUrl()
	test.ExpectHttpError(t, err, http.StatusNotFound, "error getting oauth url")
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{"GetOauthUrl()"}, []string{}, []string{}, []string{})
}

func TestCreateSession(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()

	sessionName := "sessionName"
	accountId := "accountId"
	projectName := "projectName"
	oauthCode := "oauthCode"
	uuid := "uuid"
	credentials := "credentialsJson"
	envMock.ExpUuid = uuid
	gcpApiMock.ExpOauthToken = oauth2.Token{}
	gcpApiMock.ExpCredentials = credentials

	err := gcpIamUserAccountOauthSessionActions.CreateSession(sessionName, accountId, projectName, oauthCode)
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{"GetOauthUrl(oauthCode)", "GetCredentials()"}, []string{"GenerateUuid()"}, []string{"SetSecret(credentialsJson, uuid-gcp-iam-user-account-oauth-session-credentials)"}, []string{"AddSession(sessionName)"})
}

func TestCreateSession_GcpApiGetOauthTokenReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()

	sessionName := "sessionName"
	accountId := "accountId"
	projectName := "projectName"
	oauthCode := "oauthCode"
	uuid := "uuid"
	envMock.ExpUuid = uuid
	gcpApiMock.ExpErrorOnGetOauth = true

	err := gcpIamUserAccountOauthSessionActions.CreateSession(sessionName, accountId, projectName, oauthCode)
	test.ExpectHttpError(t, err, http.StatusBadRequest, "error getting oauth token")
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{"GetOauthUrl(oauthCode)"}, []string{"GenerateUuid()"}, []string{}, []string{})
}

func TestCreateSession_KeychainSetSecretReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()

	sessionName := "sessionName"
	accountId := "accountId"
	projectName := "projectName"
	oauthCode := "oauthCode"
	uuid := "uuid"
	credentials := "credentialsJson"
	envMock.ExpUuid = uuid
	gcpApiMock.ExpOauthToken = oauth2.Token{}
	gcpApiMock.ExpCredentials = credentials
	gcpIamUserAccountOauthSessionActionsKeychainMock.ExpErrorOnSetSecret = true

	err := gcpIamUserAccountOauthSessionActions.CreateSession(sessionName, accountId, projectName, oauthCode)
	test.ExpectHttpError(t, err, http.StatusInternalServerError, "unable to set secret")
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{"GetOauthUrl(oauthCode)", "GetCredentials()"}, []string{"GenerateUuid()"}, []string{"SetSecret(credentialsJson, uuid-gcp-iam-user-account-oauth-session-credentials)"}, []string{})
}

func TestCreateSession_FacadeAddSessionReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()

	sessionName := "sessionName"
	accountId := "accountId"
	projectName := "projectName"
	oauthCode := "oauthCode"
	uuid := "uuid"
	credentials := "credentialsJson"
	envMock.ExpUuid = uuid
	gcpApiMock.ExpOauthToken = oauth2.Token{}
	gcpApiMock.ExpCredentials = credentials
	gcpIamUserAccountOauthSessionFacadeMock.ExpErrorOnAddSession = true

	err := gcpIamUserAccountOauthSessionActions.CreateSession(sessionName, accountId, projectName, oauthCode)
	test.ExpectHttpError(t, err, http.StatusConflict, "session already exist")
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{"GetOauthUrl(oauthCode)", "GetCredentials()"}, []string{"GenerateUuid()"}, []string{"SetSecret(credentialsJson, uuid-gcp-iam-user-account-oauth-session-credentials)"}, []string{"AddSession(sessionName)"})
}

func TestStartSession_NoPreviousActiveSession(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	gcpIamUserAccountOauthSessionFacadeMock.ExpGetSessions = []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Id: "ID2", Status: domain_gcp.NotActive}}
	envMock.ExpTime = "start-time"
	sessionId := "ID1"

	err := gcpIamUserAccountOauthSessionActions.StartSession(sessionId)
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GetTime()"}, []string{}, []string{"GetSessions()", "StartSession(ID1, start-time)"})
}

func TestStartSession_PreviousActiveSessionDiffersFromNewActiveSession(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	gcpIamUserAccountOauthSessionFacadeMock.ExpGetSessions = []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Id: "ID2", Status: domain_gcp.Active}}
	envMock.ExpTime = "start-time"
	sessionId := "ID1"

	err := gcpIamUserAccountOauthSessionActions.StartSession(sessionId)
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GetTime()"}, []string{}, []string{"GetSessions()", "StopSession(ID2, start-time)", "StartSession(ID1, start-time)"})
}

func TestStartSession_SessionWasAlreadyActive(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	gcpIamUserAccountOauthSessionFacadeMock.ExpGetSessions = []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Id: "ID1", Status: domain_gcp.Active}}
	envMock.ExpTime = "start-time"
	sessionId := "ID1"

	err := gcpIamUserAccountOauthSessionActions.StartSession(sessionId)
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GetTime()"}, []string{}, []string{"GetSessions()", "StartSession(ID1, start-time)"})
}

func TestStartSession_PreviousActiveSessionDifferentAndFacadeSetSessionStatusReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	gcpIamUserAccountOauthSessionFacadeMock.ExpGetSessions = []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Id: "ID2", Status: domain_gcp.Active}}
	gcpIamUserAccountOauthSessionFacadeMock.ExpErrorOnStopSession = true
	envMock.ExpTime = "start-time"
	sessionId := "ID1"

	err := gcpIamUserAccountOauthSessionActions.StartSession(sessionId)
	test.ExpectHttpError(t, err, http.StatusInternalServerError, "unable to stop the session")
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GetTime()"}, []string{}, []string{"GetSessions()", "StopSession(ID2, start-time)"})
}

func TestStartSession_FacadeSetSessionStatusReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	gcpIamUserAccountOauthSessionFacadeMock.ExpGetSessions = []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Id: "ID2", Status: domain_gcp.NotActive}}
	gcpIamUserAccountOauthSessionFacadeMock.ExpErrorOnStartSession = true
	envMock.ExpTime = "start-time"
	sessionId := "ID1"

	err := gcpIamUserAccountOauthSessionActions.StartSession(sessionId)
	test.ExpectHttpError(t, err, http.StatusInternalServerError, "unable to start the session")
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GetTime()"}, []string{}, []string{"GetSessions()", "StartSession(ID1, start-time)"})
}

func TestStopSession(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	envMock.ExpTime = "stop-time"
	sessionId := "ID"
	err := gcpIamUserAccountOauthSessionActions.StopSession(sessionId)
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GetTime()"}, []string{}, []string{"StopSession(ID, stop-time)"})
}

func TestStopSession_FacadeReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	gcpIamUserAccountOauthSessionFacadeMock.ExpErrorOnStopSession = true
	envMock.ExpTime = "stop-time"
	sessionId := "ID"
	err := gcpIamUserAccountOauthSessionActions.StopSession(sessionId)
	test.ExpectHttpError(t, err, http.StatusInternalServerError, "unable to stop the session")
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GetTime()"}, []string{}, []string{"StopSession(ID, stop-time)"})
}

func TestDeleteSession(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	sessionId := "ID"
	credentialsLabel := "credentialLabel"
	gcpIamUserAccountOauthSessionFacadeMock.ExpGetSessionById = gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{Id: "ID", CredentialsLabel: credentialsLabel}

	err := gcpIamUserAccountOauthSessionActions.DeleteSession(sessionId)
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{"DeleteSecret(credentialLabel)"}, []string{"GetSessionById(ID)", "RemoveSession(ID)"})
}

func TestDeleteSession_FacadeGetSessionByIdReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	sessionId := "ID"
	gcpIamUserAccountOauthSessionFacadeMock.ExpErrorOnGetSessionById = true

	err := gcpIamUserAccountOauthSessionActions.DeleteSession(sessionId)
	test.ExpectHttpError(t, err, http.StatusNotFound, "session not found")
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{}, []string{"GetSessionById(ID)"})
}

func TestDeleteSession_KeychainDeleteSecretReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	sessionId := "ID"
	credentialsLabel := "credentialLabel"
	gcpIamUserAccountOauthSessionFacadeMock.ExpGetSessionById = gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{Id: "ID", CredentialsLabel: credentialsLabel}
	gcpIamUserAccountOauthSessionActionsKeychainMock.ExpErrorOnDeleteSecret = true

	err := gcpIamUserAccountOauthSessionActions.DeleteSession(sessionId)
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{"DeleteSecret(credentialLabel)"}, []string{"GetSessionById(ID)", "RemoveSession(ID)"})
}

func TestDeleteSession_FacadeRemoveSessionReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	sessionId := "ID"
	credentialsLabel := "credentialLabel"
	gcpIamUserAccountOauthSessionFacadeMock.ExpGetSessionById = gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{Id: "ID", CredentialsLabel: credentialsLabel}
	gcpIamUserAccountOauthSessionFacadeMock.ExpErrorOnRemoveSession = true

	err := gcpIamUserAccountOauthSessionActions.DeleteSession(sessionId)
	test.ExpectHttpError(t, err, http.StatusNotFound, "session not found")
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{"DeleteSecret(credentialLabel)"}, []string{"GetSessionById(ID)", "RemoveSession(ID)"})
}

func TestEditSession(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()

	sessionId := "ID"
	sessionName := "sessionName"
	projectName := "projectName"

	err := gcpIamUserAccountOauthSessionActions.EditSession(sessionId, sessionName, projectName)
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{}, []string{"EditSession(ID, sessionName, projectName)"})
}

func TestEditSession_FacadeEditSessionReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()

	sessionId := "ID"
	sessionName := "sessionName"
	projectName := "projectName"
	gcpIamUserAccountOauthSessionFacadeMock.ExpErrorOnEditSession = true

	err := gcpIamUserAccountOauthSessionActions.EditSession(sessionId, sessionName, projectName)
	test.ExpectHttpError(t, err, http.StatusConflict, "unable to edit session, collision detected")

	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{}, []string{"EditSession(ID, sessionName, projectName)"})
}
