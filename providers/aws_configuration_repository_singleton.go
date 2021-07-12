package providers

import (
	"leapp_daemon/adapter/repository"
	"sync"
)

var awsConfigurationRepositorySingleton *repository.AwsConfigurationRepository
var awsConfigurationRepositoryMutex sync.Mutex

func (prov *Providers) GetAwsConfigurationRepository() *repository.AwsConfigurationRepository {
	awsConfigurationRepositoryMutex.Lock()
	defer awsConfigurationRepositoryMutex.Unlock()

	if awsConfigurationRepositorySingleton == nil {
		awsConfigurationRepositorySingleton = &repository.AwsConfigurationRepository{
			FileSystem:  prov.GetFileSystem(),
			Environment: prov.GetEnvironment(),
		}
	}
	return awsConfigurationRepositorySingleton
}
