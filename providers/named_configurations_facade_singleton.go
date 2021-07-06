package providers

import (
	"leapp_daemon/domain/domain_gcp/named_configuration"
	"sync"
)

var namedConfigurationFacadeSingleton *named_configuration.NamedConfigurationsFacade
var namedConfigurationFacadeLock sync.Mutex

func (prov *Providers) GetNamedConfigurationFacade() *named_configuration.NamedConfigurationsFacade {
	namedConfigurationFacadeLock.Lock()
	defer namedConfigurationFacadeLock.Unlock()

	if namedConfigurationFacadeSingleton == nil {
		namedConfigurationFacadeSingleton = named_configuration.NewNamedConfigurationsFacade()
	}
	return namedConfigurationFacadeSingleton
}
