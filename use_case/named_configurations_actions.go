package use_case

import (
	"leapp_daemon/domain/domain_gcp/named_configuration"
)

type NamedConfigurationsActions struct {
	Environment               Environment
	NamedConfigurationsFacade NamedConfigurationsFacade
}

func (actions *NamedConfigurationsActions) GetNamedConfigurations() []named_configuration.NamedConfiguration {
	return actions.NamedConfigurationsFacade.GetNamedConfigurations()
}

func (actions *NamedConfigurationsActions) GetNamedConfigurationById(profileId string) (named_configuration.NamedConfiguration, error) {
	return actions.NamedConfigurationsFacade.GetNamedConfigurationById(profileId)
}

func (actions *NamedConfigurationsActions) GetOrCreateNamedConfiguration(configurationName string) (named_configuration.NamedConfiguration, error) {
	if configurationName == "" {
		configurationName = named_configuration.DefaultNamedConfigurationName
	}

	facade := actions.NamedConfigurationsFacade
	namedConfiguration, err := facade.GetNamedConfigurationByName(configurationName)
	if err != nil {
		namedConfiguration = named_configuration.NamedConfiguration{
			Id:   actions.Environment.GenerateUuid(),
			Name: configurationName,
		}
		err = facade.AddNamedConfiguration(namedConfiguration)
		if err != nil {
			return named_configuration.NamedConfiguration{}, err
		}
	}

	return namedConfiguration, nil
}
