package use_case

import (
	"leapp_daemon/domain/domain_gcp"
	"leapp_daemon/domain/domain_gcp/gcp_iam_user_account_oauth"
	"leapp_daemon/domain/domain_gcp/named_configuration"
	"leapp_daemon/test/mock"
	"reflect"
	"testing"
)

var (
	gcpKeychainMock          mock.KeychainMock
	gcpRepoMock              mock.GcpConfigurationRepositoryMock
	gcpNamedConfigFacadeMock mock.NamedConfigurationsFacadeMock
	gcpCredentialsApplier    *GcpCredentialsApplier
)

func gcpCredentialsApplierSetup() {
	gcpKeychainMock = mock.NewKeychainMock()
	gcpRepoMock = mock.NewGcpConfigurationRepositoryMock()
	gcpNamedConfigFacadeMock = mock.NewNamedConfigurationsFacadeMock()
	gcpCredentialsApplier = &GcpCredentialsApplier{
		Keychain:                  &gcpKeychainMock,
		Repository:                &gcpRepoMock,
		NamedConfigurationsFacade: &gcpNamedConfigFacadeMock,
	}
}

func gcpCredentialsApplierVerifyExpectedCalls(t *testing.T, keychainMockCalls []string, repoMockCalls []string,
	namedConfigFacadeMockCalls []string) {
	if !reflect.DeepEqual(gcpKeychainMock.GetCalls(), keychainMockCalls) {
		t.Fatalf("keychainMock expectation violation.\nMock calls: %v", gcpKeychainMock.GetCalls())
	}
	if !reflect.DeepEqual(gcpRepoMock.GetCalls(), repoMockCalls) {
		t.Fatalf("gcpRepoMock expectation violation.\nMock calls: %v", gcpRepoMock.GetCalls())
	}
	if !reflect.DeepEqual(gcpNamedConfigFacadeMock.GetCalls(), namedConfigFacadeMockCalls) {
		t.Fatalf("gcpNamedConfigFacadeMock expectation violation.\nMock calls: %v", gcpNamedConfigFacadeMock.GetCalls())
	}
}

func TestUpdateGcpIamUserAccountOauthSessions_DeactivateNonDefaultSession(t *testing.T) {
	gcpCredentialsApplierSetup()
	gcpNamedConfigFacadeMock.ExpNamedConfiguration = named_configuration.NamedConfiguration{Name: "named-config"}
	oldSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{
		{Status: domain_gcp.Active, AccountId: "account-id", NamedConfigurationId: "config-id"}}
	var newSessions []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession

	gcpCredentialsApplier.UpdateGcpIamUserAccountOauthSessions(oldSessions, newSessions)
	gcpCredentialsApplierVerifyExpectedCalls(t, []string{}, []string{"RemoveConfiguration(named-config)",
		"RemoveCredentialsFromDb(account-id)", "RemoveAccessTokensFromDb(account-id)"},
		[]string{"GetNamedConfigurationById(config-id)"})
}

func TestUpdateGcpIamUserAccountOauthSessions_DeactivateDefaultSession(t *testing.T) {
	gcpCredentialsApplierSetup()
	gcpNamedConfigFacadeMock.ExpNamedConfiguration = named_configuration.NamedConfiguration{
		Name: named_configuration.DefaultNamedConfigurationName,
	}
	oldSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{
		{Status: domain_gcp.Active, AccountId: "account-id", NamedConfigurationId: "config-id"},
	}
	var newSessions []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession

	gcpCredentialsApplier.UpdateGcpIamUserAccountOauthSessions(oldSessions, newSessions)
	gcpCredentialsApplierVerifyExpectedCalls(t, []string{}, []string{"RemoveConfiguration(leapp-default)",
		"RemoveCredentialsFromDb(account-id)", "RemoveAccessTokensFromDb(account-id)", "RemoveDefaultCredentials()",
		"DeactivateConfiguration()"},
		[]string{"GetNamedConfigurationById(config-id)"})
}

func TestUpdateGcpIamUserAccountOauthSessions_ActivateNonDefaultSessions(t *testing.T) {
	gcpCredentialsApplierSetup()
	var oldSessions []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession
	newSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{
		{Status: domain_gcp.Active, AccountId: "account-id", NamedConfigurationId: "config-id",
			CredentialsLabel: "credentialLabel", ProjectName: "projectName"},
	}
	gcpNamedConfigFacadeMock.ExpNamedConfiguration = named_configuration.NamedConfiguration{Name: "named-config"}
	gcpKeychainMock.ExpGetSecret = "credentials"

	gcpCredentialsApplier.UpdateGcpIamUserAccountOauthSessions(oldSessions, newSessions)
	gcpCredentialsApplierVerifyExpectedCalls(t, []string{"GetSecret(credentialLabel)"},
		[]string{"WriteCredentialsToDb(account-id, credentials)", "CreateConfiguration(account-id, projectName, named-config)"},
		[]string{"GetNamedConfigurationById(config-id)"})
}

func TestUpdateGcpIamUserAccountOauthSessions_ActivateDefaultSessions(t *testing.T) {
	gcpCredentialsApplierSetup()
	var oldSessions []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession
	newSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{
		{Status: domain_gcp.Active, AccountId: "account-id", NamedConfigurationId: "config-id",
			CredentialsLabel: "credentialLabel", ProjectName: "projectName"},
	}
	gcpNamedConfigFacadeMock.ExpNamedConfiguration = named_configuration.NamedConfiguration{
		Name: named_configuration.DefaultNamedConfigurationName,
	}
	gcpKeychainMock.ExpGetSecret = "credentials"

	gcpCredentialsApplier.UpdateGcpIamUserAccountOauthSessions(oldSessions, newSessions)
	gcpCredentialsApplierVerifyExpectedCalls(t, []string{"GetSecret(credentialLabel)"},
		[]string{"WriteCredentialsToDb(account-id, credentials)", "CreateConfiguration(account-id, projectName, leapp-default)",
			"ActivateConfiguration()", "WriteDefaultCredentials(credentials)"},
		[]string{"GetNamedConfigurationById(config-id)"})
}

func TestUpdateGcpIamUserAccountOauthSessions_DoNothing(t *testing.T) {
	gcpCredentialsApplierSetup()
	oldSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Status: domain_gcp.NotActive}}
	newSessions := []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Status: domain_gcp.NotActive}}

	gcpCredentialsApplier.UpdateGcpIamUserAccountOauthSessions(oldSessions, newSessions)
	gcpCredentialsApplierVerifyExpectedCalls(t, []string{}, []string{}, []string{})
}
