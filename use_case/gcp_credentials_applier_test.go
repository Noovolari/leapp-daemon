package use_case

import (
	"leapp_daemon/domain/session"
	"leapp_daemon/test/mock"
	"reflect"
	"testing"
)

var (
	keychainMock              mock.KeychainMock
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

	keychainMock = mock.NewKeychainMock()
	gcpRepoMock = mock.NewGcpConfigurationRepositoryMock()
	gcpCredentialsApplier = &GcpCredentialsApplier{
		Keychain:   &keychainMock,
		Repository: &gcpRepoMock,
	}
}

func gcpCredentialsApplierVerifyExpectedCalls(t *testing.T, keychainMockCalls []string, repoMockCalls []string) {
	if !reflect.DeepEqual(keychainMock.GetCalls(), keychainMockCalls) {
		t.Fatalf("keychainMock expectation violation.\nMock calls: %v", keychainMock.GetCalls())
	}
	if !reflect.DeepEqual(gcpRepoMock.GetCalls(), repoMockCalls) {
		t.Fatalf("gcpRepoMock expectation violation.\nMock calls: %v", gcpRepoMock.GetCalls())
	}
}

func TestUpdateGcpIamUserAccountOauthSessions_OldActiveSessionAndNoNewActiveSessions(t *testing.T) {
	gcpCredentialsApplierSetup()
	oldSessions := []session.GcpIamUserAccountOauthSession{{Status: session.Active}}
	newSessions := []session.GcpIamUserAccountOauthSession{}

	gcpCredentialsApplier.UpdateGcpIamUserAccountOauthSessions(oldSessions, newSessions)
	gcpCredentialsApplierVerifyExpectedCalls(t, []string{}, expectedDeactivationCalls)
}

func TestUpdateGcpIamUserAccountOauthSessions_OldAndNewActiveSessionWithDifferentIds(t *testing.T) {
	gcpCredentialsApplierSetup()
	keychainMock.ExpGetSecret = "credentials"
	oldSessions := []session.GcpIamUserAccountOauthSession{{Id: "ID1", Status: session.Active}}
	newSessions := []session.GcpIamUserAccountOauthSession{{Id: "ID2", CredentialsLabel: "credentialsLabel", AccountId: "accountId",
		Name: "sessionName", ProjectName: "projectName", Status: session.Active}}

	gcpCredentialsApplier.UpdateGcpIamUserAccountOauthSessions(oldSessions, newSessions)
	expectedRepositoryCalls := append(expectedDeactivationCalls, expectedActivationCalls...)
	gcpCredentialsApplierVerifyExpectedCalls(t, []string{"GetSecret(credentialsLabel)"}, expectedRepositoryCalls)
}

func TestUpdateGcpIamUserAccountOauthSessions_OldAndNewActiveSessionAreEqual(t *testing.T) {
	gcpCredentialsApplierSetup()
	oldSessions := []session.GcpIamUserAccountOauthSession{{Id: "ID1", Status: session.Active}}
	newSessions := []session.GcpIamUserAccountOauthSession{{Id: "ID1", Status: session.Active}}

	gcpCredentialsApplier.UpdateGcpIamUserAccountOauthSessions(oldSessions, newSessions)
	gcpCredentialsApplierVerifyExpectedCalls(t, []string{}, []string{})
}

func TestUpdateGcpIamUserAccountOauthSessions_OldAndNewActiveSessionWithSameIdsButDifferentParams(t *testing.T) {
	gcpCredentialsApplierSetup()
	keychainMock.ExpGetSecret = "credentials"
	oldSessions := []session.GcpIamUserAccountOauthSession{{Id: "ID1", CredentialsLabel: "credentialsLabel", AccountId: "accountId",
		Name: "sessionName", ProjectName: "oldProjectName", Status: session.Active}}
	newSessions := []session.GcpIamUserAccountOauthSession{{Id: "ID1", CredentialsLabel: "credentialsLabel", AccountId: "accountId",
		Name: "sessionName", ProjectName: "projectName", Status: session.Active}}

	gcpCredentialsApplier.UpdateGcpIamUserAccountOauthSessions(oldSessions, newSessions)
	gcpCredentialsApplierVerifyExpectedCalls(t, []string{"GetSecret(credentialsLabel)"}, expectedActivationCalls)
}

func TestUpdateGcpIamUserAccountOauthSessions_NoOldActiveSessionButNewActiveSessionPresent(t *testing.T) {
	gcpCredentialsApplierSetup()
	keychainMock.ExpGetSecret = "credentials"
	oldSessions := []session.GcpIamUserAccountOauthSession{{Id: "ID1", Status: session.NotActive}}
	newSessions := []session.GcpIamUserAccountOauthSession{{Id: "ID1", CredentialsLabel: "credentialsLabel", AccountId: "accountId",
		Name: "sessionName", ProjectName: "projectName", Status: session.Active}}

	gcpCredentialsApplier.UpdateGcpIamUserAccountOauthSessions(oldSessions, newSessions)
	gcpCredentialsApplierVerifyExpectedCalls(t, []string{"GetSecret(credentialsLabel)"}, expectedActivationCalls)
}

func TestUpdateGcpIamUserAccountOauthSessions_NoActiveSessions(t *testing.T) {
	gcpCredentialsApplierSetup()
	keychainMock.ExpGetSecret = "credentials"
	oldSessions := []session.GcpIamUserAccountOauthSession{{Id: "ID1", Status: session.NotActive}}
	newSessions := []session.GcpIamUserAccountOauthSession{{Id: "ID1", Status: session.Pending}}

	gcpCredentialsApplier.UpdateGcpIamUserAccountOauthSessions(oldSessions, newSessions)
	gcpCredentialsApplierVerifyExpectedCalls(t, []string{}, []string{})
}
