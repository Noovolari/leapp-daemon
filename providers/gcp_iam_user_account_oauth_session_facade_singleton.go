package providers

import (
	"leapp_daemon/domain/session"
	"sync"
)

var gcpSessionsFacadeSingleton *session.GcpIamUserAccountOauthSessionsFacade
var gcpSessionsFacadeLock sync.Mutex

func (prov *Providers) GetGcpIamUserAccountOauthSessionFacade() *session.GcpIamUserAccountOauthSessionsFacade {
	gcpSessionsFacadeLock.Lock()
	defer gcpSessionsFacadeLock.Unlock()

	if gcpSessionsFacadeSingleton == nil {
		gcpSessionsFacadeSingleton = session.NewGcpIamUserAccountOauthSessionsFacade()
	}
	return gcpSessionsFacadeSingleton
}
