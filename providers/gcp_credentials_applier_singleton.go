package providers

import (
	"leapp_daemon/use_case"
	"sync"
)

var gcpCredentialsApplierSingleton *use_case.GcpCredentialsApplier
var gcpCredentialsApplierMutex sync.Mutex

func (prov *Providers) GetGcpCredentialsApplier() *use_case.GcpCredentialsApplier {
	gcpCredentialsApplierMutex.Lock()
	defer gcpCredentialsApplierMutex.Unlock()

	if gcpCredentialsApplierSingleton == nil {
		gcpCredentialsApplierSingleton = &use_case.GcpCredentialsApplier{
			Repository:                prov.GetGcpConfigurationRepository(),
			Keychain:                  prov.GetKeychain(),
			NamedConfigurationsFacade: prov.GetNamedConfigurationFacade(),
		}
	}
	return gcpCredentialsApplierSingleton
}
