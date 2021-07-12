package providers

import (
	"leapp_daemon/domain/domain_alibaba/alibaba_ram_role_chained"
	"sync"
)

var alibabaRamRoleChainedSessionsFacadeSingleton *alibaba_ram_role_chained.AlibabaRamRoleChainedSessionsFacade
var alibabaRamRoleChainedSessionsFacadeLock sync.Mutex

func (prov *Providers) GetAlibabaRamRoleChainedSessionFacade() *alibaba_ram_role_chained.AlibabaRamRoleChainedSessionsFacade {
	alibabaRamRoleChainedSessionsFacadeLock.Lock()
	defer alibabaRamRoleChainedSessionsFacadeLock.Unlock()

	if alibabaRamRoleChainedSessionsFacadeSingleton == nil {
		alibabaRamRoleChainedSessionsFacadeSingleton = alibaba_ram_role_chained.GetAlibabaRamRoleChainedSessionsFacade()
	}
	return alibabaRamRoleChainedSessionsFacadeSingleton
}
