package providers

import (
	"leapp_daemon/adapter/repository"
	"sync"
)

var fileConfigurationRepositorySingleton *repository.FileConfigurationRepository
var fileConfigurationRepositoryLock sync.Mutex

func (prov *Providers) GetFileConfigurationRepository() *repository.FileConfigurationRepository {
	fileConfigurationRepositoryLock.Lock()
	defer fileConfigurationRepositoryLock.Unlock()

	if fileConfigurationRepositorySingleton == nil {
		fileConfigurationRepositorySingleton = &repository.FileConfigurationRepository{
			FileSystem: prov.GetFileSystem(),
			Encryption: prov.GetEncryption(),
		}
	}
	return fileConfigurationRepositorySingleton
}
