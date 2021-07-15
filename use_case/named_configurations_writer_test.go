package use_case

import (
	"leapp_daemon/domain/domain_gcp/named_configuration"
	"leapp_daemon/test/mock"
	"reflect"
	"testing"
)

var (
	oldNamedConfigurations         []named_configuration.NamedConfiguration
	newNamedConfigurations         []named_configuration.NamedConfiguration
	namedConfigurationFileRepoMock mock.FileConfigurationRepositoryMock
	namedConfigurationsWriter      *NamedConfigurationsWriter
)

func namedConfigurationsWriterSetup() {
	oldNamedConfigurations = []named_configuration.NamedConfiguration{}
	newNamedConfigurations = []named_configuration.NamedConfiguration{}

	namedConfigurationFileRepoMock = mock.NewFileConfigurationRepositoryMock()
	namedConfigurationsWriter = &NamedConfigurationsWriter{
		ConfigurationRepository: &namedConfigurationFileRepoMock,
	}
}

func namedConfigurationsWriterVerifyExpectedCalls(t *testing.T, fileRepoMockCalls []string) {
	if !reflect.DeepEqual(namedConfigurationFileRepoMock.GetCalls(), fileRepoMockCalls) {
		t.Fatalf("FileRepoMock expectation violation.\nMock calls: %v", namedConfigurationFileRepoMock.GetCalls())
	}
}

func TestUpdateNamedConfigurations(t *testing.T) {
	namedConfigurationsWriterSetup()

	namedConfigurationsWriter.UpdateNamedConfigurations(oldNamedConfigurations, newNamedConfigurations)
	namedConfigurationsWriterVerifyExpectedCalls(t, []string{"GetConfiguration()", "UpdateConfiguration()"})
}

func TestUpdateNamedConfigurations_ErrorGettingConfiguration(t *testing.T) {
	namedConfigurationsWriterSetup()
	namedConfigurationFileRepoMock.ExpErrorOnGetConfiguration = true

	namedConfigurationsWriter.UpdateNamedConfigurations(oldNamedConfigurations, newNamedConfigurations)
	namedConfigurationsWriterVerifyExpectedCalls(t, []string{"GetConfiguration()"})
}

func TestUpdateNamedConfigurations_ErrorUpdatingConfiguration(t *testing.T) {
	namedConfigurationsWriterSetup()
	namedConfigurationFileRepoMock.ExpErrorOnUpdateConfiguration = true

	namedConfigurationsWriter.UpdateNamedConfigurations(oldNamedConfigurations, newNamedConfigurations)
	namedConfigurationsWriterVerifyExpectedCalls(t, []string{"GetConfiguration()", "UpdateConfiguration()"})
}
