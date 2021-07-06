package use_case

import (
	"leapp_daemon/domain/domain_gcp/named_configuration"
)

type NamedConfigurationsWriter struct {
	ConfigurationRepository ConfigurationRepository
}

func (writer *NamedConfigurationsWriter) UpdateNamedConfigurations(oldNamedConfigurations []named_configuration.NamedConfiguration,
	newNamedConfigurations []named_configuration.NamedConfiguration) error {
	config, err := writer.ConfigurationRepository.GetConfiguration()
	if err != nil {
		return err
	}

	config.NamedConfigurations = newNamedConfigurations
	err = writer.ConfigurationRepository.UpdateConfiguration(config)

	return err
}
