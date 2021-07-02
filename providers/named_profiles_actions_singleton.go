package providers

import (
	"leapp_daemon/use_case"
	"sync"
)

var namedProfileActionsSingleton *use_case.NamedProfilesActions
var namedProfileMutex sync.Mutex

func (prov *Providers) GetNamedProfilesActions() *use_case.NamedProfilesActions {
	namedProfileMutex.Lock()
	defer namedProfileMutex.Unlock()

	if namedProfileActionsSingleton == nil {
		namedProfileActionsSingleton = &use_case.NamedProfilesActions{
			Environment:         prov.GetEnvironment(),
			NamedProfilesFacade: prov.GetNamedProfilesFacade(),
		}
	}
	return namedProfileActionsSingleton
}
