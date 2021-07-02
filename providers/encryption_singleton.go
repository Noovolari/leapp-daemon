package providers

import (
	"leapp_daemon/infrastructure/encryption"
	"sync"
)

var encryptionSingleton *encryption.Encryption
var encryptionMutex sync.Mutex

func (prov *Providers) GetEncryption() *encryption.Encryption {
	encryptionMutex.Lock()
	defer encryptionMutex.Unlock()

	if encryptionSingleton == nil {
		encryptionSingleton = &encryption.Encryption{}
	}
	return encryptionSingleton
}
