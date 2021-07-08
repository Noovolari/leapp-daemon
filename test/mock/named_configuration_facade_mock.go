package mock

import (
	"errors"
	"fmt"
	"leapp_daemon/domain/domain_gcp/named_configuration"
	"leapp_daemon/infrastructure/http/http_error"
)

type NamedConfigurationsFacadeMock struct {
	calls                                 []string
	ExpErrorOnGetNamedConfigurationById   bool
	ExpErrorOnGetNamedConfigurationByName bool
	ExpErrorOnAddNamedConfiguration       bool
	ExpNamedConfiguration                 named_configuration.NamedConfiguration
	ExpNamedConfigurations                []named_configuration.NamedConfiguration
}

func NewNamedConfigurationsFacadeMock() NamedConfigurationsFacadeMock {
	return NamedConfigurationsFacadeMock{calls: []string{}}
}

func (actions *NamedConfigurationsFacadeMock) GetCalls() []string {
	return actions.calls
}

func (actions *NamedConfigurationsFacadeMock) GetNamedConfigurations() []named_configuration.NamedConfiguration {
	actions.calls = append(actions.calls, "GetNamedConfigurations()")
	return actions.ExpNamedConfigurations
}

func (actions *NamedConfigurationsFacadeMock) GetNamedConfigurationById(configurationId string) (named_configuration.NamedConfiguration, error) {
	actions.calls = append(actions.calls, fmt.Sprintf("GetNamedConfigurationById(%v)", configurationId))
	if actions.ExpErrorOnGetNamedConfigurationById {
		return named_configuration.NamedConfiguration{}, http_error.NewNotFoundError(errors.New("named configuration not found"))
	}

	return actions.ExpNamedConfiguration, nil
}

func (actions *NamedConfigurationsFacadeMock) GetNamedConfigurationByName(configurationName string) (named_configuration.NamedConfiguration, error) {
	actions.calls = append(actions.calls, fmt.Sprintf("GetNamedConfigurationByName(%v)", configurationName))
	if actions.ExpErrorOnGetNamedConfigurationByName {
		return named_configuration.NamedConfiguration{}, http_error.NewNotFoundError(errors.New("named configuration not found"))
	}

	return actions.ExpNamedConfiguration, nil
}

func (actions *NamedConfigurationsFacadeMock) AddNamedConfiguration(namedConfiguration named_configuration.NamedConfiguration) error {
	actions.calls = append(actions.calls, fmt.Sprintf("AddNamedConfiguration(%v)", namedConfiguration))
	if actions.ExpErrorOnAddNamedConfiguration {
		return http_error.NewConflictError(errors.New("unable to add named configuration"))
	}

	return nil
}
