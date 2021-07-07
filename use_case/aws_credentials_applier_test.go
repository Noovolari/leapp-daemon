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

	oldSessions := []aws_iam_user.AwsIamUserSession{{Id: "ID1", Status: domain_aws.Active}}
	newSessions := []aws_iam_user.AwsIamUserSession{
		{Id: "ID1", Status: domain_aws.Active, SessionTokenLabel: "sessionTokenLabel", NamedProfileId: "ProfileId", Region: "region-1"},
		{Id: "ID2", Status: domain_aws.NotActive}}

	awsCredentialApplier.UpdateAwsIamUserSessions(oldSessions, newSessions)
	awsCredentialApplierVerifyExpectedCalls(t, []string{"GetSecret(sessionTokenLabel)"},
		[]string{"GetNamedProfileById(ProfileId)"},
		[]string{"WriteCredentials([{profileName accessKeyId secretAccessKey sessionToken expiration region-1}])"})
}

func TestAwsCredentialsApplier_UpdateAwsIamUserSessions_NoActiveSessions(t *testing.T) {
	awsCredentialApplierSetup()
	oldSessions := []aws_iam_user.AwsIamUserSession{{Id: "ID1", Status: domain_aws.Active}}
	newSessions := []aws_iam_user.AwsIamUserSession{{Id: "ID1", Status: domain_aws.NotActive}}

	awsCredentialApplier.UpdateAwsIamUserSessions(oldSessions, newSessions)
	awsCredentialApplierVerifyExpectedCalls(t, []string{}, []string{}, []string{"WriteCredentials([])"})
}

func TestAwsCredentialsApplier_UpdateAwsImUserSessions_FacadeGetNamedProfileByIdReturnsError(t *testing.T) {
	awsCredentialApplierSetup()
	oldSessions := []aws_iam_user.AwsIamUserSession{}
	newSessions := []aws_iam_user.AwsIamUserSession{{Id: "ID1", Status: domain_aws.Active, NamedProfileId: "ProfileId"}}
	awsCredentialApplierNamedProfileFacadeMock.ExpErrorOnGetNamedProfileById = true

	awsCredentialApplier.UpdateAwsIamUserSessions(oldSessions, newSessions)
	awsCredentialApplierVerifyExpectedCalls(t, []string{}, []string{"GetNamedProfileById(ProfileId)"}, []string{})
}

func TestAwsCredentialsApplier_UpdateAwsImUserSessions_KeychainGetSecretReturnsError(t *testing.T) {
	awsCredentialApplierSetup()
	oldSessions := []aws_iam_user.AwsIamUserSession{}
	newSessions := []aws_iam_user.AwsIamUserSession{
		{Id: "ID1", Status: domain_aws.Active, NamedProfileId: "ProfileId", SessionTokenLabel: "sessionTokenLabel"},
	}
	awsCredentialApplierNamedProfileFacadeMock.ExpNamedProfile = named_profile.NamedProfile{Id: "ProfileId", Name: "profileName"}
	awsCredentialApplierKeychainMock.ExpErrorOnGetSecret = true

	awsCredentialApplier.UpdateAwsIamUserSessions(oldSessions, newSessions)
	awsCredentialApplierVerifyExpectedCalls(t, []string{"GetSecret(sessionTokenLabel)"},
		[]string{"GetNamedProfileById(ProfileId)"}, []string{})
}

func TestAwsCredentialsApplier_UpdateAwsImUserSessions_UnmarshallSessionTokenReturnsError(t *testing.T) {
	awsCredentialApplierSetup()
	oldSessions := []aws_iam_user.AwsIamUserSession{}
	newSessions := []aws_iam_user.AwsIamUserSession{
		{Id: "ID1", Status: domain_aws.Active, NamedProfileId: "ProfileId", SessionTokenLabel: "sessionTokenLabel"},
	}
	awsCredentialApplierNamedProfileFacadeMock.ExpNamedProfile = named_profile.NamedProfile{Id: "ProfileId", Name: "profileName"}
	awsCredentialApplierKeychainMock.ExpGetSecret = "{this is not a json}"

	awsCredentialApplier.UpdateAwsIamUserSessions(oldSessions, newSessions)
	awsCredentialApplierVerifyExpectedCalls(t, []string{"GetSecret(sessionTokenLabel)"},
		[]string{"GetNamedProfileById(ProfileId)"}, []string{})
}

func TestAwsCredentialsApplier_UpdateAwsImUserSessions_RepositoryWriteCredentialsReturnsError(t *testing.T) {
	awsCredentialApplierSetup()
	oldSessions := []aws_iam_user.AwsIamUserSession{}
	newSessions := []aws_iam_user.AwsIamUserSession{
		{Id: "ID1", Status: domain_aws.Active, NamedProfileId: "ProfileId", SessionTokenLabel: "sessionTokenLabel"},
	}
	awsCredentialApplierNamedProfileFacadeMock.ExpNamedProfile = named_profile.NamedProfile{Id: "ProfileId", Name: "profileName"}
	awsCredentialApplierKeychainMock.ExpGetSecret = "{\"empty\":\"json\"}"

	awsCredentialApplier.UpdateAwsIamUserSessions(oldSessions, newSessions)
	awsCredentialApplierVerifyExpectedCalls(t, []string{"GetSecret(sessionTokenLabel)"},
		[]string{"GetNamedProfileById(ProfileId)"}, []string{"WriteCredentials([{profileName     }])"})
}
