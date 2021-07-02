package repository

import (
	"fmt"
	"leapp_daemon/infrastructure/http/http_error"
	"path/filepath"
)

type GcpConfigFileSystem interface {
	DoesFileExist(path string) bool
	GetHomeDir() (string, error)
	WriteToFile(path string, data []byte) error
	RemoveFile(path string) error
}

type GcpEnvironment interface {
	GetEnvironmentVariable(variable string) string
	IsCommandAvailable(command string) bool
	IsWindows() bool
}

type GcpCredentialsTable interface {
	WriteCredentials(sqlFilePath string, accountId string, value string) error
	RemoveCredentials(sqlFilePath string, accountId string) error
}

type GcpAccessTokensTable interface {
	RemoveAccessToken(sqlFilePath string, accountId string) error
}

type GcpConfigurationRepository struct {
	FileSystem        GcpConfigFileSystem
	Environment       GcpEnvironment
	CredentialsTable  GcpCredentialsTable
	AccessTokensTable GcpAccessTokensTable
}

func (repo *GcpConfigurationRepository) gcloudConfigDir() (string, error) {
	gcloudDirectory := "gcloud"
	if repo.Environment.IsWindows() {
		return filepath.Join(repo.Environment.GetEnvironmentVariable("APPDATA"), gcloudDirectory), nil
	}

	dir, err := repo.FileSystem.GetHomeDir()
	if err != nil {
		return "", http_error.NewInternalServerError(err)
	}

	return filepath.Join(dir, ".config", gcloudDirectory), nil
}

func (repo *GcpConfigurationRepository) getGcloudConfigFilePath() (string, error) {
	configDir, err := repo.gcloudConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "configurations", "config_leapp"), nil
}

func (repo *GcpConfigurationRepository) getGcloudActiveConfigFilePath() (string, error) {
	configDir, err := repo.gcloudConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "active_config"), nil
}

func (repo *GcpConfigurationRepository) getGcloudDefaultCredentialsFilePath() (string, error) {
	configDir, err := repo.gcloudConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "application_default_credentials.json"), nil
}

func (repo *GcpConfigurationRepository) getGcloudCredentialsDbFilePath() (string, error) {
	configDir, err := repo.gcloudConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "credentials.db"), nil
}

func (repo *GcpConfigurationRepository) getGcloudAccessTokensDbFilePath() (string, error) {
	configDir, err := repo.gcloudConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "access_tokens.db"), nil
}

func (repo *GcpConfigurationRepository) isGcloudCliAvailable() bool {
	return repo.Environment.IsCommandAvailable("gcloud")
}

func (repo *GcpConfigurationRepository) DoesGcloudConfigFolderExist() (bool, error) {
	configDir, err := repo.gcloudConfigDir()
	if err != nil {
		return false, err
	}

	return repo.FileSystem.DoesFileExist(configDir), nil
}

func (repo *GcpConfigurationRepository) CreateConfiguration(account string, project string) error {
	configFilePath, err := repo.getGcloudConfigFilePath()
	if err != nil {
		return err
	}

	configFileContent := fmt.Sprintf("[core]\naccount = %v\nproject = %v\n", account, project)
	err = repo.FileSystem.WriteToFile(configFilePath, []byte(configFileContent))
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func (repo *GcpConfigurationRepository) RemoveConfiguration() error {
	configFilePath, err := repo.getGcloudConfigFilePath()
	if err != nil {
		return err
	}

	err = repo.FileSystem.RemoveFile(configFilePath)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}
	return nil
}

func (repo *GcpConfigurationRepository) ActivateConfiguration() error {
	activeConfigFilePath, err := repo.getGcloudActiveConfigFilePath()
	if err != nil {
		return err
	}

	err = repo.FileSystem.WriteToFile(activeConfigFilePath, []byte("leapp"))
	if err != nil {
		return http_error.NewInternalServerError(err)
	}
	return nil
}

func (repo *GcpConfigurationRepository) DeactivateConfiguration() error {
	activeConfigFilePath, err := repo.getGcloudActiveConfigFilePath()
	if err != nil {
		return err
	}

	err = repo.FileSystem.RemoveFile(activeConfigFilePath)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func (repo *GcpConfigurationRepository) WriteDefaultCredentials(credentialsJson string) error {
	defaultCredentialsFilePath, err := repo.getGcloudDefaultCredentialsFilePath()
	if err != nil {
		return err
	}

	err = repo.FileSystem.WriteToFile(defaultCredentialsFilePath, []byte(credentialsJson))
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func (repo *GcpConfigurationRepository) RemoveDefaultCredentials() error {
	defaultCredentialFilePath, err := repo.getGcloudDefaultCredentialsFilePath()
	if err != nil {
		return err
	}

	err = repo.FileSystem.RemoveFile(defaultCredentialFilePath)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func (repo *GcpConfigurationRepository) WriteCredentialsToDb(accountId string, credentialsJson string) error {
	credentialsDbFilePath, err := repo.getGcloudCredentialsDbFilePath()
	if err != nil {
		return err
	}

	err = repo.CredentialsTable.WriteCredentials(credentialsDbFilePath, accountId, credentialsJson)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func (repo *GcpConfigurationRepository) RemoveCredentialsFromDb(accountId string) error {
	credentialsDbFilePath, err := repo.getGcloudCredentialsDbFilePath()
	if err != nil {
		return err
	}

	if !repo.FileSystem.DoesFileExist(credentialsDbFilePath) {
		return nil
	}

	err = repo.CredentialsTable.RemoveCredentials(credentialsDbFilePath, accountId)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func (repo *GcpConfigurationRepository) RemoveAccessTokensFromDb(accountId string) error {
	accessTokensDbFilePath, err := repo.getGcloudAccessTokensDbFilePath()
	if err != nil {
		return err
	}

	if !repo.FileSystem.DoesFileExist(accessTokensDbFilePath) {
		return nil
	}

	err = repo.AccessTokensTable.RemoveAccessToken(accessTokensDbFilePath, accountId)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}
