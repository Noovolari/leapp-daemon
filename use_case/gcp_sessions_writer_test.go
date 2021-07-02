package use_case

import (
	"leapp_daemon/domain/session"
	"leapp_daemon/test/mock"
	"reflect"
	"testing"
)

var (
	oldSessions       []session.GcpIamUserAccountOauthSession
	newSessions       []session.GcpIamUserAccountOauthSession
	fileRepoMock      mock.FileConfigurationRepositoryMock
	gcpSessionsWriter *GcpSessionsWriter
)

func gcpSessionsWriterSetup() {
	oldSessions = []session.GcpIamUserAccountOauthSession{}
	newSessions = []session.GcpIamUserAccountOauthSession{{Id: "ID"}}

	fileRepoMock = mock.NewFileConfigurationRepositoryMock()
	gcpSessionsWriter = &GcpSessionsWriter{
		ConfigurationRepository: &fileRepoMock,
	}
}

func gcpSessionsWriterVerifyExpectedCalls(t *testing.T, fileRepoMockCalls []string) {
	if !reflect.DeepEqual(fileRepoMock.GetCalls(), fileRepoMockCalls) {
		t.Fatalf("fileRepoMock expectation violation.\nMock calls: %v", fileRepoMock.GetCalls())
	}
}

func TestUpdateGcpIamUserAccountOauthSessions(t *testing.T) {
	gcpSessionsWriterSetup()

	gcpSessionsWriter.UpdateGcpIamUserAccountOauthSessions(oldSessions, newSessions)
	gcpSessionsWriterVerifyExpectedCalls(t, []string{"GetConfiguration()", "UpdateConfiguration()"})
}

func TestUpdateGcpIamUserAccountOauthSessions_ErrorGettingConfiguration(t *testing.T) {
	gcpSessionsWriterSetup()
	fileRepoMock.ExpErrorOnGetConfiguration = true

	gcpSessionsWriter.UpdateGcpIamUserAccountOauthSessions(oldSessions, newSessions)
	gcpSessionsWriterVerifyExpectedCalls(t, []string{"GetConfiguration()"})
}

func TestUpdateGcpIamUserAccountOauthSessions_ErrorUpdatingConfiguration(t *testing.T) {
	gcpSessionsWriterSetup()
	fileRepoMock.ExpErrorOnUpdateConfiguration = true

	gcpSessionsWriter.UpdateGcpIamUserAccountOauthSessions(oldSessions, newSessions)
	gcpSessionsWriterVerifyExpectedCalls(t, []string{"GetConfiguration()", "UpdateConfiguration()"})
}
