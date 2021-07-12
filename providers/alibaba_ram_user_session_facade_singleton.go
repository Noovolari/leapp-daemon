package providers

import (
	"leapp_daemon/domain/domain_alibaba/alibaba_ram_user"
	"sync"
)

var alibabaRamUserSessionsFacadeSingleton *alibaba_ram_user.AlibabaRamUserSessionsFacade
var alibabaRamUserSessionsFacadeLock sync.Mutex

func (prov *Providers) GetAlibabaRamUserSessionFacade() *alibaba_ram_user.AlibabaRamUserSessionsFacade {
	alibabaRamUserSessionsFacadeLock.Lock()
	defer alibabaRamUserSessionsFacadeLock.Unlock()

	if alibabaRamUserSessionsFacadeSingleton == nil {
		alibabaRamUserSessionsFacadeSingleton = alibaba_ram_user.GetAlibabaRamUserSessionsFacade()
	}
	return alibabaRamUserSessionsFacadeSingleton
}
