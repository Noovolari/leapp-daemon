package providers

import (
	"leapp_daemon/interface/repository"
	"sync"
)

var gcpConfigurationRepositorySingleton *repository.GcpConfigurationRepository
var gcpConfigurationRepositoryMutex sync.Mutex

func (prov *Providers) GetGcpConfigurationRepository() *repository.GcpConfigurationRepository {
	gcpConfigurationRepositoryMutex.Lock()
	defer gcpConfigurationRepositoryMutex.Unlock()

	if gcpConfigurationRepositorySingleton == nil {
		gcpConfigurationRepositorySingleton = &repository.GcpConfigurationRepository{
			FileSystem:        prov.GetFileSystem(),
			Environment:       prov.GetEnvironment(),
			CredentialsTable:  prov.GetGcpCredentialsTable(),
			AccessTokensTable: prov.GetGcpAccessTokensTable(),
		}
	}
	return gcpConfigurationRepositorySingleton
}
