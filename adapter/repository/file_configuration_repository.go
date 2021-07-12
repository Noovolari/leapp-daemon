package repository

import (
	"encoding/json"
	"fmt"
	"leapp_daemon/domain"
	"leapp_daemon/infrastructure/http/http_error"
	"sync"
)

const configurationFilePath = `.Leapp/Leapp-lock.json`

// The zero value is an unlocked mutex
var configurationFileMutex sync.Mutex

type FileSystem interface {
	DoesFileExist(path string) bool
	GetHomeDir() (string, error)
	ReadFile(path string) ([]byte, error)
	RemoveFile(path string) error
	WriteToFile(path string, data []byte) error
}

type Encryption interface {
	Decrypt(encryptedText string) (string, error)
	Encrypt(plainText string) (string, error)
}

type FileConfigurationRepository struct {
	FileSystem FileSystem
	Encryption Encryption
}

func (repository *FileConfigurationRepository) CreateConfiguration(config domain.Configuration) error {
	configurationFileMutex.Lock()
	defer configurationFileMutex.Unlock()

	homeDir, err := repository.FileSystem.GetHomeDir()
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	configurationFilePath := buildConfigurationFilePath(homeDir, configurationFilePath)

	if repository.FileSystem.DoesFileExist(configurationFilePath) {
		err = repository.FileSystem.RemoveFile(configurationFilePath)
		if err != nil {
			return http_error.NewInternalServerError(err)
		}
	}

	configurationJson, err := json.Marshal(config)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	encryptedConfigurationJson, err := repository.Encryption.Encrypt(string(configurationJson))
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = repository.FileSystem.WriteToFile(configurationFilePath, []byte(encryptedConfigurationJson))
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func (repository *FileConfigurationRepository) GetConfiguration() (domain.Configuration, error) {
	var config domain.Configuration

	homeDir, err := repository.FileSystem.GetHomeDir()
	if err != nil {
		return config, http_error.NewInternalServerError(err)
	}

	configurationFilePath := fmt.Sprintf("%s/%s", homeDir, configurationFilePath)

	if !repository.FileSystem.DoesFileExist(configurationFilePath) {
		return domain.GetDefaultConfiguration(), nil
	}

	encryptedText, err := repository.FileSystem.ReadFile(configurationFilePath)
	if err != nil {
		return config, http_error.NewInternalServerError(err)
	}

	plainText, err := repository.Encryption.Decrypt(string(encryptedText))
	if err != nil {
		return config, http_error.NewInternalServerError(err)
	}

	config = domain.FromJson(plainText)
	return config, nil
}

func (repository *FileConfigurationRepository) UpdateConfiguration(config domain.Configuration) error {
	configurationFileMutex.Lock()
	defer configurationFileMutex.Unlock()

	homeDir, err := repository.FileSystem.GetHomeDir()
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	configurationFilePath := buildConfigurationFilePath(homeDir, configurationFilePath)

	configurationJson, err := json.Marshal(config)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	encryptedConfigurationJson, err := repository.Encryption.Encrypt(string(configurationJson))
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = repository.FileSystem.WriteToFile(configurationFilePath, []byte(encryptedConfigurationJson))
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func buildConfigurationFilePath(homeDir string, configurationFilePath string) string {
	return fmt.Sprintf("%s/%s", homeDir, configurationFilePath)
}
