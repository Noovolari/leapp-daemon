package providers

import (
	"leapp_daemon/use_case"
	"sync"
)

var alibabaCredentialsApplierSingleton *use_case.AlibabaCredentialsApplier
var alibabaCredentialsApplierMutex sync.Mutex

func (prov *Providers) GetAlibabaCredentialsApplier() *use_case.AlibabaCredentialsApplier {
	alibabaCredentialsApplierMutex.Lock()
	defer alibabaCredentialsApplierMutex.Unlock()

	if alibabaCredentialsApplierSingleton == nil {
		alibabaCredentialsApplierSingleton = &use_case.AlibabaCredentialsApplier{
			FileSystem:          prov.GetFileSystem(),
			Keychain:            prov.GetKeychain(),
			NamedProfilesFacade: prov.GetNamedProfilesFacade(),
		}
	}
	return alibabaCredentialsApplierSingleton
}
