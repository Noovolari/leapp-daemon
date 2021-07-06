package providers

import (
	"leapp_daemon/use_case"
	"sync"
)

var alibabaRamUserSessionActionsSingleton *use_case.AlibabaRamUserSessionActions
var alibabaRamUserSessionActionsMutex sync.Mutex

func (prov *Providers) GetAlibabaRamUserSessionActions() *use_case.AlibabaRamUserSessionActions {
	alibabaRamUserSessionActionsMutex.Lock()
	defer alibabaRamUserSessionActionsMutex.Unlock()

	if alibabaRamUserSessionActionsSingleton == nil {
		alibabaRamUserSessionActionsSingleton = &use_case.AlibabaRamUserSessionActions{
			NamedProfilesActions:         prov.GetNamedProfilesActions(),
			Environment:                  prov.GetEnvironment(),
			Keychain:                     prov.GetKeychain(),
			AlibabaRamUserSessionsFacade: prov.GetAlibabaRamUserSessionFacade(),
		}
	}
	return alibabaRamUserSessionActionsSingleton
}
