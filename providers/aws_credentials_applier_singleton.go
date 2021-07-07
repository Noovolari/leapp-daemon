package providers

import (
	"leapp_daemon/use_case"
	"sync"
)

var awsCredentialsApplierSingleton *use_case.AwsCredentialsApplier
var awsCredentialsApplierMutex sync.Mutex

func (prov *Providers) GetAwsCredentialsApplier() *use_case.AwsCredentialsApplier {
	awsCredentialsApplierMutex.Lock()
	defer awsCredentialsApplierMutex.Unlock()

	if awsCredentialsApplierSingleton == nil {
		awsCredentialsApplierSingleton = &use_case.AwsCredentialsApplier{
			Keychain:            prov.GetKeychain(),
			NamedProfilesFacade: prov.GetNamedProfilesFacade(),
			Repository:          prov.GetAwsConfigurationRepository(),
		}
	}
	return awsCredentialsApplierSingleton
}
