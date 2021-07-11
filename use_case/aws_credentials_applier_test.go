package use_case

import (
	"leapp_daemon/domain/domain_aws"
	"leapp_daemon/domain/domain_aws/aws_iam_user"
	"leapp_daemon/domain/domain_aws/named_profile"
	"leapp_daemon/test/mock"
	"reflect"
	"testing"
)

var (
	awsCredentialApplierKeychainMock           mock.KeychainMock
	awsCredentialApplierNamedProfileFacadeMock mock.NamedProfilesFacadeMock
	awsCredentialApplierFileRepoMock           mock.AwsConfigurationRepositoryMock
	awsCredentialApplier                       *AwsCredentialsApplier
)

func awsCredentialApplierSetup() {
	awsCredentialApplierKeychainMock = mock.NewKeychainMock()
	awsCredentialApplierNamedProfileFacadeMock = mock.NewNamedProfilesFacadeMock()
	awsCredentialApplierFileRepoMock = mock.NewAwsConfigurationRepositoryMock()
	awsCredentialApplier = &AwsCredentialsApplier{
		Keychain:            &awsCredentialApplierKeychainMock,
		NamedProfilesFacade: &awsCredentialApplierNamedProfileFacadeMock,
		Repository:          &awsCredentialApplierFileRepoMock,
	}
}

func awsCredentialApplierVerifyExpectedCalls(t *testing.T, keychainMockCalls []string, namedProfileFacadeMockCalls []string, repoMockCalls []string) {
	if !reflect.DeepEqual(awsCredentialApplierKeychainMock.GetCalls(), keychainMockCalls) {
		t.Fatalf("keychainMock expectation violation.\nMock calls: %v", awsCredentialApplierKeychainMock.GetCalls())
	}
	if !reflect.DeepEqual(awsCredentialApplierNamedProfileFacadeMock.GetCalls(), namedProfileFacadeMockCalls) {
		t.Fatalf("namedProfileFacadeMock expectation violation.\nMock calls: %v", awsCredentialApplierNamedProfileFacadeMock.GetCalls())
	}
	if !reflect.DeepEqual(awsCredentialApplierFileRepoMock.GetCalls(), repoMockCalls) {
		t.Fatalf("awsRepoMock expectation violation.\nMock calls: %v", awsCredentialApplierFileRepoMock.GetCalls())
	}
}

func TestAwsCredentialsApplier_UpdateAwsIamUserSessions(t *testing.T) {
	awsCredentialApplierSetup()
	awsCredentialApplierKeychainMock.ExpGetSecret = "{" +
		"\"AccessKeyId\":\"accessKeyId\", " +
		"\"SecretAccessKey\":\"secretAccessKey\", " +
		"\"SessionToken\":\"sessionToken\", " +
		"\"Expiration\":\"expiration\"" +
		"}"
	awsCredentialApplierNamedProfileFacadeMock.ExpNamedProfile = named_profile.NamedProfile{Id: "ProfileId", Name: "profileName"}

	oldSessions := []aws_iam_user.AwsIamUserSession{{ID: "ID1", Status: domain_aws.Active}}
	newSessions := []aws_iam_user.AwsIamUserSession{
		{ID: "ID1", Status: domain_aws.Active, SessionTokenLabel: "sessionTokenLabel", NamedProfileID: "ProfileId", Region: "region-1"},
		{ID: "ID2", Status: domain_aws.NotActive}}

	awsCredentialApplier.UpdateAwsIamUserSessions(oldSessions, newSessions)
	awsCredentialApplierVerifyExpectedCalls(t, []string{"GetSecret(sessionTokenLabel)"},
		[]string{"GetNamedProfileById(ProfileId)"},
		[]string{"WriteCredentials([{profileName accessKeyId secretAccessKey sessionToken expiration region-1}])"})
}

func TestAwsCredentialsApplier_UpdateAwsIamUserSessions_NoActiveSessions(t *testing.T) {
	awsCredentialApplierSetup()
	oldSessions := []aws_iam_user.AwsIamUserSession{{ID: "ID1", Status: domain_aws.Active}}
	newSessions := []aws_iam_user.AwsIamUserSession{{ID: "ID1", Status: domain_aws.NotActive}}

	awsCredentialApplier.UpdateAwsIamUserSessions(oldSessions, newSessions)
	awsCredentialApplierVerifyExpectedCalls(t, []string{}, []string{}, []string{"WriteCredentials([])"})
}
