package use_case

import (
	"golang.org/x/oauth2"
	"leapp_daemon/domain/domain_gcp"
	"leapp_daemon/domain/domain_gcp/gcp_iam_user_account_oauth"
	"leapp_daemon/domain/domain_gcp/named_configuration"
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
	namedConfigurationsActionsMock                   mock.NamedConfigurationsActionsMock
	gcpIamUserAccountOauthSessionActions             *GcpIamUserAccountOauthSessionActions
)

func gcpIamUserAccountOauthSessionActionsSetup() {
	gcpApiMock = mock.NewGcpApiMock()
	envMock = mock.NewEnvironmentMock()
	gcpIamUserAccountOauthSessionActionsKeychainMock = mock.NewKeychainMock()
	gcpIamUserAccountOauthSessionFacadeMock = mock.NewGcpIamUserAccountOauthSessionsFacadeMock()
	namedConfigurationsActionsMock = mock.NewNamedConfigurationsActionsMock()
	gcpIamUserAccountOauthSessionActions = &GcpIamUserAccountOauthSessionActions{
		GcpApi:                              &gcpApiMock,
		Environment:                         &envMock,
		Keychain:                            &gcpIamUserAccountOauthSessionActionsKeychainMock,
		GcpIamUserAccountOauthSessionFacade: &gcpIamUserAccountOauthSessionFacadeMock,
		NamedConfigurationsActions:          &namedConfigurationsActionsMock,
	}
}

func gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t *testing.T, gcpApiMockCalls, envMockCalls,
	keychainMockCalls, facadeMockCalls, namedConfigurationsActionsMockCalls []string) {
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
	if !reflect.DeepEqual(namedConfigurationsActionsMock.GetCalls(), namedConfigurationsActionsMockCalls) {
		t.Fatalf("namedConfigurationsActionsMock expectation violation.\nMock calls: %v", namedConfigurationsActionsMock.GetCalls())
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
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{},
		[]string{"GetSessionById(ID)"}, []string{})
}

func TestGetSession_SessionFacadeReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	gcpIamUserAccountOauthSessionFacadeMock.ExpErrorOnGetSessionById = true

	_, err := gcpIamUserAccountOauthSessionActions.GetSession("ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session not found")
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{},
		[]string{"GetSessionById(ID)"}, []string{})
}

func TestGetOAuthUrl(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	gcpApiMock.ExpOauthUrl = "url"

	actualOauthUrl, err := gcpIamUserAccountOauthSessionActions.GetOAuthUrl()
	if err != nil && !reflect.DeepEqual("url", actualOauthUrl) {
		t.Fatalf("Returned unexpected oauth url")
	}
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{"GetOauthUrl()"}, []string{}, []string{},
		[]string{}, []string{})
}

func TestGetOAuthUrl_GcpApiReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	gcpApiMock.ExpErrorOnGetOauthUrl = true

	_, err := gcpIamUserAccountOauthSessionActions.GetOAuthUrl()
	test.ExpectHttpError(t, err, http.StatusNotFound, "error getting oauth url")
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{"GetOauthUrl()"}, []string{}, []string{},
		[]string{}, []string{})
}

func TestCreateSession(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()

	envMock.ExpUuid = "uuid"
	gcpApiMock.ExpOauthToken = oauth2.Token{}
	gcpApiMock.ExpCredentials = "credentialsJson"

	err := gcpIamUserAccountOauthSessionActions.CreateSession("sessionName", "accountId", "projectName", "configurationName", "oauthCode")
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{"GetOauthUrl(oauthCode)", "GetCredentials()"},
		[]string{"GenerateUuid()"}, []string{"SetSecret(credentialsJson, uuid-gcp-iam-user-account-oauth-session-credentials)"},
		[]string{"AddSession(sessionName)"}, []string{"GetOrCreateNamedConfiguration(configurationName)"})
}

func TestCreateSession_NamedConfigurationGetOrCreateReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()

	envMock.ExpUuid = "uuid"
	gcpApiMock.ExpOauthToken = oauth2.Token{}
	gcpApiMock.ExpCredentials = "credentialsJson"
	namedConfigurationsActionsMock.ExpErrorOnGetOrCreateNamedConfiguration = true

	err := gcpIamUserAccountOauthSessionActions.CreateSession("sessionName", "accountId", "projectName", "configurationName", "oauthCode")
	test.ExpectHttpError(t, err, http.StatusBadRequest, "configuration name is invalid")

	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GenerateUuid()"}, []string{},
		[]string{}, []string{"GetOrCreateNamedConfiguration(configurationName)"})
}

func TestCreateSession_GcpApiGetOauthTokenReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()

	envMock.ExpUuid = "uuid"
	gcpApiMock.ExpErrorOnGetOauth = true

	err := gcpIamUserAccountOauthSessionActions.CreateSession("sessionName", "accountId", "projectName", "configurationName", "oauthCode")
	test.ExpectHttpError(t, err, http.StatusBadRequest, "error getting oauth token")
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{"GetOauthUrl(oauthCode)"},
		[]string{"GenerateUuid()"}, []string{}, []string{}, []string{"GetOrCreateNamedConfiguration(configurationName)"})
}

func TestCreateSession_KeychainSetSecretReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()

	envMock.ExpUuid = "uuid"
	gcpApiMock.ExpOauthToken = oauth2.Token{}
	gcpApiMock.ExpCredentials = "credentialsJson"
	gcpIamUserAccountOauthSessionActionsKeychainMock.ExpErrorOnSetSecret = true

	err := gcpIamUserAccountOauthSessionActions.CreateSession("sessionName", "accountId", "projectName", "configurationName", "oauthCode")
	test.ExpectHttpError(t, err, http.StatusInternalServerError, "unable to set secret")
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{"GetOauthUrl(oauthCode)", "GetCredentials()"},
		[]string{"GenerateUuid()"}, []string{"SetSecret(credentialsJson, uuid-gcp-iam-user-account-oauth-session-credentials)"},
		[]string{}, []string{"GetOrCreateNamedConfiguration(configurationName)"})
}

func TestCreateSession_FacadeAddSessionReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()

	envMock.ExpUuid = "uuid"
	gcpApiMock.ExpOauthToken = oauth2.Token{}
	gcpApiMock.ExpCredentials = "credentialsJson"
	gcpIamUserAccountOauthSessionFacadeMock.ExpErrorOnAddSession = true

	err := gcpIamUserAccountOauthSessionActions.CreateSession("sessionName", "accountId", "projectName", "configurationName", "oauthCode")
	test.ExpectHttpError(t, err, http.StatusConflict, "session already exist")
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{"GetOauthUrl(oauthCode)", "GetCredentials()"},
		[]string{"GenerateUuid()"}, []string{"SetSecret(credentialsJson, uuid-gcp-iam-user-account-oauth-session-credentials)"},
		[]string{"AddSession(sessionName)"}, []string{"GetOrCreateNamedConfiguration(configurationName)"})
}

func TestStartSession_NoPreviousActiveSession(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	gcpIamUserAccountOauthSessionFacadeMock.ExpGetSessions = []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{
		{Id: "ID2", Status: domain_gcp.NotActive},
	}
	envMock.ExpTime = "start-time"
	err := gcpIamUserAccountOauthSessionActions.StartSession("ID1")
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GetTime()"}, []string{},
		[]string{"GetSessionById(ID1)", "GetSessions()", "StartSession(ID1, start-time)"}, []string{})
}

func TestStartSession_PreviousActiveSessionWithDifferentNamedConfiguration(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	gcpIamUserAccountOauthSessionFacadeMock.ExpGetSessionById = gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{
		Id: "ID1", Status: domain_gcp.Active, NamedConfigurationId: "NAMED1"}

	gcpIamUserAccountOauthSessionFacadeMock.ExpGetSessions = []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{
		{Id: "ID2", Status: domain_gcp.Active, NamedConfigurationId: "NAMED2"},
	}
	envMock.ExpTime = "start-time"
	err := gcpIamUserAccountOauthSessionActions.StartSession("ID1")
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GetTime()"}, []string{},
		[]string{"GetSessionById(ID1)", "GetSessions()", "StartSession(ID1, start-time)"}, []string{})
}

func TestStartSession_PreviousActiveSessionWithSameNamedConfiguration(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	gcpIamUserAccountOauthSessionFacadeMock.ExpGetSessionById = gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{
		Id: "ID1", Status: domain_gcp.Active, NamedConfigurationId: "NAMED1"}

	gcpIamUserAccountOauthSessionFacadeMock.ExpGetSessions = []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{
		{Id: "ID2", Status: domain_gcp.Active, NamedConfigurationId: "NAMED1"},
	}
	envMock.ExpTime = "start-time"
	err := gcpIamUserAccountOauthSessionActions.StartSession("ID1")
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GetTime()", "GetTime()"}, []string{},
		[]string{"GetSessionById(ID1)", "GetSessions()", "StopSession(ID2, start-time)", "StartSession(ID1, start-time)"}, []string{})
}

func TestStartSession_GetSessionByIdReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	gcpIamUserAccountOauthSessionFacadeMock.ExpErrorOnGetSessionById = true

	err := gcpIamUserAccountOauthSessionActions.StartSession("ID1")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session not found")
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{},
		[]string{"GetSessionById(ID1)"}, []string{})
}

func TestStartSession_StopSessionReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	gcpIamUserAccountOauthSessionFacadeMock.ExpGetSessionById = gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{
		Id: "ID1", Status: domain_gcp.Active, NamedConfigurationId: "NAMED1"}

	gcpIamUserAccountOauthSessionFacadeMock.ExpGetSessions = []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{
		{Id: "ID2", Status: domain_gcp.Active, NamedConfigurationId: "NAMED1"},
	}
	gcpIamUserAccountOauthSessionFacadeMock.ExpErrorOnStopSession = true
	envMock.ExpTime = "start-time"

	err := gcpIamUserAccountOauthSessionActions.StartSession("ID1")
	test.ExpectHttpError(t, err, http.StatusInternalServerError, "unable to stop the session")
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GetTime()"}, []string{},
		[]string{"GetSessionById(ID1)", "GetSessions()", "StopSession(ID2, start-time)"}, []string{})
}

func TestStartSession_StartSessionReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	gcpIamUserAccountOauthSessionFacadeMock.ExpGetSessionById = gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{
		Id: "ID1", Status: domain_gcp.Active, NamedConfigurationId: "NAMED1"}
	gcpIamUserAccountOauthSessionFacadeMock.ExpErrorOnStartSession = true
	envMock.ExpTime = "start-time"

	err := gcpIamUserAccountOauthSessionActions.StartSession("ID1")
	test.ExpectHttpError(t, err, http.StatusInternalServerError, "unable to start the session")
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GetTime()"}, []string{},
		[]string{"GetSessionById(ID1)", "GetSessions()", "StartSession(ID1, start-time)"}, []string{})
}

func TestStopSession(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	envMock.ExpTime = "stop-time"
	err := gcpIamUserAccountOauthSessionActions.StopSession("ID")
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GetTime()"}, []string{},
		[]string{"StopSession(ID, stop-time)"}, []string{})
}

func TestStopSession_FacadeReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	gcpIamUserAccountOauthSessionFacadeMock.ExpErrorOnStopSession = true
	envMock.ExpTime = "stop-time"
	err := gcpIamUserAccountOauthSessionActions.StopSession("ID")
	test.ExpectHttpError(t, err, http.StatusInternalServerError, "unable to stop the session")
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GetTime()"}, []string{},
		[]string{"StopSession(ID, stop-time)"}, []string{})
}

func TestDeleteSession(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	gcpIamUserAccountOauthSessionFacadeMock.ExpGetSessionById = gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{
		Id:               "ID",
		CredentialsLabel: "credentialLabel",
	}

	err := gcpIamUserAccountOauthSessionActions.DeleteSession("ID")
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{},
		[]string{"DeleteSecret(credentialLabel)"}, []string{"GetSessionById(ID)", "RemoveSession(ID)"}, []string{})
}

func TestDeleteSession_FacadeGetSessionByIdReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	gcpIamUserAccountOauthSessionFacadeMock.ExpErrorOnGetSessionById = true

	err := gcpIamUserAccountOauthSessionActions.DeleteSession("ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session not found")
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{},
		[]string{"GetSessionById(ID)"}, []string{})
}

func TestDeleteSession_KeychainDeleteSecretReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	gcpIamUserAccountOauthSessionFacadeMock.ExpGetSessionById = gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{
		Id:               "ID",
		CredentialsLabel: "credentialLabel",
	}
	gcpIamUserAccountOauthSessionActionsKeychainMock.ExpErrorOnDeleteSecret = true

	err := gcpIamUserAccountOauthSessionActions.DeleteSession("ID")
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{},
		[]string{"DeleteSecret(credentialLabel)"}, []string{"GetSessionById(ID)", "RemoveSession(ID)"}, []string{})
}

func TestDeleteSession_FacadeRemoveSessionReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()
	gcpIamUserAccountOauthSessionFacadeMock.ExpGetSessionById = gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{
		Id:               "ID",
		CredentialsLabel: "credentialLabel",
	}
	gcpIamUserAccountOauthSessionFacadeMock.ExpErrorOnRemoveSession = true

	err := gcpIamUserAccountOauthSessionActions.DeleteSession("ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session not found")
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{},
		[]string{"DeleteSecret(credentialLabel)"}, []string{"GetSessionById(ID)", "RemoveSession(ID)"}, []string{})
}

func TestEditSession(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()

	configurationName := "configurationName"
	namedConfigurationsActionsMock.ExpNamedConfiguration = named_configuration.NamedConfiguration{
		Id:   "ConfigurationID",
		Name: configurationName,
	}

	err := gcpIamUserAccountOauthSessionActions.EditSession("ID", "sessionName", "projectName", configurationName)
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{},
		[]string{"GetSessionById(ID)", "EditSession(ID, sessionName, projectName, ConfigurationID)"},
		[]string{"GetOrCreateNamedConfiguration(configurationName)"})
}

func TestEditSession_StopPreviousActiveSessionWithSameNamedConfiguration(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()

	configurationId := "config-id-1"
	configurationName := "config-name"
	namedConfigurationsActionsMock.ExpNamedConfiguration = named_configuration.NamedConfiguration{
		Id:   configurationId,
		Name: configurationName,
	}

	gcpIamUserAccountOauthSessionFacadeMock.ExpGetSessionById =
		gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{Id: "session-id1", Status: domain_gcp.Active}

	gcpIamUserAccountOauthSessionFacadeMock.ExpGetSessions =
		[]gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{
			{Id: "session-id2", Status: domain_gcp.NotActive},
			{Id: "session-id3", Status: domain_gcp.Active, NamedConfigurationId: "config-id-2"},
			{Id: "session-id4", Status: domain_gcp.Active, NamedConfigurationId: configurationId},
		}

	envMock.ExpTime = "stop-time"

	err := gcpIamUserAccountOauthSessionActions.EditSession("session-id1", "sessionName", "projectName", configurationName)
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GetTime()"}, []string{},
		[]string{"GetSessionById(session-id1)", "GetSessions()", "StopSession(session-id4, stop-time)",
			"EditSession(session-id1, sessionName, projectName, config-id-1)"},
		[]string{"GetOrCreateNamedConfiguration(config-name)"})
}

func TestEditSession_NamedConfigurationsActionsGetOrCreateReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()

	namedConfigurationsActionsMock.ExpErrorOnGetOrCreateNamedConfiguration = true

	err := gcpIamUserAccountOauthSessionActions.EditSession("ID", "sessionName", "projectName", "configurationName")
	test.ExpectHttpError(t, err, http.StatusBadRequest, "configuration name is invalid")

	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{}, []string{},
		[]string{"GetOrCreateNamedConfiguration(configurationName)"})
}

func TestEditSession_GetSessionByIdReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()

	gcpIamUserAccountOauthSessionFacadeMock.ExpErrorOnGetSessionById = true

	err := gcpIamUserAccountOauthSessionActions.EditSession("ID", "sessionName", "projectName", "configurationName")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session not found")

	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{},
		[]string{"GetSessionById(ID)"}, []string{"GetOrCreateNamedConfiguration(configurationName)"})
}

func TestEditSession_StopPreviousActiveSessionReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()

	configurationId := "config-id-1"
	configurationName := "config-name"
	namedConfigurationsActionsMock.ExpNamedConfiguration = named_configuration.NamedConfiguration{
		Id:   configurationId,
		Name: configurationName,
	}

	gcpIamUserAccountOauthSessionFacadeMock.ExpGetSessionById =
		gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{Id: "session-id1", Status: domain_gcp.Active}

	gcpIamUserAccountOauthSessionFacadeMock.ExpGetSessions =
		[]gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{
			{Id: "session-id4", Status: domain_gcp.Active, NamedConfigurationId: configurationId},
		}

	gcpIamUserAccountOauthSessionFacadeMock.ExpErrorOnStopSession = true

	envMock.ExpTime = "stop-time"

	err := gcpIamUserAccountOauthSessionActions.EditSession("session-id1", "sessionName", "projectName", configurationName)
	test.ExpectHttpError(t, err, http.StatusInternalServerError, "unable to stop the session")
	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GetTime()"}, []string{},
		[]string{"GetSessionById(session-id1)", "GetSessions()", "StopSession(session-id4, stop-time)"},
		[]string{"GetOrCreateNamedConfiguration(config-name)"})
}

func TestEditSession_FacadeEditSessionReturnsError(t *testing.T) {
	gcpIamUserAccountOauthSessionActionsSetup()

	configurationName := "configurationName"
	gcpIamUserAccountOauthSessionFacadeMock.ExpErrorOnEditSession = true
	namedConfigurationsActionsMock.ExpNamedConfiguration = named_configuration.NamedConfiguration{
		Id:   "ConfigurationID",
		Name: configurationName,
	}

	err := gcpIamUserAccountOauthSessionActions.EditSession("ID", "sessionName", "projectName", configurationName)
	test.ExpectHttpError(t, err, http.StatusConflict, "unable to edit session, collision detected")

	gcpIamUserAccountOauthSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{},
		[]string{"GetSessionById(ID)", "EditSession(ID, sessionName, projectName, ConfigurationID)"},
		[]string{"GetOrCreateNamedConfiguration(configurationName)"})
}
