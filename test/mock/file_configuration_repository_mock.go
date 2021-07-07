package mock

import (
	"errors"
	"leapp_daemon/domain"
	"leapp_daemon/infrastructure/http/http_error"
)

type FileConfigurationRepositoryMock struct {
	calls                         []string
	ExpErrorOnGetConfiguration    bool
	ExpErrorOnUpdateConfiguration bool
	ExpGetConfiguration           domain.Configuration
}

func NewFileConfigurationRepositoryMock() FileConfigurationRepositoryMock {
	return FileConfigurationRepositoryMock{calls: []string{}}
}

func (repo *FileConfigurationRepositoryMock) GetCalls() []string {
	return repo.calls
}

func (repo *FileConfigurationRepositoryMock) CreateConfiguration(configuration domain.Configuration) error {
	repo.calls = append(repo.calls, "CreateConfiguration()")
	return nil
}

func (repo *FileConfigurationRepositoryMock) GetConfiguration() (domain.Configuration, error) {
	repo.calls = append(repo.calls, "GetConfiguration()")
	if repo.ExpErrorOnGetConfiguration {
		return domain.Configuration{}, http_error.NewInternalServerError(errors.New("error getting configuration"))
	}

	return repo.ExpGetConfiguration, nil
}

func (repo *FileConfigurationRepositoryMock) UpdateConfiguration(domain.Configuration) error {
	repo.calls = append(repo.calls, "UpdateConfiguration()")

	if repo.ExpErrorOnUpdateConfiguration {
		return http_error.NewInternalServerError(errors.New("error updating configuration"))
	}

	return nil
}
