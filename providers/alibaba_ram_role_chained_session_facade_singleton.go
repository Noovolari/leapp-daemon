package providers

import (
	"leapp_daemon/domain/session"
	"sync"
)

var alibabaRamRoleChainedSessionsFacadeSingleton *session.AlibabaRamRoleChainedSessionsFacade
var alibabaRamRoleChainedSessionsFacadeLock sync.Mutex

func (prov *Providers) GetAlibabaRamRoleChainedSessionFacade() *session.AlibabaRamRoleChainedSessionsFacade {
	alibabaRamRoleChainedSessionsFacadeLock.Lock()
	defer alibabaRamRoleChainedSessionsFacadeLock.Unlock()

	if alibabaRamRoleChainedSessionsFacadeSingleton == nil {
		alibabaRamRoleChainedSessionsFacadeSingleton = session.GetAlibabaRamRoleChainedSessionsFacade()
	}
	return alibabaRamRoleChainedSessionsFacadeSingleton
}
