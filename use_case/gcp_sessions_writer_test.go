package use_case

import (
	"leapp_daemon/domain/domain_gcp/gcp_iam_user_account_oauth"
	"leapp_daemon/test/mock"
	"reflect"
	"testing"
)

var (
	gcpOldSessions    []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession
	gcpNewSessions    []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession
	gcpFileRepoMock   mock.FileConfigurationRepositoryMock
	gcpSessionsWriter *GcpSessionsWriter
)

func gcpSessionsWriterSetup() {
	gcpOldSessions = []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{}
	gcpNewSessions = []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession{{Id: "ID"}}

	gcpFileRepoMock = mock.NewFileConfigurationRepositoryMock()
	gcpSessionsWriter = &GcpSessionsWriter{
		ConfigurationRepository: &gcpFileRepoMock,
	}
}

func gcpSessionsWriterVerifyExpectedCalls(t *testing.T, fileRepoMockCalls []string) {
	if !reflect.DeepEqual(gcpFileRepoMock.GetCalls(), fileRepoMockCalls) {
		t.Fatalf("FileRepoMock expectation violation.\nMock calls: %v", gcpFileRepoMock.GetCalls())
	}
}

func TestUpdateGcpIamUserAccountOauthSessions(t *testing.T) {
	gcpSessionsWriterSetup()

	gcpSessionsWriter.UpdateGcpIamUserAccountOauthSessions(gcpOldSessions, gcpNewSessions)
	gcpSessionsWriterVerifyExpectedCalls(t, []string{"GetConfiguration()", "UpdateConfiguration()"})
}

func TestUpdateGcpIamUserAccountOauthSessions_ErrorGettingConfiguration(t *testing.T) {
	gcpSessionsWriterSetup()
	gcpFileRepoMock.ExpErrorOnGetConfiguration = true

	gcpSessionsWriter.UpdateGcpIamUserAccountOauthSessions(gcpOldSessions, gcpNewSessions)
	gcpSessionsWriterVerifyExpectedCalls(t, []string{"GetConfiguration()"})
}

func TestUpdateGcpIamUserAccountOauthSessions_ErrorUpdatingConfiguration(t *testing.T) {
	gcpSessionsWriterSetup()
	gcpFileRepoMock.ExpErrorOnUpdateConfiguration = true

	gcpSessionsWriter.UpdateGcpIamUserAccountOauthSessions(gcpOldSessions, gcpNewSessions)
	gcpSessionsWriterVerifyExpectedCalls(t, []string{"GetConfiguration()", "UpdateConfiguration()"})
}
