package use_case

import (
	"leapp_daemon/domain/domain_gcp"
	"leapp_daemon/domain/domain_gcp/gcp_iam_user_account_oauth"
	"leapp_daemon/test/mock"
	"reflect"
	"testing"
)

var (
	gcpKeychainMock           mock.KeychainMock
	gcpRepoMock               mock.GcpConfigurationRepositoryMock
	gcpCredentialsApplier     *GcpCredentialsApplier
	expectedDeactivationCalls []string
	expectedActivationCalls   []string
)

func gcpCredentialsApplierSetup() {
	expectedDeactivationCalls = []string{"RemoveDefaultCredentials()", "DeactivateConfiguration()",
		"RemoveCredentialsFromDb()", "RemoveAccessTokensFromDb()", "RemoveConfiguration()"}
	expectedActivationCalls = []string{"WriteDefaultCredentials(accountId, credentials)",
		"CreateConfiguration(accountId, projectName)", "ActivateConfiguration()", "WriteDefaultCredentials(credentials)"}

	gcpKeychainMock = mock.NewKeychainMock()
	gcpRepoMock = mock.NewGcpConfigurationRepositoryMock()
	gcpCredentialsApplier = &GcpCredentialsApplier{
		Keychain:   &gcpKeychainMock,
		Repository: &gcpRepoMock,
	}
}

func gcpCredentialsApplierVerifyExpectedCalls(t *testing.T, keychainMockCalls []string, repoMockCalls []string) {
	if !reflect.DeepEqual(gcpKeychainMock.GetCalls(), keychainMockCalls) {
		t.Fatalf("keychainMock expectation violation.\nMock calls: %v", gcpKeychainMock.GetCalls())
	}
	if !reflect.DeepEqual(gcpRepoMock.GetCalls(), repoMockCalls) {
		t.Fatalf("gcpRepoMock expectation violation.\nMock calls: %v", gcpRepoMock.GetCalls())
	}
}

func TestUpdateGcpIamUserAccountOauthSessions_OldActiveSessionAndNoNewActiveSessions(t *testing.T) {
	gcpCredentialsApplierSetup()
	oldSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Status: domain_gcp.Active}}
	newSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{}

	gcpCredentialsApplier.UpdateGcpIamUserAccountOauthSessions(oldSessions, newSessions)
	gcpCredentialsApplierVerifyExpectedCalls(t, []string{}, expectedDeactivationCalls)
}

func TestUpdateGcpIamUserAccountOauthSessions_OldAndNewActiveSessionWithDifferentIds(t *testing.T) {
	gcpCredentialsApplierSetup()
	gcpKeychainMock.ExpGetSecret = "credentials"
	oldSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Id: "ID1", Status: domain_gcp.Active}}
	newSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Id: "ID2", CredentialsLabel: "credentialsLabel", AccountId: "accountId",
		Name: "sessionName", ProjectName: "projectName", Status: domain_gcp.Active}}

	gcpCredentialsApplier.UpdateGcpIamUserAccountOauthSessions(oldSessions, newSessions)
	expectedRepositoryCalls := append(expectedDeactivationCalls, expectedActivationCalls...)
	gcpCredentialsApplierVerifyExpectedCalls(t, []string{"GetSecret(credentialsLabel)"}, expectedRepositoryCalls)
}

func TestUpdateGcpIamUserAccountOauthSessions_OldAndNewActiveSessionAreEqual(t *testing.T) {
	gcpCredentialsApplierSetup()
	oldSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Id: "ID1", Status: domain_gcp.Active}}
	newSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Id: "ID1", Status: domain_gcp.Active}}

	gcpCredentialsApplier.UpdateGcpIamUserAccountOauthSessions(oldSessions, newSessions)
	gcpCredentialsApplierVerifyExpectedCalls(t, []string{}, []string{})
}

func TestUpdateGcpIamUserAccountOauthSessions_OldAndNewActiveSessionWithSameIdsButDifferentParams(t *testing.T) {
	gcpCredentialsApplierSetup()
	gcpKeychainMock.ExpGetSecret = "credentials"
	oldSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Id: "ID1", CredentialsLabel: "credentialsLabel", AccountId: "accountId",
		Name: "sessionName", ProjectName: "oldProjectName", Status: domain_gcp.Active}}
	newSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Id: "ID1", CredentialsLabel: "credentialsLabel", AccountId: "accountId",
		Name: "sessionName", ProjectName: "projectName", Status: domain_gcp.Active}}

	gcpCredentialsApplier.UpdateGcpIamUserAccountOauthSessions(oldSessions, newSessions)
	gcpCredentialsApplierVerifyExpectedCalls(t, []string{"GetSecret(credentialsLabel)"}, expectedActivationCalls)
}

func TestUpdateGcpIamUserAccountOauthSessions_NoOldActiveSessionButNewActiveSessionPresent(t *testing.T) {
	gcpCredentialsApplierSetup()
	gcpKeychainMock.ExpGetSecret = "credentials"
	oldSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Id: "ID1", Status: domain_gcp.NotActive}}
	newSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Id: "ID1", CredentialsLabel: "credentialsLabel", AccountId: "accountId",
		Name: "sessionName", ProjectName: "projectName", Status: domain_gcp.Active}}

	gcpCredentialsApplier.UpdateGcpIamUserAccountOauthSessions(oldSessions, newSessions)
	gcpCredentialsApplierVerifyExpectedCalls(t, []string{"GetSecret(credentialsLabel)"}, expectedActivationCalls)
}

func TestUpdateGcpIamUserAccountOauthSessions_NoActiveSessions(t *testing.T) {
	gcpCredentialsApplierSetup()
	gcpKeychainMock.ExpGetSecret = "credentials"
	oldSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Id: "ID1", Status: domain_gcp.NotActive}}
	newSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Id: "ID1", Status: domain_gcp.NotActive}}

	gcpCredentialsApplier.UpdateGcpIamUserAccountOauthSessions(oldSessions, newSessions)
	gcpCredentialsApplierVerifyExpectedCalls(t, []string{}, []string{})
}
