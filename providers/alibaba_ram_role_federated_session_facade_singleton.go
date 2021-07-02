package providers

import (
	"leapp_daemon/domain/session"
	"sync"
)

var alibabaRamRoleFederatedSessionsFacadeSingleton *session.AlibabaRamRoleFederatedSessionsFacade
var alibabaRamRoleFederatedSessionsFacadeLock sync.Mutex

func (prov *Providers) GetAlibabaRamRoleFederatedSessionFacade() *session.AlibabaRamRoleFederatedSessionsFacade {
	alibabaRamRoleFederatedSessionsFacadeLock.Lock()
	defer alibabaRamRoleFederatedSessionsFacadeLock.Unlock()

	if alibabaRamRoleFederatedSessionsFacadeSingleton == nil {
		alibabaRamRoleFederatedSessionsFacadeSingleton = session.GetAlibabaRamRoleFederatedSessionsFacade()
	}
	return alibabaRamRoleFederatedSessionsFacadeSingleton
}
