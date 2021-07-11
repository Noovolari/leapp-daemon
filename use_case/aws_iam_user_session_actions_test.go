package use_case

import (
	"github.com/aws/aws-sdk-go/service/sts"
	"leapp_daemon/domain/domain_aws"
	"leapp_daemon/domain/domain_aws/aws_iam_user"
	"leapp_daemon/domain/domain_aws/named_profile"
	"leapp_daemon/test"
	"leapp_daemon/test/mock"
	"net/http"
	"reflect"
	"testing"
	"time"
)

var (
	stsApiMock                                mock.StsApiMock
	awsIamUserSessionActionsEnvMock           mock.EnvironmentMock
	awsIamUserSessionActionsKeychainMock      mock.KeychainMock
	awsIamUserSessionActionsFacadeMock        mock.AwsIamUserSessionsFacadeMock
	awsIamUserSessionNamedProfilesActionsMock mock.NamedProfilesActionsMock
	awsIamUserSessionActions                  *AwsIamUserSessionActions
)

func awsIamUserSessionActionsSetup() {
	stsApiMock = mock.NewStsApiMock()
	awsIamUserSessionActionsEnvMock = mock.NewEnvironmentMock()
	awsIamUserSessionActionsKeychainMock = mock.NewKeychainMock()
	awsIamUserSessionActionsFacadeMock = mock.NewAwsIamUserSessionsFacadeMock()
	awsIamUserSessionNamedProfilesActionsMock = mock.NewNamedProfilesActionsMock()

	awsIamUserSessionActions = &AwsIamUserSessionActions{
		Environment:              &awsIamUserSessionActionsEnvMock,
		Keychain:                 &awsIamUserSessionActionsKeychainMock,
		StsApi:                   &stsApiMock,
		NamedProfilesActions:     &awsIamUserSessionNamedProfilesActionsMock,
		AwsIamUserSessionsFacade: &awsIamUserSessionActionsFacadeMock,
	}
}

func awsIamUserSessionActionsVerifyExpectedCalls(t *testing.T, stsApiMockCalls, envMockCalls, keychainMockCalls,
	facadeMockCalls []string, namedProfileActionsMockCalls []string) {
	if !reflect.DeepEqual(stsApiMock.GetCalls(), stsApiMockCalls) {
		t.Fatalf("stsApiMock expectation violation.\nMock calls: %v", stsApiMock.GetCalls())
	}
	if !reflect.DeepEqual(awsIamUserSessionActionsEnvMock.GetCalls(), envMockCalls) {
		t.Fatalf("envMock expectation violation.\nMock calls: %v", awsIamUserSessionActionsEnvMock.GetCalls())
	}
	if !reflect.DeepEqual(awsIamUserSessionActionsKeychainMock.GetCalls(), keychainMockCalls) {
		t.Fatalf("keychainMock expectation violation.\nMock calls: %v", awsIamUserSessionActionsKeychainMock.GetCalls())
	}
	if !reflect.DeepEqual(awsIamUserSessionActionsFacadeMock.GetCalls(), facadeMockCalls) {
		t.Fatalf("facadeMock expectation violation.\nMock calls: %v", awsIamUserSessionActionsFacadeMock.GetCalls())
	}
	if !reflect.DeepEqual(awsIamUserSessionNamedProfilesActionsMock.GetCalls(), namedProfileActionsMockCalls) {
		t.Fatalf("facadeMock expectation violation.\nMock calls: %v", awsIamUserSessionNamedProfilesActionsMock.GetCalls())
	}
}

func TestAwsIamUserSessionActions_GetSession(t *testing.T) {
	awsIamUserSessionActionsSetup()

	session := aws_iam_user.AwsIamUserSession{Name: "test_session"}
	awsIamUserSessionActionsFacadeMock.ExpGetSessionById = session

	actualSession, err := awsIamUserSessionActions.GetSession("ID")
	if err != nil && !reflect.DeepEqual(session, actualSession) {
		t.Fatalf("Returned unexpected session")
	}
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{},
		[]string{"GetSessionById(ID)"}, []string{})
}

func TestAwsIamUserSessionActions_GetSession_SessionFacadeReturnsError(t *testing.T) {
	awsIamUserSessionActionsSetup()
	awsIamUserSessionActionsFacadeMock.ExpErrorOnGetSessionById = true

	_, err := awsIamUserSessionActions.GetSession("ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session not found")
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{}, []string{"GetSessionById(ID)"}, []string{})
}

func TestAwsIamUserSessionActions_CreateSession(t *testing.T) {
	awsIamUserSessionActionsSetup()

	sessionName := "sessionName"
	region := "region"
	accountNumber := "accountNumber"
	userName := "userName"
	accessKeyId := "accessKeyId"
	secretKey := "secretKey"
	mfaDevice := "mfaDevice"
	profileName := "profileName"

	awsIamUserSessionActionsEnvMock.ExpUuid = "uuid"
	awsIamUserSessionNamedProfilesActionsMock.ExpNamedProfile = named_profile.NamedProfile{Id: "ID", Name: profileName}

	err := awsIamUserSessionActions.CreateSession(sessionName, region, accountNumber, userName, accessKeyId, secretKey,
		mfaDevice, profileName)
	if err != nil {
		t.Fatalf("Unexpected error")
	}

	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{},
		[]string{"GenerateUuid()"},
		[]string{"SetSecret(accessKeyId, uuid-aws-iam-user-session-access-key-id)", "SetSecret(secretKey, uuid-aws-iam-user-session-secret-key)"},
		[]string{"AddSession(sessionName)"},
		[]string{"GetOrCreateNamedProfile(profileName)"})
}

func TestAwsIamUserSessionActions_CreateSession_NamedProfileActionsReturnsError(t *testing.T) {
	awsIamUserSessionActionsSetup()

	sessionName := "sessionName"
	region := "region"
	accountNumber := "accountNumber"
	userName := "userName"
	accessKeyId := "accessKeyId"
	secretKey := "secretKey"
	mfaDevice := "mfaDevice"
	profileName := "profileName"

	awsIamUserSessionActionsEnvMock.ExpUuid = "uuid"
	awsIamUserSessionNamedProfilesActionsMock.ExpErrorOnGetOrCreateNamedProfile = true

	err := awsIamUserSessionActions.CreateSession(sessionName, region, accountNumber, userName, accessKeyId, secretKey,
		mfaDevice, profileName)
	test.ExpectHttpError(t, err, http.StatusNotFound, "named profile not found")
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GenerateUuid()"}, []string{}, []string{},
		[]string{"GetOrCreateNamedProfile(profileName)"})
}

func TestAwsIamUserSessionActions_CreateSession_KeychainSetSecretReturnsError(t *testing.T) {
	awsIamUserSessionActionsSetup()

	sessionName := "sessionName"
	region := "region"
	accountNumber := "accountNumber"
	userName := "userName"
	accessKeyId := "accessKeyId"
	secretKey := "secretKey"
	mfaDevice := "mfaDevice"
	profileName := "profileName"

	awsIamUserSessionActionsEnvMock.ExpUuid = "uuid"
	awsIamUserSessionNamedProfilesActionsMock.ExpNamedProfile = named_profile.NamedProfile{Id: "ID", Name: profileName}
	awsIamUserSessionActionsKeychainMock.ExpErrorOnSetSecret = true

	err := awsIamUserSessionActions.CreateSession(sessionName, region, accountNumber, userName, accessKeyId, secretKey,
		mfaDevice, profileName)
	test.ExpectHttpError(t, err, http.StatusInternalServerError, "unable to set secret")
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GenerateUuid()"},
		[]string{"SetSecret(accessKeyId, uuid-aws-iam-user-session-access-key-id)"}, []string{},
		[]string{"GetOrCreateNamedProfile(profileName)"})
}

func TestAwsIamUserSessionActions_CreateSession_FacadeAddSessionReturnsError(t *testing.T) {
	awsIamUserSessionActionsSetup()

	sessionName := "sessionName"
	region := "region"
	accountNumber := "accountNumber"
	userName := "userName"
	accessKeyId := "accessKeyId"
	secretKey := "secretKey"
	mfaDevice := "mfaDevice"
	profileName := "profileName"

	awsIamUserSessionActionsEnvMock.ExpUuid = "uuid"
	awsIamUserSessionNamedProfilesActionsMock.ExpNamedProfile = named_profile.NamedProfile{Id: "ID", Name: profileName}
	awsIamUserSessionActionsFacadeMock.ExpErrorOnAddSession = true

	err := awsIamUserSessionActions.CreateSession(sessionName, region, accountNumber, userName, accessKeyId, secretKey,
		mfaDevice, profileName)
	test.ExpectHttpError(t, err, http.StatusConflict, "session already exist")
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GenerateUuid()"},
		[]string{"SetSecret(accessKeyId, uuid-aws-iam-user-session-access-key-id)", "SetSecret(secretKey, uuid-aws-iam-user-session-secret-key)"},
		[]string{"AddSession(sessionName)"}, []string{"GetOrCreateNamedProfile(profileName)"})
}

func TestAwsIamUserSessionActions_StartSession(t *testing.T) {
	awsIamUserSessionActionsSetup()
	sessionId := "ID1"
	awsIamUserSessionActionsFacadeMock.ExpGetSessionById = aws_iam_user.AwsIamUserSession{
		ID:                     sessionId,
		Status:                 domain_aws.NotActive,
		SessionTokenLabel:      "sessionTokenLabel",
		SessionTokenExpiration: "2020-01-01T12:00:00Z",
	}
	awsIamUserSessionActionsEnvMock.ExpTime = "2020-01-01T11:00:00Z"
	awsIamUserSessionActionsKeychainMock.ExpSecretExist = true
	err := awsIamUserSessionActions.StartSession(sessionId)
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GetTime()"}, []string{"DoesSecretExist(sessionTokenLabel)"},
		[]string{"GetSessionById(ID1)", "StartingSession(ID1)", "GetSessions()", "StartSession(ID1, 2020-01-01T11:00:00Z)"}, []string{})
}

func TestAwsIamUserSessionActions_StartSession_PreviousActiveSessionWithTheSameNamedProfile(t *testing.T) {
	awsIamUserSessionActionsSetup()
	sessionId := "ID1"
	awsIamUserSessionActionsFacadeMock.ExpGetSessionById = aws_iam_user.AwsIamUserSession{
		ID:                     sessionId,
		Status:                 domain_aws.NotActive,
		SessionTokenLabel:      "sessionTokenLabel",
		SessionTokenExpiration: "2020-01-01T12:00:00Z",
		NamedProfileID:         "ProfileId",
	}
	awsIamUserSessionActionsFacadeMock.ExpGetSessions = []aws_iam_user.AwsIamUserSession{{
		ID:             "ID2",
		Status:         domain_aws.Active,
		NamedProfileID: "ProfileId",
	}}
	awsIamUserSessionActionsEnvMock.ExpTime = "2020-01-01T11:00:00Z"
	awsIamUserSessionActionsKeychainMock.ExpSecretExist = true
	err := awsIamUserSessionActions.StartSession(sessionId)
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GetTime()", "GetTime()"},
		[]string{"DoesSecretExist(sessionTokenLabel)"},
		[]string{"GetSessionById(ID1)", "StartingSession(ID1)", "GetSessions()", "StopSession(ID2, 2020-01-01T11:00:00Z)",
			"StartSession(ID1, 2020-01-01T11:00:00Z)"}, []string{})
}

func TestAwsIamUserSessionActions_StartSession_ErrorStoppingSessionWithSameNamedProfile(t *testing.T) {
	awsIamUserSessionActionsSetup()
	sessionId := "ID1"
	awsIamUserSessionActionsFacadeMock.ExpGetSessionById = aws_iam_user.AwsIamUserSession{
		ID:                     sessionId,
		Status:                 domain_aws.NotActive,
		SessionTokenLabel:      "sessionTokenLabel",
		SessionTokenExpiration: "2020-01-01T12:00:00Z",
		NamedProfileID:         "ProfileId",
	}
	awsIamUserSessionActionsFacadeMock.ExpGetSessions = []aws_iam_user.AwsIamUserSession{{
		ID:             "ID2",
		Status:         domain_aws.Active,
		NamedProfileID: "ProfileId",
	}}
	awsIamUserSessionActionsFacadeMock.ExpErrorOnStopSession = true
	awsIamUserSessionActionsEnvMock.ExpTime = "2020-01-01T11:00:00Z"
	awsIamUserSessionActionsKeychainMock.ExpSecretExist = true
	err := awsIamUserSessionActions.StartSession(sessionId)
	test.ExpectHttpError(t, err, http.StatusInternalServerError, "unable to stop the session")
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GetTime()", "GetTime()"},
		[]string{"DoesSecretExist(sessionTokenLabel)"},
		[]string{"GetSessionById(ID1)", "StartingSession(ID1)", "GetSessions()", "StopSession(ID2, 2020-01-01T11:00:00Z)",
			"StopSession(ID1, 2020-01-01T11:00:00Z)"}, []string{})
}

func TestAwsIamUserSessionActions_StartSession_SessionNotFound(t *testing.T) {
	awsIamUserSessionActionsSetup()
	sessionId := "ID1"
	awsIamUserSessionActionsFacadeMock.ExpErrorOnGetSessionById = true

	err := awsIamUserSessionActions.StartSession(sessionId)
	test.ExpectHttpError(t, err, http.StatusNotFound, "session not found")

	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{},
		[]string{"GetSessionById(ID1)"}, []string{})
}

func TestAwsIamUserSessionActions_StartSession_FacadeStartingSessionReturnsError(t *testing.T) {
	awsIamUserSessionActionsSetup()
	sessionId := "ID1"
	awsIamUserSessionActionsFacadeMock.ExpGetSessionById = aws_iam_user.AwsIamUserSession{
		ID:                     sessionId,
		Status:                 domain_aws.NotActive,
		SessionTokenLabel:      "sessionTokenLabel",
		SessionTokenExpiration: "2020-01-01T12:00:00Z",
	}
	awsIamUserSessionActionsFacadeMock.ExpErrorOnStartingSession = true
	err := awsIamUserSessionActions.StartSession(sessionId)
	test.ExpectHttpError(t, err, http.StatusInternalServerError, "starting session failed")

	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{},
		[]string{"GetSessionById(ID1)", "StartingSession(ID1)"}, []string{})
}

func TestAwsIamUserSessionActions_StartSession_RefreshSessionTokenReturnsError(t *testing.T) {
	awsIamUserSessionActionsSetup()
	sessionId := "ID1"
	awsIamUserSessionActionsFacadeMock.ExpGetSessionById = aws_iam_user.AwsIamUserSession{
		ID:                     sessionId,
		Status:                 domain_aws.NotActive,
		SessionTokenLabel:      "sessionTokenLabel",
		SessionTokenExpiration: "2020-01-01T12:00:00Z",
		AccessKeyIDLabel:       "accessKeyId1",
	}
	awsIamUserSessionActionsEnvMock.ExpTime = "2020-01-01T13:00:00Z"
	awsIamUserSessionActionsKeychainMock.ExpSecretExist = true
	awsIamUserSessionActionsKeychainMock.ExpErrorOnGetSecret = true
	err := awsIamUserSessionActions.StartSession(sessionId)
	test.ExpectHttpError(t, err, http.StatusInternalServerError, "unable to get secret")

	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GetTime()"},
		[]string{"DoesSecretExist(sessionTokenLabel)", "GetSecret(accessKeyId1)"},
		[]string{"GetSessionById(ID1)", "StartingSession(ID1)", "StopSession(ID1, 2020-01-01T13:00:00Z)"}, []string{})
}

func TestAwsIamUserSessionActions_StartSession_FacadeStartSessionReturnsError(t *testing.T) {
	awsIamUserSessionActionsSetup()
	sessionId := "ID1"
	awsIamUserSessionActionsFacadeMock.ExpGetSessionById = aws_iam_user.AwsIamUserSession{
		ID:                     sessionId,
		Status:                 domain_aws.NotActive,
		SessionTokenLabel:      "sessionTokenLabel",
		SessionTokenExpiration: "2020-01-01T12:00:00Z",
	}
	awsIamUserSessionActionsFacadeMock.ExpErrorOnStartSession = true
	awsIamUserSessionActionsEnvMock.ExpTime = "2020-01-01T11:00:00Z"
	awsIamUserSessionActionsKeychainMock.ExpSecretExist = true
	err := awsIamUserSessionActions.StartSession(sessionId)
	test.ExpectHttpError(t, err, http.StatusInternalServerError, "unable to start the session")

	awsIamUserSessionActionsVerifyExpectedCalls(t,
		[]string{},
		[]string{"GetTime()"},
		[]string{"DoesSecretExist(sessionTokenLabel)"},
		[]string{"GetSessionById(ID1)", "StartingSession(ID1)", "GetSessions()", "StartSession(ID1, 2020-01-01T11:00:00Z)",
			"StopSession(ID1, 2020-01-01T11:00:00Z)"},
		[]string{})
}

func TestAwsIamUserSessionActions_StopSession(t *testing.T) {
	awsIamUserSessionActionsSetup()
	awsIamUserSessionActionsEnvMock.ExpTime = "stop-time"

	err := awsIamUserSessionActions.StopSession("ID")
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GetTime()"}, []string{},
		[]string{"StopSession(ID, stop-time)"}, []string{})
}

func TestAwsIamUserSessionActions_StopSession_FacadeReturnsError(t *testing.T) {
	awsIamUserSessionActionsSetup()
	awsIamUserSessionActionsFacadeMock.ExpErrorOnStopSession = true
	awsIamUserSessionActionsEnvMock.ExpTime = "stop-time"

	err := awsIamUserSessionActions.StopSession("ID")
	test.ExpectHttpError(t, err, http.StatusInternalServerError, "unable to stop the session")
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GetTime()"}, []string{},
		[]string{"StopSession(ID, stop-time)"}, []string{})
}

func TestAwsIamUserSessionActions_DeleteSession(t *testing.T) {
	awsIamUserSessionActionsSetup()
	sessionId := "ID"
	awsIamUserSessionActionsFacadeMock.ExpGetSessionById = aws_iam_user.AwsIamUserSession{
		ID:                "ID",
		AccessKeyIDLabel:  "accessKeyIdLabel",
		SecretKeyLabel:    "secretKeyLabel",
		SessionTokenLabel: "sessionTokenLabel",
	}

	err := awsIamUserSessionActions.DeleteSession(sessionId)
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{},
		[]string{"DeleteSecret(accessKeyIdLabel)", "DeleteSecret(secretKeyLabel)", "DeleteSecret(sessionTokenLabel)"},
		[]string{"GetSessionById(ID)", "RemoveSession(ID)"}, []string{})
}

func TestAwsIamUserSessionActions_DeleteSession_FacadeGetSessionByIdReturnsError(t *testing.T) {
	awsIamUserSessionActionsSetup()
	awsIamUserSessionActionsFacadeMock.ExpErrorOnGetSessionById = true

	err := awsIamUserSessionActions.DeleteSession("ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session not found")
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{}, []string{"GetSessionById(ID)"},
		[]string{})
}

func TestAwsIamUserSessionActions_DeleteSession_KeychainDeleteSecretReturnsError(t *testing.T) {
	awsIamUserSessionActionsSetup()
	awsIamUserSessionActionsFacadeMock.ExpGetSessionById = aws_iam_user.AwsIamUserSession{
		ID:                "ID",
		AccessKeyIDLabel:  "accessKeyIdLabel",
		SecretKeyLabel:    "secretKeyLabel",
		SessionTokenLabel: "sessionTokenLabel",
	}

	awsIamUserSessionActionsKeychainMock.ExpErrorOnDeleteSecret = true

	err := awsIamUserSessionActions.DeleteSession("ID")
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{},
		[]string{"DeleteSecret(accessKeyIdLabel)", "DeleteSecret(secretKeyLabel)", "DeleteSecret(sessionTokenLabel)"},
		[]string{"GetSessionById(ID)", "RemoveSession(ID)"}, []string{})
}

func TestAwsIamUserSessionActions_DeleteSession_FacadeRemoveSessionReturnsError(t *testing.T) {
	awsIamUserSessionActionsSetup()
	awsIamUserSessionActionsFacadeMock.ExpGetSessionById = aws_iam_user.AwsIamUserSession{
		ID:                "ID",
		AccessKeyIDLabel:  "accessKeyIdLabel",
		SecretKeyLabel:    "secretKeyLabel",
		SessionTokenLabel: "sessionTokenLabel",
	}
	awsIamUserSessionActionsFacadeMock.ExpErrorOnRemoveSession = true

	err := awsIamUserSessionActions.DeleteSession("ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session not found")
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{},
		[]string{"DeleteSecret(accessKeyIdLabel)", "DeleteSecret(secretKeyLabel)", "DeleteSecret(sessionTokenLabel)"},
		[]string{"GetSessionById(ID)", "RemoveSession(ID)"}, []string{})
}

func TestAwsIamUserSessionActions_EditSession(t *testing.T) {
	awsIamUserSessionActionsSetup()
	awsIamUserSessionActionsFacadeMock.ExpGetSessionById = aws_iam_user.AwsIamUserSession{
		ID:                "ID",
		AccessKeyIDLabel:  "accessKeyIdLabel",
		SecretKeyLabel:    "secretKeyLabel",
		SessionTokenLabel: "sessionTokenLabel",
		Status:            domain_aws.NotActive,
	}
	awsIamUserSessionNamedProfilesActionsMock.ExpNamedProfile = named_profile.NamedProfile{Id: "ProfileId", Name: "profileName"}

	err := awsIamUserSessionActions.EditSession("ID", "sessionName", "region",
		"accountNumber", "userName", "awsAccessKeyId", "awsSecretKey",
		"mfaDevice", "profileName")
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{},
		[]string{"SetSecret(awsAccessKeyId, accessKeyIdLabel)", "SetSecret(awsSecretKey, secretKeyLabel)"},
		[]string{"GetSessionById(ID)", "EditSession(ID, sessionName, region, accountNumber, userName, mfaDevice, ProfileId)"},
		[]string{"GetOrCreateNamedProfile(profileName)"})
}

func TestAwsIamUserSessionActions_EditSession_PreviousActiveSessionWithTheSameNamedProfile(t *testing.T) {
	awsIamUserSessionActionsSetup()

	awsIamUserSessionActionsFacadeMock.ExpGetSessionById = aws_iam_user.AwsIamUserSession{
		ID:                "ID",
		AccessKeyIDLabel:  "accessKeyIdLabel",
		SecretKeyLabel:    "secretKeyLabel",
		SessionTokenLabel: "sessionTokenLabel",
		Status:            domain_aws.Active,
		NamedProfileID:    "ProfileId",
	}
	awsIamUserSessionActionsFacadeMock.ExpGetSessions = []aws_iam_user.AwsIamUserSession{{
		ID:             "ID2",
		Status:         domain_aws.Active,
		NamedProfileID: "ProfileId",
	}}
	awsIamUserSessionActionsEnvMock.ExpTime = "stop-time"
	awsIamUserSessionNamedProfilesActionsMock.ExpNamedProfile = named_profile.NamedProfile{Id: "ProfileId", Name: "profileName"}

	err := awsIamUserSessionActions.EditSession("ID", "sessionName", "region",
		"accountNumber", "userName", "awsAccessKeyId", "awsSecretKey",
		"mfaDevice", "profileName")
	if err != nil {
		t.Fatalf("Unexpected error")
	}
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GetTime()"},
		[]string{"SetSecret(awsAccessKeyId, accessKeyIdLabel)", "SetSecret(awsSecretKey, secretKeyLabel)"},
		[]string{"GetSessionById(ID)", "GetSessions()", "StopSession(ID2, stop-time)", "EditSession(ID, sessionName, region, accountNumber, userName, mfaDevice, ProfileId)"},
		[]string{"GetOrCreateNamedProfile(profileName)"})
}

func TestAwsIamUserSessionActions_EditSession_ErrorOnStoppingPreviousActiveSessionWithTheSameNamedProfile(t *testing.T) {
	awsIamUserSessionActionsSetup()

	awsIamUserSessionActionsFacadeMock.ExpGetSessionById = aws_iam_user.AwsIamUserSession{
		ID:                "ID",
		AccessKeyIDLabel:  "accessKeyIdLabel",
		SecretKeyLabel:    "secretKeyLabel",
		SessionTokenLabel: "sessionTokenLabel",
		Status:            domain_aws.Active,
		NamedProfileID:    "ProfileId",
	}
	awsIamUserSessionActionsFacadeMock.ExpGetSessions = []aws_iam_user.AwsIamUserSession{{
		ID:             "ID2",
		Status:         domain_aws.Active,
		NamedProfileID: "ProfileId",
	}}
	awsIamUserSessionActionsFacadeMock.ExpErrorOnStopSession = true
	awsIamUserSessionActionsEnvMock.ExpTime = "stop-time"
	awsIamUserSessionNamedProfilesActionsMock.ExpNamedProfile = named_profile.NamedProfile{Id: "ProfileId", Name: "profileName"}

	err := awsIamUserSessionActions.EditSession("ID", "sessionName", "region",
		"accountNumber", "userName", "awsAccessKeyId", "awsSecretKey",
		"mfaDevice", "profileName")
	test.ExpectHttpError(t, err, http.StatusInternalServerError, "unable to stop the session")
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GetTime()"},
		[]string{"SetSecret(awsAccessKeyId, accessKeyIdLabel)", "SetSecret(awsSecretKey, secretKeyLabel)"},
		[]string{"GetSessionById(ID)", "GetSessions()", "StopSession(ID2, stop-time)"},
		[]string{"GetOrCreateNamedProfile(profileName)"})
}

func TestAwsIamUserSessionActions_EditSession_FacadeGetSessionByIdReturnsError(t *testing.T) {
	awsIamUserSessionActionsSetup()

	sessionId := "ID"
	sessionName := "sessionName"
	region := "region"
	accountNumber := "accountNumber"
	userName := "userName"
	awsAccessKeyId := "awsAccessKeyId"
	awsSecretKey := "awsSecretKey"
	mfaDevice := "mfaDevice"
	profileName := "profileName"
	awsIamUserSessionActionsFacadeMock.ExpErrorOnGetSessionById = true

	err := awsIamUserSessionActions.EditSession(sessionId, sessionName, region, accountNumber, userName, awsAccessKeyId,
		awsSecretKey, mfaDevice, profileName)
	test.ExpectHttpError(t, err, http.StatusNotFound, "session not found")

	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{}, []string{"GetSessionById(ID)"},
		[]string{})
}

func TestAwsIamUserSessionActions_EditSession_KeychainSetSecretReturnsError(t *testing.T) {
	awsIamUserSessionActionsSetup()

	sessionId := "ID"
	sessionName := "sessionName"
	region := "region"
	accountNumber := "accountNumber"
	userName := "userName"
	awsAccessKeyId := "awsAccessKeyId"
	awsSecretKey := "awsSecretKey"
	mfaDevice := "mfaDevice"
	profileName := "profileName"
	awsIamUserSessionActionsFacadeMock.ExpGetSessionById = aws_iam_user.AwsIamUserSession{
		ID:                "ID",
		AccessKeyIDLabel:  "accessKeyIdLabel",
		SecretKeyLabel:    "secretKeyLabel",
		SessionTokenLabel: "sessionTokenLabel",
	}
	awsIamUserSessionActionsKeychainMock.ExpErrorOnSetSecret = true

	err := awsIamUserSessionActions.EditSession(sessionId, sessionName, region, accountNumber, userName, awsAccessKeyId,
		awsSecretKey, mfaDevice, profileName)
	test.ExpectHttpError(t, err, http.StatusUnprocessableEntity, "unable to set secret")
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{},
		[]string{"SetSecret(awsAccessKeyId, accessKeyIdLabel)"},
		[]string{"GetSessionById(ID)"},
		[]string{})
}

func TestAwsIamUserSessionActions_EditSession_GetOrCreateNamedProfileReturnsError(t *testing.T) {
	awsIamUserSessionActionsSetup()

	sessionId := "ID"
	sessionName := "sessionName"
	region := "region"
	accountNumber := "accountNumber"
	userName := "userName"
	awsAccessKeyId := "awsAccessKeyId"
	awsSecretKey := "awsSecretKey"
	mfaDevice := "mfaDevice"
	profileName := "profileName"
	awsIamUserSessionActionsFacadeMock.ExpGetSessionById = aws_iam_user.AwsIamUserSession{
		ID:                "ID",
		AccessKeyIDLabel:  "accessKeyIdLabel",
		SecretKeyLabel:    "secretKeyLabel",
		SessionTokenLabel: "sessionTokenLabel",
	}
	awsIamUserSessionNamedProfilesActionsMock.ExpErrorOnGetOrCreateNamedProfile = true
	err := awsIamUserSessionActions.EditSession(sessionId, sessionName, region, accountNumber, userName, awsAccessKeyId,
		awsSecretKey, mfaDevice, profileName)
	test.ExpectHttpError(t, err, http.StatusNotFound, "named profile not found")
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{},
		[]string{"SetSecret(awsAccessKeyId, accessKeyIdLabel)", "SetSecret(awsSecretKey, secretKeyLabel)"},
		[]string{"GetSessionById(ID)"},
		[]string{"GetOrCreateNamedProfile(profileName)"})
}

func TestAwsIamUserSessionActions_EditSession_FacadeEditSessionReturnsError(t *testing.T) {
	awsIamUserSessionActionsSetup()

	sessionId := "ID"
	sessionName := "sessionName"
	region := "region"
	accountNumber := "accountNumber"
	userName := "userName"
	awsAccessKeyId := "awsAccessKeyId"
	awsSecretKey := "awsSecretKey"
	mfaDevice := "mfaDevice"
	profileName := "profileName"
	awsIamUserSessionActionsFacadeMock.ExpGetSessionById = aws_iam_user.AwsIamUserSession{
		ID:                "ID",
		AccessKeyIDLabel:  "accessKeyIdLabel",
		SecretKeyLabel:    "secretKeyLabel",
		SessionTokenLabel: "sessionTokenLabel",
	}
	awsIamUserSessionActionsFacadeMock.ExpErrorOnEditSession = true
	awsIamUserSessionNamedProfilesActionsMock.ExpNamedProfile = named_profile.NamedProfile{Id: "ProfileId", Name: profileName}

	err := awsIamUserSessionActions.EditSession(sessionId, sessionName, region, accountNumber, userName, awsAccessKeyId,
		awsSecretKey, mfaDevice, profileName)
	test.ExpectHttpError(t, err, http.StatusConflict, "unable to edit session, collision detected")
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{},
		[]string{"SetSecret(awsAccessKeyId, accessKeyIdLabel)", "SetSecret(awsSecretKey, secretKeyLabel)"},
		[]string{"GetSessionById(ID)", "EditSession(ID, sessionName, region, accountNumber, userName, mfaDevice, ProfileId)"},
		[]string{"GetOrCreateNamedProfile(profileName)"})
}

func TestRotateSessionTokens(t *testing.T) {
	awsIamUserSessionActionsSetup()
	awsIamUserSessionActionsFacadeMock.ExpGetSessions = []aws_iam_user.AwsIamUserSession{{
		ID:                     "ID1",
		Status:                 domain_aws.Active,
		SessionTokenLabel:      "sessionTokenLabel1",
		SessionTokenExpiration: "2020-01-01T12:00:00Z",
	},
		{
			ID:     "ID2",
			Status: domain_aws.NotActive,
		},
	}
	awsIamUserSessionActionsEnvMock.ExpTime = "2020-01-01T11:00:00Z"
	awsIamUserSessionActionsKeychainMock.ExpSecretExist = true

	awsIamUserSessionActions.RotateSessionTokens()
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GetTime()"},
		[]string{"DoesSecretExist(sessionTokenLabel1)"},
		[]string{"GetSessions()"}, []string{})
}

func TestRotateSessionTokens_RefreshSessionTokenIfNeededReturnsError(t *testing.T) {
	awsIamUserSessionActionsSetup()
	awsIamUserSessionActionsFacadeMock.ExpGetSessions = []aws_iam_user.AwsIamUserSession{{
		ID:                     "ID1",
		Status:                 domain_aws.Active,
		AccessKeyIDLabel:       "accessKeyIdLabel1",
		SessionTokenLabel:      "sessionTokenLabel1",
		SessionTokenExpiration: "2020-01-01T12:00:00Z",
	},
		{
			ID:     "ID2",
			Status: domain_aws.NotActive,
		},
	}

	awsIamUserSessionActionsEnvMock.ExpTime = "2020-01-01T13:00:00Z"
	awsIamUserSessionActionsKeychainMock.ExpSecretExist = true
	awsIamUserSessionActionsKeychainMock.ExpErrorOnGetSecret = true

	awsIamUserSessionActions.RotateSessionTokens()
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{"GetTime()"},
		[]string{"DoesSecretExist(sessionTokenLabel1)", "GetSecret(accessKeyIdLabel1)"},
		[]string{"GetSessions()"}, []string{})
}

func TestIsSessionTokenValid(t *testing.T) {
	awsIamUserSessionActionsSetup()

	sessionTokenLabel := "sessionTokenLabel"
	sessionTokenExpiration := "2020-01-01T12:00:00Z"
	currentTime := "2020-01-01T11:00:00Z"
	awsIamUserSessionActionsKeychainMock.ExpSecretExist = true

	isSessionTokenValid := awsIamUserSessionActions.isSessionTokenValid(sessionTokenLabel, sessionTokenExpiration, currentTime)
	if !isSessionTokenValid {
		t.Fatalf("Unexpected result")
	}
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{"DoesSecretExist(sessionTokenLabel)"},
		[]string{}, []string{})
}

func TestIsSessionTokenValid_SecretDoesNotExist(t *testing.T) {
	awsIamUserSessionActionsSetup()

	sessionTokenLabel := "sessionTokenLabel"
	sessionTokenExpiration := "2020-01-01T12:00:00Z"
	currentTime := "2020-01-01T11:00:00Z"

	isSessionTokenValid := awsIamUserSessionActions.isSessionTokenValid(sessionTokenLabel, sessionTokenExpiration, currentTime)
	if isSessionTokenValid {
		t.Fatalf("Unexpected result")
	}
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{"DoesSecretExist(sessionTokenLabel)"},
		[]string{}, []string{})
}

func TestIsSessionTokenValid_DoesSecretExistReturnsError(t *testing.T) {
	awsIamUserSessionActionsSetup()

	sessionTokenLabel := "sessionTokenLabel"
	sessionTokenExpiration := "2020-01-01T12:00:00Z"
	currentTime := "2020-01-01T11:00:00Z"
	awsIamUserSessionActionsKeychainMock.ExpErrorOnSecretExist = true

	isSessionTokenValid := awsIamUserSessionActions.isSessionTokenValid(sessionTokenLabel, sessionTokenExpiration, currentTime)
	if isSessionTokenValid {
		t.Fatalf("Unexpected result")
	}
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{"DoesSecretExist(sessionTokenLabel)"},
		[]string{}, []string{})
}

func TestIsSessionTokenValid_sessionTokenExpirationIsEmpty(t *testing.T) {
	awsIamUserSessionActionsSetup()

	sessionTokenLabel := "sessionTokenLabel"
	currentTime := "2020-01-01T11:00:00Z"
	awsIamUserSessionActionsKeychainMock.ExpSecretExist = true

	isSessionTokenValid := awsIamUserSessionActions.isSessionTokenValid(sessionTokenLabel, "", currentTime)
	if isSessionTokenValid {
		t.Fatalf("Unexpected result")
	}
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{"DoesSecretExist(sessionTokenLabel)"},
		[]string{}, []string{})
}

func TestIsSessionTokenValid_CurrentTimeIsNotParsable(t *testing.T) {
	awsIamUserSessionActionsSetup()

	sessionTokenLabel := "sessionTokenLabel"
	currentTime := "Haloa!"
	sessionTokenExpiration := "2020-01-01T12:00:00Z"
	awsIamUserSessionActionsKeychainMock.ExpSecretExist = true

	isSessionTokenValid := awsIamUserSessionActions.isSessionTokenValid(sessionTokenLabel, sessionTokenExpiration, currentTime)
	if isSessionTokenValid {
		t.Fatalf("Unexpected result")
	}
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{"DoesSecretExist(sessionTokenLabel)"},
		[]string{}, []string{})
}

func TestIsSessionTokenValid_SessionTokenExpirationNotParsable(t *testing.T) {
	awsIamUserSessionActionsSetup()

	sessionTokenLabel := "sessionTokenLabel"
	currentTime := "2020-01-01T12:00:00Z"
	sessionTokenExpiration := "Hello!"
	awsIamUserSessionActionsKeychainMock.ExpSecretExist = true

	isSessionTokenValid := awsIamUserSessionActions.isSessionTokenValid(sessionTokenLabel, sessionTokenExpiration, currentTime)
	if isSessionTokenValid {
		t.Fatalf("Unexpected result")
	}
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{"DoesSecretExist(sessionTokenLabel)"},
		[]string{}, []string{})
}

func TestIsSessionTokenValid_SessionIsExpired(t *testing.T) {
	awsIamUserSessionActionsSetup()

	sessionTokenLabel := "sessionTokenLabel"
	currentTime := "2020-01-01T12:00:00Z"
	sessionTokenExpiration := "2020-01-01T11:00:00Z"
	awsIamUserSessionActionsKeychainMock.ExpSecretExist = true

	isSessionTokenValid := awsIamUserSessionActions.isSessionTokenValid(sessionTokenLabel, sessionTokenExpiration, currentTime)
	if isSessionTokenValid {
		t.Fatalf("Unexpected result")
	}
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{"DoesSecretExist(sessionTokenLabel)"},
		[]string{}, []string{})
}

func TestRefreshSessionToken(t *testing.T) {
	awsIamUserSessionActionsSetup()
	session := aws_iam_user.AwsIamUserSession{
		ID:                "ID",
		AccessKeyIDLabel:  "accessKeyLabel",
		SecretKeyLabel:    "secretKeyLabel",
		SessionTokenLabel: "sessionTokenLabel",
		Region:            "region-1",
		MfaDevice:         "mfaDevice",
	}
	awsIamUserSessionActionsKeychainMock.ExpGetSecret = "label"
	currentTime, _ := time.Parse(time.RFC3339, "2020-01-01T11:00:00Z")
	stsApiMock.ExpCredentials = sts.Credentials{Expiration: &currentTime}

	err := awsIamUserSessionActions.refreshSessionToken(session)
	if err != nil {
		t.Fatalf("Unexpected error")
	}

	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{"GenerateNewSessionToken(label, label, region-1, mfaDevice, <nil>)"},
		[]string{}, []string{"GetSecret(accessKeyLabel)", "GetSecret(secretKeyLabel)", "SetSecret({\"AccessKeyId\":null,\"Expiration\":\"2020-01-01T11:00:00Z\",\"SecretAccessKey\":null,\"SessionToken\":null}, sessionTokenLabel)"},
		[]string{"SetSessionTokenExpiration(ID, 2020-01-01T11:00:00Z)"}, []string{})
}

func TestRefreshSessionToken_KeychainGetSecretReturnsError(t *testing.T) {
	awsIamUserSessionActionsSetup()
	session := aws_iam_user.AwsIamUserSession{
		ID:                "ID",
		AccessKeyIDLabel:  "accessKeyLabel",
		SecretKeyLabel:    "secretKeyLabel",
		SessionTokenLabel: "sessionTokenLabel",
		Region:            "region-1",
		MfaDevice:         "mfaDevice",
	}
	awsIamUserSessionActionsKeychainMock.ExpErrorOnGetSecret = true
	currentTime, _ := time.Parse(time.RFC3339, "2020-01-01T11:00:00Z")
	stsApiMock.ExpCredentials = sts.Credentials{Expiration: &currentTime}

	err := awsIamUserSessionActions.refreshSessionToken(session)
	test.ExpectHttpError(t, err, http.StatusInternalServerError, "unable to get secret")
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{}, []string{}, []string{"GetSecret(accessKeyLabel)"},
		[]string{}, []string{})
}

func TestRefreshSessionToken_StsApiGenerateNewSessionTokenReturnsError(t *testing.T) {
	awsIamUserSessionActionsSetup()
	session := aws_iam_user.AwsIamUserSession{
		ID:                "ID",
		AccessKeyIDLabel:  "accessKeyLabel",
		SecretKeyLabel:    "secretKeyLabel",
		SessionTokenLabel: "sessionTokenLabel",
		Region:            "region-1",
		MfaDevice:         "mfaDevice",
	}
	awsIamUserSessionActionsKeychainMock.ExpGetSecret = "label"
	stsApiMock.ExpErrorOnGenerateNewSessionToken = true

	err := awsIamUserSessionActions.refreshSessionToken(session)
	test.ExpectHttpError(t, err, http.StatusUnprocessableEntity, "unable to generate new session token")
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{"GenerateNewSessionToken(label, label, region-1, mfaDevice, <nil>)"},
		[]string{}, []string{"GetSecret(accessKeyLabel)", "GetSecret(secretKeyLabel)"}, []string{}, []string{})
}

func TestRefreshSessionToken_KeychainSetSecretReturnsError(t *testing.T) {
	awsIamUserSessionActionsSetup()
	session := aws_iam_user.AwsIamUserSession{
		ID:                "ID",
		AccessKeyIDLabel:  "accessKeyLabel",
		SecretKeyLabel:    "secretKeyLabel",
		SessionTokenLabel: "sessionTokenLabel",
		Region:            "region-1",
		MfaDevice:         "mfaDevice",
	}
	awsIamUserSessionActionsKeychainMock.ExpGetSecret = "label"
	awsIamUserSessionActionsKeychainMock.ExpErrorOnSetSecret = true
	currentTime, _ := time.Parse(time.RFC3339, "2020-01-01T11:00:00Z")
	stsApiMock.ExpCredentials = sts.Credentials{Expiration: &currentTime}

	err := awsIamUserSessionActions.refreshSessionToken(session)
	test.ExpectHttpError(t, err, http.StatusUnprocessableEntity, "unable to set secret")
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{"GenerateNewSessionToken(label, label, region-1, mfaDevice, <nil>)"},
		[]string{}, []string{"GetSecret(accessKeyLabel)", "GetSecret(secretKeyLabel)", "SetSecret({\"AccessKeyId\":null,\"Expiration\":\"2020-01-01T11:00:00Z\",\"SecretAccessKey\":null,\"SessionToken\":null}, sessionTokenLabel)"},
		[]string{}, []string{})
}

func TestRefreshSessionToken_FacadeSetSessionTokenExpirationReturnsError(t *testing.T) {
	awsIamUserSessionActionsSetup()
	session := aws_iam_user.AwsIamUserSession{
		ID:                "ID",
		AccessKeyIDLabel:  "accessKeyLabel",
		SecretKeyLabel:    "secretKeyLabel",
		SessionTokenLabel: "sessionTokenLabel",
		Region:            "region-1",
		MfaDevice:         "mfaDevice",
	}
	awsIamUserSessionActionsFacadeMock.ExpErrorOnSetSessionTokenExpiration = true
	awsIamUserSessionActionsKeychainMock.ExpGetSecret = "label"
	currentTime, _ := time.Parse(time.RFC3339, "2020-01-01T11:00:00Z")
	stsApiMock.ExpCredentials = sts.Credentials{Expiration: &currentTime}

	err := awsIamUserSessionActions.refreshSessionToken(session)
	test.ExpectHttpError(t, err, http.StatusInternalServerError, "unable to set token expiration")
	awsIamUserSessionActionsVerifyExpectedCalls(t, []string{"GenerateNewSessionToken(label, label, region-1, mfaDevice, <nil>)"},
		[]string{}, []string{"GetSecret(accessKeyLabel)", "GetSecret(secretKeyLabel)", "SetSecret({\"AccessKeyId\":null,\"Expiration\":\"2020-01-01T11:00:00Z\",\"SecretAccessKey\":null,\"SessionToken\":null}, sessionTokenLabel)"},
		[]string{"SetSessionTokenExpiration(ID, 2020-01-01T11:00:00Z)"}, []string{})
}
