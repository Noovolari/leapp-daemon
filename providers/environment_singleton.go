package providers

import (
	"leapp_daemon/infrastructure/environment"
	"sync"
)

var environmentSingleton *environment.Environment
var environmentMutex sync.Mutex

func (prov *Providers) GetEnvironment() *environment.Environment {
	environmentMutex.Lock()
	defer environmentMutex.Unlock()

	if environmentSingleton == nil {
		environmentSingleton = &environment.Environment{}
	}
	return environmentSingleton
}
