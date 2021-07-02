package providers

import (
	"leapp_daemon/infrastructure/keychain"
	"sync"
)

var keychainSingleton *keychain.Keychain
var keychainMutex sync.Mutex

func (prov *Providers) GetKeychain() *keychain.Keychain {
	keychainMutex.Lock()
	defer keychainMutex.Unlock()

	if keychainSingleton == nil {
		keychainSingleton = &keychain.Keychain{}
	}
	return keychainSingleton
}
