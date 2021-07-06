package named_configuration

import (
	"fmt"
	"leapp_daemon/infrastructure/http/http_error"
	"regexp"
	"sync"
)

var namedConfigurationsLock sync.Mutex

type NamedConfigurationsObserver interface {
	UpdateNamedConfigurations(oldNamedConfigurations []NamedConfiguration, newNamedConfigurations []NamedConfiguration) error
}

type NamedConfigurationsFacade struct {
	namedConfigurations []NamedConfiguration
	observers           []NamedConfigurationsObserver
}

func NewNamedConfigurationsFacade() *NamedConfigurationsFacade {
	return &NamedConfigurationsFacade{
		namedConfigurations: make([]NamedConfiguration, 0),
	}
}

func (fac *NamedConfigurationsFacade) Subscribe(observer NamedConfigurationsObserver) {
	fac.observers = append(fac.observers, observer)
}

func (fac *NamedConfigurationsFacade) GetNamedConfigurations() []NamedConfiguration {
	return fac.namedConfigurations
}

func (fac *NamedConfigurationsFacade) SetNamedConfigurations(namedConfigurations []NamedConfiguration) {
	fac.namedConfigurations = namedConfigurations
}

func (fac *NamedConfigurationsFacade) AddNamedConfiguration(namedConfiguration NamedConfiguration) error {
	namedConfigurationsLock.Lock()
	defer namedConfigurationsLock.Unlock()

	err := fac.checkConfigurationNameValidity(namedConfiguration.Name)
	if err != nil {
		return err
	}

	namedConfigurations := fac.GetNamedConfigurations()
	for _, np := range namedConfigurations {
		if namedConfiguration.Id == np.Id {
			return http_error.NewConflictError(fmt.Errorf("a ConfigurationName with id " + namedConfiguration.Id +
				" is already present"))
		}
		if namedConfiguration.Name == np.Name {
			return http_error.NewConflictError(fmt.Errorf("a ConfigurationName with name " + namedConfiguration.Name +
				" is already present"))
		}
	}

	namedConfigurations = append(namedConfigurations, namedConfiguration)
	err = fac.updateState(namedConfigurations)
	if err != nil {
		return err
	}

	return nil
}

func (fac *NamedConfigurationsFacade) GetNamedConfigurationByName(name string) (NamedConfiguration, error) {
	for _, namedConfiguration := range fac.namedConfigurations {
		if namedConfiguration.Name == name {
			return namedConfiguration, nil
		}
	}
	return NamedConfiguration{}, http_error.NewNotFoundError(fmt.Errorf("named configuration with name %v not found", name))
}

func (fac *NamedConfigurationsFacade) GetNamedConfigurationById(id string) (NamedConfiguration, error) {
	for _, namedConfiguration := range fac.namedConfigurations {
		if namedConfiguration.Id == id {
			return namedConfiguration, nil
		}
	}
	return NamedConfiguration{}, http_error.NewNotFoundError(fmt.Errorf("named configuration with id %v not found", id))
}

func (fac *NamedConfigurationsFacade) updateState(newState []NamedConfiguration) error {
	oldNamedConfigurations := fac.GetNamedConfigurations()
	fac.SetNamedConfigurations(newState)

	for _, observer := range fac.observers {
		err := observer.UpdateNamedConfigurations(oldNamedConfigurations, newState)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fac *NamedConfigurationsFacade) checkConfigurationNameValidity(configurationName string) error {
	isValid, _ := regexp.MatchString("^[a-z][a-z0-9-]*$", configurationName)
	if !isValid {
		return http_error.NewBadRequestError(fmt.Errorf("configuration names must start with a lower case " +
			"letter and contain only lower case letters a-z, digits 0-9, and hyphens '-'"))
	}
	return nil
}
