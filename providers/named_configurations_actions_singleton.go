package providers

import (
	"leapp_daemon/use_case"
	"sync"
)

var namedConfigurationActionsSingleton *use_case.NamedConfigurationsActions
var namedConfigurationMutex sync.Mutex

func (prov *Providers) GetNamedConfigurationsActions() *use_case.NamedConfigurationsActions {
	namedConfigurationMutex.Lock()
	defer namedConfigurationMutex.Unlock()

	if namedConfigurationActionsSingleton == nil {
		namedConfigurationActionsSingleton = &use_case.NamedConfigurationsActions{
			Environment:               prov.GetEnvironment(),
			NamedConfigurationsFacade: prov.GetNamedConfigurationFacade(),
		}
	}
	return namedConfigurationActionsSingleton
}
