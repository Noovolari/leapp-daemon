package providers

import (
	"leapp_daemon/domain/session"
	"sync"
)

var alibabaRamUserSessionsFacadeSingleton *session.AlibabaRamUserSessionsFacade
var alibabaRamUserSessionsFacadeLock sync.Mutex

func (prov *Providers) GetAlibabaRamUserSessionFacade() *session.AlibabaRamUserSessionsFacade {
	alibabaRamUserSessionsFacadeLock.Lock()
	defer alibabaRamUserSessionsFacadeLock.Unlock()

	if alibabaRamUserSessionsFacadeSingleton == nil {
		alibabaRamUserSessionsFacadeSingleton = session.GetAlibabaRamUserSessionsFacade()
	}
	return alibabaRamUserSessionsFacadeSingleton
}
