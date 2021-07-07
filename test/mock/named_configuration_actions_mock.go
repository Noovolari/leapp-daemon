package mock

import (
	"errors"
	"fmt"
	"leapp_daemon/domain/domain_gcp/named_configuration"
	"leapp_daemon/infrastructure/http/http_error"
)

type NamedConfigurationsActionsMock struct {
	calls                                   []string
	ExpErrorOnGetNamedConfigurationById     bool
	ExpErrorOnGetOrCreateNamedConfiguration bool
	ExpNamedConfiguration                   named_configuration.NamedConfiguration
	ExpNamedConfigurations                  []named_configuration.NamedConfiguration
}

func NewNamedConfigurationsActionsMock() NamedConfigurationsActionsMock {
	return NamedConfigurationsActionsMock{calls: []string{}}
}

func (actions *NamedConfigurationsActionsMock) GetCalls() []string {
	return actions.calls
}

func (actions *NamedConfigurationsActionsMock) GetNamedConfigurations() []named_configuration.NamedConfiguration {
	actions.calls = append(actions.calls, "GetNamedConfigurations()")
	return actions.ExpNamedConfigurations
}

func (actions *NamedConfigurationsActionsMock) GetNamedConfigurationById(profileId string) (named_configuration.NamedConfiguration, error) {
	actions.calls = append(actions.calls, fmt.Sprintf("GetNamedConfigurationById(%v)", profileId))
	if actions.ExpErrorOnGetNamedConfigurationById {
		return named_configuration.NamedConfiguration{}, http_error.NewNotFoundError(errors.New("named profile not found"))
	}

	return actions.ExpNamedConfiguration, nil
}

func (actions *NamedConfigurationsActionsMock) GetOrCreateNamedConfiguration(configurationName string) (named_configuration.NamedConfiguration, error) {
	actions.calls = append(actions.calls, fmt.Sprintf("GetOrCreateNamedConfiguration(%v)", configurationName))
	if actions.ExpErrorOnGetOrCreateNamedConfiguration {
		return named_configuration.NamedConfiguration{}, http_error.NewBadRequestError(errors.New("configuration name is invalid"))
	}

	return actions.ExpNamedConfiguration, nil
}
