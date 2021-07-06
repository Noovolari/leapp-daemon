package mock

import (
	"errors"
	"fmt"
	"leapp_daemon/infrastructure/http/http_error"
	"leapp_daemon/interface/aws"
)

type AwsConfigurationRepositoryMock struct {
	calls                      []string
	ExpAwsConfigFolderExist    bool
	ExpErrorOnWriteCredentials bool
}

func NewAwsConfigurationRepositoryMock() AwsConfigurationRepositoryMock {
	return AwsConfigurationRepositoryMock{calls: []string{}}
}

func (repo *AwsConfigurationRepositoryMock) GetCalls() []string {
	return repo.calls
}

func (repo *AwsConfigurationRepositoryMock) WriteCredentials(credentials []aws.AwsTempCredentials) error {
	repo.calls = append(repo.calls, fmt.Sprintf("WriteCredentials(%v)", credentials))
	if repo.ExpErrorOnWriteCredentials {
		return http_error.NewInternalServerError(errors.New("WriteCredentials failed"))
	}
	return nil
}
