package mock

import (
	"errors"
	"leapp_daemon/domain/configuration"
	"leapp_daemon/infrastructure/http/http_error"
)

type FileConfigurationRepositoryMock struct {
	calls                         []string
	ExpErrorOnGetConfiguration    bool
	ExpErrorOnUpdateConfiguration bool
	ExpGetConfiguration           configuration.Configuration
}

func NewFileConfigurationRepositoryMock() FileConfigurationRepositoryMock {
	return FileConfigurationRepositoryMock{calls: []string{}}
}

func (repo *FileConfigurationRepositoryMock) GetCalls() []string {
	return repo.calls
}

func (repo *FileConfigurationRepositoryMock) CreateConfiguration(configuration configuration.Configuration) error {
	repo.calls = append(repo.calls, "WriteDefaultCredentials()")
	return nil
}

func (repo *FileConfigurationRepositoryMock) GetConfiguration() (configuration.Configuration, error) {
	repo.calls = append(repo.calls, "GetConfiguration()")
	if repo.ExpErrorOnGetConfiguration {
		return configuration.Configuration{}, http_error.NewInternalServerError(errors.New("error getting configuration"))
	}

	return repo.ExpGetConfiguration, nil
}

func (repo *FileConfigurationRepositoryMock) UpdateConfiguration(configuration.Configuration) error {
	repo.calls = append(repo.calls, "UpdateConfiguration()")

	if repo.ExpErrorOnUpdateConfiguration {
		return http_error.NewInternalServerError(errors.New("error updating configuration"))
	}

	return nil
}
