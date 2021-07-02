package providers

import (
	"leapp_daemon/use_case"
	"sync"
)

var alibabaRamRoleFederatedSessionActionsSingleton *use_case.AlibabaRamRoleFederatedSessionActions
var alibabaRamRoleFederatedSessionActionsMutex sync.Mutex

func (prov *Providers) GetAlibabaRamRoleFederatedSessionActions() *use_case.AlibabaRamRoleFederatedSessionActions {
	alibabaRamRoleFederatedSessionActionsMutex.Lock()
	defer alibabaRamRoleFederatedSessionActionsMutex.Unlock()

	if alibabaRamRoleFederatedSessionActionsSingleton == nil {
		alibabaRamRoleFederatedSessionActionsSingleton = &use_case.AlibabaRamRoleFederatedSessionActions{
			NamedProfilesActions:     prov.GetNamedProfilesActions(),
			Environment:              prov.GetEnvironment(),
			Keychain:                 prov.GetKeychain(),
			/*AlibabaRamRoleFederatedSessionsFacade: prov.GetAlibabaRamRoleFederatedSessionFacade(),*/
		}
	}
	return alibabaRamRoleFederatedSessionActionsSingleton
}
