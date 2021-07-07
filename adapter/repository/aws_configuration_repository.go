package repository

import (
	"fmt"
	"gopkg.in/ini.v1"
	"leapp_daemon/infrastructure/http/http_error"
	"path/filepath"
	"strings"
)

type AwsConfigFileSystem interface {
	DoesFileExist(path string) bool
	GetHomeDir() (string, error)
	WriteToFile(path string, data []byte) error
	RemoveFile(path string) error
	RenameFile(from string, to string) error
}

type AwsEnvironment interface {
	GetEnvironmentVariable(variable string) string
	IsCommandAvailable(command string) bool
	IsWindows() bool
	GetFormattedTime(format string) string
}

type AwsConfigurationRepository struct {
	FileSystem  AwsConfigFileSystem
	Environment AwsEnvironment
}

type AwsTempCredentials struct {
	ProfileName  string
	AccessKeyId  string
	SecretKey    string
	SessionToken string
	Expiration   string
	Region       string
}

func (repo *AwsConfigurationRepository) getAwsDirPath() (string, error) {
	homeDir, err := repo.FileSystem.GetHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".aws"), nil
}

func (repo *AwsConfigurationRepository) getCredentialsFilePath() (string, error) {
	awsDir, err := repo.getAwsDirPath()
	if err != nil {
		return "", err
	}

	return filepath.Join(awsDir, "credentials"), nil
}

func (repo *AwsConfigurationRepository) getCredentialsBackupFilePath() (string, error) {
	awsDir, err := repo.getAwsDirPath()
	if err != nil {
		return "", err
	}
	timeText := repo.Environment.GetFormattedTime("2006-01-02 15-04-05")
	return filepath.Join(awsDir, fmt.Sprintf("credentials_%v.backup", timeText)), nil
}

func (repo *AwsConfigurationRepository) createProfileSection(credentialsIniFile *ini.File, credentials AwsTempCredentials) error {

	credentialsIniFile.DeleteSection(credentials.ProfileName)
	section, err := credentialsIniFile.NewSection(credentials.ProfileName)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}
	section.Comment = "Credentials managed by Leapp, do not manually edit this file"

	_, err = section.NewKey("aws_access_key_id", credentials.AccessKeyId)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	_, err = section.NewKey("aws_secret_access_key", credentials.SecretKey)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	_, err = section.NewKey("aws_session_token", credentials.SessionToken)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	_, err = section.NewKey("expiration", credentials.Expiration)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	if credentials.Region != "" {
		_, err = section.NewKey("region", credentials.Region)
		if err != nil {
			return http_error.NewInternalServerError(err)
		}
	}
	return nil
}

func (repo *AwsConfigurationRepository) backupIfNeeded(iniFile *ini.File, credentialsFilePath string) error {
	iniSections := iniFile.Sections()
	backupIniFile := false
	for _, iniSection := range iniSections {
		if iniSection.Name() != ini.DefaultSection && !strings.Contains(iniSection.Comment, "Leapp") {
			backupIniFile = true
			break
		}
	}

	if backupIniFile {
		err := repo.backupCredentialsFile(credentialsFilePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *AwsConfigurationRepository) backupCredentialsFile(credentialsFilePath string) error {
	backupFilePath, err := repo.getCredentialsBackupFilePath()
	err = repo.FileSystem.RenameFile(credentialsFilePath, backupFilePath)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}
	return nil
}

func (repo *AwsConfigurationRepository) WriteCredentials(credentials []AwsTempCredentials) error {

	credentialsFilePath, err := repo.getCredentialsFilePath()
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	var iniFile *ini.File
	if repo.FileSystem.DoesFileExist(credentialsFilePath) {
		iniFile, err = ini.Load(credentialsFilePath)
		if err != nil {
			return err
		}

		err = repo.backupIfNeeded(iniFile, credentialsFilePath)
		if err != nil {
			return err
		}
	}

	iniFile = ini.Empty()
	for _, credential := range credentials {
		err = repo.createProfileSection(iniFile, credential)
		if err != nil {
			return err
		}
	}

	err = iniFile.SaveTo(credentialsFilePath)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}
	return nil
}
