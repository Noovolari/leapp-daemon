package providers

import (
	"leapp_daemon/use_case"
	"sync"
)

var alibabaRamRoleChainedSessionActionsSingleton *use_case.AlibabaRamRoleChainedSessionActions
var alibabaRamRoleChainedSessionActionsMutex sync.Mutex

func (prov *Providers) GetAlibabaRamRoleChainedSessionActions() *use_case.AlibabaRamRoleChainedSessionActions {
	alibabaRamRoleChainedSessionActionsMutex.Lock()
	defer alibabaRamRoleChainedSessionActionsMutex.Unlock()

	if alibabaRamRoleChainedSessionActionsSingleton == nil {
		alibabaRamRoleChainedSessionActionsSingleton = &use_case.AlibabaRamRoleChainedSessionActions{
			NamedProfilesActions:                prov.GetNamedProfilesActions(),
			Environment:                         prov.GetEnvironment(),
			Keychain:                            prov.GetKeychain(),
			AlibabaRamRoleChainedSessionsFacade: prov.GetAlibabaRamRoleChainedSessionFacade(),
		}
	}
	return alibabaRamRoleChainedSessionActionsSingleton
}
