package use_case

import (
	"leapp_daemon/domain/domain_aws/aws_iam_user"
	"leapp_daemon/test/mock"
	"reflect"
	"testing"
)

var (
	awsOldSessions    []aws_iam_user.AwsIamUserSession
	awsNewSessions    []aws_iam_user.AwsIamUserSession
	awsFileRepoMock   mock.FileConfigurationRepositoryMock
	awsSessionsWriter *AwsSessionsWriter
)

func awsSessionsWriterSetup() {
	awsOldSessions = []aws_iam_user.AwsIamUserSession{}
	awsNewSessions = []aws_iam_user.AwsIamUserSession{{ID: "ID"}}

	awsFileRepoMock = mock.NewFileConfigurationRepositoryMock()
	awsSessionsWriter = &AwsSessionsWriter{
		ConfigurationRepository: &awsFileRepoMock,
	}
}

func awsSessionsWriterVerifyExpectedCalls(t *testing.T, fileRepoMockCalls []string) {
	if !reflect.DeepEqual(awsFileRepoMock.GetCalls(), fileRepoMockCalls) {
		t.Fatalf("FileRepoMock expectation violation.\nMock calls: %v", awsFileRepoMock.GetCalls())
	}
}

func TestUpdateAwsIamUserSessions(t *testing.T) {
	awsSessionsWriterSetup()

	awsSessionsWriter.UpdateAwsIamUserSessions(awsOldSessions, awsNewSessions)
	awsSessionsWriterVerifyExpectedCalls(t, []string{"GetConfiguration()", "UpdateConfiguration()"})
}

func TestUpdateAwsIamUserSessions_ErrorGettingConfiguration(t *testing.T) {
	awsSessionsWriterSetup()
	awsFileRepoMock.ExpErrorOnGetConfiguration = true

	awsSessionsWriter.UpdateAwsIamUserSessions(awsOldSessions, awsNewSessions)
	awsSessionsWriterVerifyExpectedCalls(t, []string{"GetConfiguration()"})
}

func TestUpdateAwsIamUserSessions_ErrorUpdatingConfiguration(t *testing.T) {
	awsSessionsWriterSetup()
	awsFileRepoMock.ExpErrorOnUpdateConfiguration = true

	awsSessionsWriter.UpdateAwsIamUserSessions(awsOldSessions, awsNewSessions)
	awsSessionsWriterVerifyExpectedCalls(t, []string{"GetConfiguration()", "UpdateConfiguration()"})
}
