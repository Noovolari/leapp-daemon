package providers

import (
	"leapp_daemon/use_case"
	"sync"
)

var gcpIamUserAccountOauthSessionActionsSingleton *use_case.GcpIamUserAccountOauthSessionActions
var gcpIamUserAccountOauthSessionActionsMutex sync.Mutex

func (prov *Providers) GetGcpIamUserAccountOauthSessionActions() *use_case.GcpIamUserAccountOauthSessionActions {
	gcpIamUserAccountOauthSessionActionsMutex.Lock()
	defer gcpIamUserAccountOauthSessionActionsMutex.Unlock()

	if gcpIamUserAccountOauthSessionActionsSingleton == nil {
		gcpIamUserAccountOauthSessionActionsSingleton = &use_case.GcpIamUserAccountOauthSessionActions{
			GcpApi:                              prov.GetGcpApi(),
			Environment:                         prov.GetEnvironment(),
			Keychain:                            prov.GetKeychain(),
			GcpIamUserAccountOauthSessionFacade: prov.GetGcpIamUserAccountOauthSessionFacade(),
		}
	}
	return gcpIamUserAccountOauthSessionActionsSingleton
}
