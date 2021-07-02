package use_case

import (
	"fmt"
	"gopkg.in/ini.v1"
	"leapp_daemon/domain/constant"
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/http/http_error"
	"leapp_daemon/infrastructure/logging"
	"os"
)

type AwsCredentialsApplier struct {
	FileSystem          FileSystem
	Keychain            Keychain
	NamedProfilesFacade NamedProfilesFacade
}

func (applier *AwsCredentialsApplier) UpdateAwsIamUserSessions(oldSessions []session.AwsIamUserSession, newSessions []session.AwsIamUserSession) error {
	for i, oldSess := range oldSessions {
		if i < len(newSessions) {
			newSess := newSessions[i]

			if oldSess.Status == session.NotActive && newSess.Status == session.Pending {

				homeDir, err := applier.FileSystem.GetHomeDir()
				if err != nil {
					return err
				}

				credentialsFilePath := homeDir + "/" + constant.AwsCredentialsFilePath
				namedProfile, err := applier.NamedProfilesFacade.GetNamedProfileById(newSess.Account.NamedProfileId)
				if err != nil {
					return err
				}

				profileName := namedProfile.Name
				region := newSess.Account.Region

				accessKeyId, secretAccessKey, err := applier.getAccessKeys(newSess.Id)
				if err != nil {
					return err
				}

				sessionToken, err := applier.getSessionToken(newSess.Id)
				if err != nil {
					return err
				}

				if applier.FileSystem.DoesFileExist(credentialsFilePath) {
					credentialsFile, err := ini.Load(credentialsFilePath)
					if err != nil {
						return err
					}

					section, err := credentialsFile.GetSection(profileName)
					if err != nil && err.Error() != fmt.Sprintf("section %q does not exist", profileName) {
						return err
					}

					if section == nil {
						_, err = applier.createNamedProfileSection(credentialsFile, profileName, accessKeyId,
							secretAccessKey, sessionToken, region)
						if err != nil {
							return err
						}

						err = applier.appendToFile(credentialsFile, credentialsFilePath)
						if err != nil {
							return err
						}
					} else {
						credentialsFile.DeleteSection(profileName)

						_, err = applier.createNamedProfileSection(credentialsFile, profileName, accessKeyId, secretAccessKey,
							sessionToken, region)
						if err != nil {
							return err
						}

						err = applier.overwriteFile(credentialsFile, credentialsFilePath)
						if err != nil {
							return err
						}
					}
				} else {
					credentialsFile := ini.Empty()

					_, err = applier.createNamedProfileSection(credentialsFile, profileName, accessKeyId, secretAccessKey,
						sessionToken, region)
					if err != nil {
						return err
					}

					err = applier.overwriteFile(credentialsFile, credentialsFilePath)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (applier *AwsCredentialsApplier) getAccessKeys(sessionId string) (accessKeyId string, secretAccessKey string, error error) {
	accessKeyId = ""
	secretAccessKey = ""

	accessKeyIdSecretName := sessionId + "-aws-iam-user-session-access-key-id"
	accessKeyId, err := applier.Keychain.GetSecret(accessKeyIdSecretName)
	if err != nil {
		return accessKeyId, secretAccessKey, http_error.NewUnprocessableEntityError(err)
	}

	secretAccessKeySecretName := sessionId + "-aws-iam-user-session-secret-access-key"
	secretAccessKey, err = applier.Keychain.GetSecret(secretAccessKeySecretName)
	if err != nil {
		return accessKeyId, secretAccessKey, http_error.NewUnprocessableEntityError(err)
	}

	return accessKeyId, secretAccessKey, nil
}

func (applier *AwsCredentialsApplier) getSessionToken(sessionId string) (sessionToken string, error error) {
	sessionToken = ""

	sessionTokenSecretName := sessionId + "-aws-iam-user-session-session-token"
	sessionToken, err := applier.Keychain.GetSecret(sessionTokenSecretName)
	if err != nil {
		return sessionToken, http_error.NewUnprocessableEntityError(err)
	}

	return sessionToken, nil
}

func (applier *AwsCredentialsApplier) createNamedProfileSection(credentialsFile *ini.File, profileName string, accessKeyId string,
	secretAccessKey string, sessionToken string, region string) (*ini.Section, error) {

	section, err := credentialsFile.NewSection(profileName)
	if err != nil {
		return nil, http_error.NewInternalServerError(err)
	}

	_, err = section.NewKey("aws_access_key_id", accessKeyId)
	if err != nil {
		return nil, http_error.NewInternalServerError(err)
	}

	_, err = section.NewKey("aws_secret_access_key", secretAccessKey)
	if err != nil {
		return nil, http_error.NewInternalServerError(err)
	}

	_, err = section.NewKey("aws_session_token", sessionToken)
	if err != nil {
		return nil, http_error.NewInternalServerError(err)
	}

	if region != "" {
		_, err = section.NewKey("region", region)
		if err != nil {
			return nil, http_error.NewInternalServerError(err)
		}
	}

	return section, nil
}

func (applier *AwsCredentialsApplier) appendToFile(file *ini.File, path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		return http_error.NewNotFoundError(err)
	}

	_, err = file.WriteTo(f)
	if err != nil {
		return http_error.NewUnprocessableEntityError(err)
	}

	return nil
}

func (applier *AwsCredentialsApplier) overwriteFile(file *ini.File, path string) error {
	logging.Entry().Error("flag 3")

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return http_error.NewNotFoundError(err)
	}

	_, err = file.WriteTo(f)
	if err != nil {
		return http_error.NewUnprocessableEntityError(err)
	}

	return nil
}
