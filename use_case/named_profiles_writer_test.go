package use_case

import (
	"leapp_daemon/domain/domain_aws/named_profile"
	"leapp_daemon/test/mock"
	"reflect"
	"testing"
)

var (
	oldNamedProfiles         []named_profile.NamedProfile
	newNamedProfiles         []named_profile.NamedProfile
	namedProfileFileRepoMock mock.FileConfigurationRepositoryMock
	namedProfilesWriter      *NamedProfilesWriter
)

func namedProfilesWriterSetup() {
	oldNamedProfiles = []named_profile.NamedProfile{}
	newNamedProfiles = []named_profile.NamedProfile{}

	namedProfileFileRepoMock = mock.NewFileConfigurationRepositoryMock()
	namedProfilesWriter = &NamedProfilesWriter{
		ConfigurationRepository: &namedProfileFileRepoMock,
	}
}

func namedProfilesWriterVerifyExpectedCalls(t *testing.T, fileRepoMockCalls []string) {
	if !reflect.DeepEqual(namedProfileFileRepoMock.GetCalls(), fileRepoMockCalls) {
		t.Fatalf("FileRepoMock expectation violation.\nMock calls: %v", namedProfileFileRepoMock.GetCalls())
	}
}

func TestUpdateNamedProfiles(t *testing.T) {
	namedProfilesWriterSetup()

	namedProfilesWriter.UpdateNamedProfiles(oldNamedProfiles, newNamedProfiles)
	namedProfilesWriterVerifyExpectedCalls(t, []string{"GetConfiguration()", "UpdateConfiguration()"})
}

func TestUpdateNamedProfiles_ErrorGettingProfile(t *testing.T) {
	namedProfilesWriterSetup()
	namedProfileFileRepoMock.ExpErrorOnGetConfiguration = true

	namedProfilesWriter.UpdateNamedProfiles(oldNamedProfiles, newNamedProfiles)
	namedProfilesWriterVerifyExpectedCalls(t, []string{"GetConfiguration()"})
}

func TestUpdateNamedProfiles_ErrorUpdatingProfile(t *testing.T) {
	namedProfilesWriterSetup()
	namedProfileFileRepoMock.ExpErrorOnUpdateConfiguration = true

	namedProfilesWriter.UpdateNamedProfiles(oldNamedProfiles, newNamedProfiles)
	namedProfilesWriterVerifyExpectedCalls(t, []string{"GetConfiguration()", "UpdateConfiguration()"})
}
