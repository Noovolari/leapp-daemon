package session_token

import (
	"encoding/json"
	//"fmt"
	"github.com/aws/aws-sdk-go/service/sts"
	"leapp_daemon/domain/constant"
	"leapp_daemon/infrastructure/aws/sts_client"
	"leapp_daemon/infrastructure/http/http_error"
	"leapp_daemon/infrastructure/keychain"
	"sync"
	"time"
)

// The zero value is an unlocked mutex
var keychainMutex sync.Mutex

// The zero value is an unlocked mutex
var iniFileMutex sync.Mutex

func DoExist(accountName string) (bool, error) {
	doesSessionTokenExpirationExist, err := (&keychain.Keychain{}).DoesSecretExist(accountName + "-aws-iam-user-session-session-token-expiration")
	if err != nil {
		return false, err
	}

	doesSessionTokenExist, err := (&keychain.Keychain{}).DoesSecretExist(accountName + "-aws-iam-user-session-session-token")
	if err != nil {
		return false, err
	}

	return doesSessionTokenExpirationExist && doesSessionTokenExist, nil
}

func IsExpired(accountName string) (bool, error) {
	secret, err := (&keychain.Keychain{}).GetSecret(accountName + "-aws-iam-user-session-session-token-expiration")
	if err != nil {
		return false, err
	}

	t, err := time.Parse(time.RFC3339, secret)
	if err != nil {
		return false, err
	}

	return time.Now().After(t), nil
}

func Generate(accountName string, region string, mfaDevice string, mfaToken *string) (*sts.Credentials, error) {
	accessKeyId, secretAccessKey, err := Get(accountName)
	if err != nil {
		return nil, err
	}

	stsClient, err := sts_client.GetStaticCredentialsClient(accessKeyId, secretAccessKey, &region)
	if err != nil {
		return nil, err
	}

	durationSeconds := constant.SessionTokenDurationInSeconds
	var getSessionTokenInput sts.GetSessionTokenInput

	if mfaToken == nil {
		getSessionTokenInput = sts.GetSessionTokenInput{
			DurationSeconds: &durationSeconds,
			SerialNumber:    nil,
			TokenCode:       nil,
		}
	} else {
		getSessionTokenInput = sts.GetSessionTokenInput{
			DurationSeconds: &durationSeconds,
			SerialNumber:    &mfaDevice,
			TokenCode:       mfaToken,
		}
	}

	getSessionTokenOutput, err := stsClient.GetSessionToken(&getSessionTokenInput)
	if err != nil {
		return nil, err
	}

	return getSessionTokenOutput.Credentials, nil
}

func Get(accountName string) (string, string, error) {
	sessionTokenSecretName := accountName + "-aws-iam-user-session-session-token"

	sessionToken, err := (&keychain.Keychain{}).GetSecret(sessionTokenSecretName)
	if err != nil {
		return "", "", http_error.NewUnprocessableEntityError(err)
	}

	sessionTokenExpirationSecretName := accountName + "-aws-iam-user-session-session-token-expiration"

	sessionTokenExpiration, err := (&keychain.Keychain{}).GetSecret(sessionTokenExpirationSecretName)
	if err != nil {
		return "", "", http_error.NewUnprocessableEntityError(err)
	}

	return sessionToken, sessionTokenExpiration, nil
}

func SaveInKeychain(accountName string, credentials *sts.Credentials) error {
	keychainMutex.Lock()
	defer keychainMutex.Unlock()

	credentialsJson, err := json.Marshal(credentials)
	if err != nil {
		return err
	}

	err = (&keychain.Keychain{}).SetSecret(string(credentialsJson),
		accountName+"-aws-iam-user-session-session-token")
	if err != nil {
		return err
	}

	err = (&keychain.Keychain{}).SetSecret(credentials.Expiration.Format(time.RFC3339),
		accountName+"-aws-iam-user-session-session-token-expiration")
	if err != nil {
		return err
	}

	return nil
}

func SaveInIniFile(accessKeyId string, secretAccessKey string, sessionToken string, region string, profileName string) error {
	/*
		iniFileMutex.Lock()
		defer iniFileMutex.Unlock()

		homeDir, err := file_system.GetHomeDir()
		if err != nil {
			return err
		}

		credentialsFilePath := homeDir + "/" + constant.CredentialsFilePath

		if file_system.DoesFileExist(credentialsFilePath) {
			credentialsFile, err := ini.Load(credentialsFilePath)
			if err != nil {
				return err
			}

			section, err := credentialsFile.GetSection(profileName)
			if err != nil && err.Error() != fmt.Sprintf("section %q does not exist", profileName){
				return err
			}

			if section == nil {
				_, err = credentials_ini_file.CreateNamedProfileSection(credentialsFile, profileName, accessKeyId, secretAccessKey,
					sessionToken, region)
				if err != nil {
					return err
				}

				err = credentials_ini_file.AppendToFile(credentialsFile, credentialsFilePath)
				if err != nil {
					return err
				}
			} else {
				credentialsFile.DeleteSection(profileName)

				_, err = credentials_ini_file.CreateNamedProfileSection(credentialsFile, profileName, accessKeyId, secretAccessKey,
					sessionToken, region)
				if err != nil {
					return err
				}

				err = credentials_ini_file.OverwriteFile(credentialsFile, credentialsFilePath)
				if err != nil {
					return err
				}
			}
		} else {
			credentialsFile := ini.Empty()

			_, err = credentials_ini_file.CreateNamedProfileSection(credentialsFile, profileName, accessKeyId, secretAccessKey,
				sessionToken, region)
			if err != nil {
				return err
			}

			err = credentials_ini_file.OverwriteFile(credentialsFile, credentialsFilePath)
			if err != nil {
				return err
			}
		}
	*/

	return nil
}

func RemoveFromIniFile(profileName string) error {
	/*
		iniFileMutex.Lock()
		defer iniFileMutex.Unlock()

		homeDir, err := file_system.GetHomeDir()
		if err != nil {
			return err
		}

		credentialsFilePath := homeDir + "/" + constant.CredentialsFilePath

		if file_system.DoesFileExist(credentialsFilePath) {
			credentialsFile, err := ini.Load(credentialsFilePath)
			if err != nil {
				return err
			}

			credentialsFile.DeleteSection(profileName)

			err = credentials_ini_file.OverwriteFile(credentialsFile, credentialsFilePath)
			if err != nil {
				return err
			}
		}
	*/

	return nil
}
