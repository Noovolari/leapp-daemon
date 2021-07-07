package providers

import (
	"leapp_daemon/domain/domain_aws/named_profile"
	"sync"
)

var namedProfilesFacadeSingleton *named_profile.NamedProfilesFacade
var namedProfilesFacadeLock sync.Mutex

func (prov *Providers) GetNamedProfilesFacade() *named_profile.NamedProfilesFacade {
	namedProfilesFacadeLock.Lock()
	defer namedProfilesFacadeLock.Unlock()

	if namedProfilesFacadeSingleton == nil {
		namedProfilesFacadeSingleton = named_profile.NewNamedProfilesFacade()
	}
	return namedProfilesFacadeSingleton
}
