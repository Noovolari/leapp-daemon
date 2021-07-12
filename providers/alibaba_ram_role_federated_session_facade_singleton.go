package providers

import (
	"leapp_daemon/domain/domain_alibaba/alibaba_ram_role_federated"
	"sync"
)

var alibabaRamRoleFederatedSessionsFacadeSingleton *alibaba_ram_role_federated.AlibabaRamRoleFederatedSessionsFacade
var alibabaRamRoleFederatedSessionsFacadeLock sync.Mutex

func (prov *Providers) GetAlibabaRamRoleFederatedSessionFacade() *alibaba_ram_role_federated.AlibabaRamRoleFederatedSessionsFacade {
	alibabaRamRoleFederatedSessionsFacadeLock.Lock()
	defer alibabaRamRoleFederatedSessionsFacadeLock.Unlock()

	if alibabaRamRoleFederatedSessionsFacadeSingleton == nil {
		alibabaRamRoleFederatedSessionsFacadeSingleton = alibaba_ram_role_federated.GetAlibabaRamRoleFederatedSessionsFacade()
	}
	return alibabaRamRoleFederatedSessionsFacadeSingleton
}
