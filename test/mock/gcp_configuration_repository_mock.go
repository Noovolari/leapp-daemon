package mock

import (
	"fmt"
)

type GcpConfigurationRepositoryMock struct {
	calls                   []string
	ExpGcpConfigFolderExist bool
}

func NewGcpConfigurationRepositoryMock() GcpConfigurationRepositoryMock {
	return GcpConfigurationRepositoryMock{calls: []string{}}
}

func (repo *GcpConfigurationRepositoryMock) GetCalls() []string {
	return repo.calls
}

func (repo *GcpConfigurationRepositoryMock) DoesGcloudConfigFolderExist() (bool, error) {
	repo.calls = append(repo.calls, "DoesGcloudConfigFolderExist()")
	return repo.ExpGcpConfigFolderExist, nil
}

func (repo *GcpConfigurationRepositoryMock) CreateConfiguration(account string, project string) error {
	repo.calls = append(repo.calls, fmt.Sprintf("CreateConfiguration(%v, %v)", account, project))
	return nil
}

func (repo *GcpConfigurationRepositoryMock) RemoveConfiguration() error {
	repo.calls = append(repo.calls, "RemoveConfiguration()")
	return nil
}

func (repo *GcpConfigurationRepositoryMock) ActivateConfiguration() error {
	repo.calls = append(repo.calls, "ActivateConfiguration()")
	return nil
}

func (repo *GcpConfigurationRepositoryMock) DeactivateConfiguration() error {
	repo.calls = append(repo.calls, "DeactivateConfiguration()")
	return nil
}

func (repo *GcpConfigurationRepositoryMock) WriteDefaultCredentials(credentialsJson string) error {
	repo.calls = append(repo.calls, fmt.Sprintf("WriteDefaultCredentials(%v)", credentialsJson))
	return nil
}

func (repo *GcpConfigurationRepositoryMock) RemoveDefaultCredentials() error {
	repo.calls = append(repo.calls, "RemoveDefaultCredentials()")
	return nil
}

func (repo *GcpConfigurationRepositoryMock) WriteCredentialsToDb(accountId string, credentialsJson string) error {
	repo.calls = append(repo.calls, fmt.Sprintf("WriteDefaultCredentials(%v, %v)", accountId, credentialsJson))
	return nil
}

func (repo *GcpConfigurationRepositoryMock) RemoveCredentialsFromDb(accountId string) error {
	repo.calls = append(repo.calls, fmt.Sprintf("RemoveCredentialsFromDb(%v)", accountId))
	return nil
}

func (repo *GcpConfigurationRepositoryMock) RemoveAccessTokensFromDb(accountId string) error {
	repo.calls = append(repo.calls, fmt.Sprintf("RemoveAccessTokensFromDb(%v)", accountId))
	return nil
}
