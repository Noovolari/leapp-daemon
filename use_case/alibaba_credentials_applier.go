package use_case

import (
	"encoding/json"
	"io/ioutil"
	"leapp_daemon/domain/constant"
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/http/http_error"
	"os"
)

type CredentialsFile struct {
	Current   string                `json:"current"`
	Profiles  []NamedProfileSection `json:"profiles"`
	Meta_path string                `json:"meta_path"`
}

type NamedProfileSection struct {
	Name              string `json:"name"`
	Mode              string `json:"mode"`
	Access_key_id     string `json:"access_key_id"`
	Access_key_secret string `json:"access_key_secret"`
	Sts_token         string `json:"sts_token"`
	Ram_role_name     string `json:"ram_role_name"`
	Ram_role_arn      string `json:"ram_role_arn"`
	Ram_session_name  string `json:"ram_session_name"`
	Private_key       string `json:"private_key"`
	Key_pair_name     string `json:"key_pair_name"`
	Expired_seconds   int    `json:"expired_seconds"`
	Verified          string `json:"verified"`
	Region_id         string `json:"region_id"`
	Output_format     string `json:"output_format"`
	Language          string `json:"language"`
	Site              string `json:"site"`
	Retry_timeout     int    `json:"retry_timeout"`
	Connect_timeout   int    `json:"connect_timeout"`
	Retry_count       int    `json:"retry_count"`
	Process_command   string `json:"process_command"`
}

type AlibabaCredentialsApplier struct {
	FileSystem          FileSystem
	Keychain            Keychain
	NamedProfilesFacade NamedProfilesFacade
}

func (alibabaCredentialsApplier *AlibabaCredentialsApplier) UpdateAlibabaRamUserSessions(oldAlibabaRamUserSessions []session.AlibabaRamUserSession, newAlibabaRamUserSessions []session.AlibabaRamUserSession) error {
	for i, oldSess := range oldAlibabaRamUserSessions {
		if i < len(newAlibabaRamUserSessions) {
			newSess := newAlibabaRamUserSessions[i]

			if oldSess.Status == session.NotActive && newSess.Status == session.Pending {

				homeDir, err := alibabaCredentialsApplier.FileSystem.GetHomeDir()
				if err != nil {
					return err
				}

				credentialsFilePath := homeDir + "/" + constant.AlibabaCredentialsFilePath
				profile, err := alibabaCredentialsApplier.NamedProfilesFacade.GetNamedProfileById(newSess.Account.NamedProfileId)
				if err != nil {
					return err
				}
				profileName := profile.Name
				region := newSess.Account.Region

				accessKeyId, secretAccessKey, err := alibabaCredentialsApplier.getPlainCreds(newSess.Id)
				if err != nil {
					return err
				}

				namedProfileSection := NamedProfileSection{Name: profileName, Mode: "AK", Access_key_id: accessKeyId, Access_key_secret: secretAccessKey, Region_id: region, Output_format: "json", Language: "en"}
				profiles := []NamedProfileSection{namedProfileSection}
				config := CredentialsFile{Current: namedProfileSection.Name, Profiles: profiles}
				out, _ := json.MarshalIndent(config, "", "  ")
				alibabaCredentialsApplier.overwriteFile(out, credentialsFilePath)
			}

			if oldSess.Status == session.Active && newSess.Status == session.NotActive {
				err := alibabaCredentialsApplier.deleteConfig()
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (alibabaCredentialsApplier *AlibabaCredentialsApplier) UpdateAlibabaRamRoleFederatedSessions(oldAlibabaRamRoleFederatedSessions []session.AlibabaRamRoleFederatedSession, newAlibabaRamRoleFederatedSessions []session.AlibabaRamRoleFederatedSession) error {
	for i, oldSess := range oldAlibabaRamRoleFederatedSessions {
		if i < len(newAlibabaRamRoleFederatedSessions) {
			newSess := newAlibabaRamRoleFederatedSessions[i]

			if oldSess.Status == session.NotActive && newSess.Status == session.Pending {

				homeDir, err := alibabaCredentialsApplier.FileSystem.GetHomeDir()
				if err != nil {
					return err
				}

				credentialsFilePath := homeDir + "/" + constant.AlibabaCredentialsFilePath
				profile, err := alibabaCredentialsApplier.NamedProfilesFacade.GetNamedProfileById(newSess.Account.NamedProfileId)
				if err != nil {
					return err
				}
				profileName := profile.Name
				region := newSess.Account.Region

				accessKeyId, secretAccessKey, stsToken, err := alibabaCredentialsApplier.getFederatedCreds(newSess.Id)
				if err != nil {
					return err
				}

				namedProfileSection := NamedProfileSection{Name: profileName, Mode: "StsToken", Access_key_id: accessKeyId, Access_key_secret: secretAccessKey, Sts_token: stsToken, Region_id: region, Output_format: "json", Language: "en"}
				profiles := []NamedProfileSection{namedProfileSection}
				config := CredentialsFile{Current: namedProfileSection.Name, Profiles: profiles}
				out, _ := json.MarshalIndent(config, "", "  ")
				alibabaCredentialsApplier.overwriteFile(out, credentialsFilePath)
			}

			if oldSess.Status == session.Active && newSess.Status == session.NotActive {
				err := alibabaCredentialsApplier.deleteConfig()
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (alibabaCredentialsApplier *AlibabaCredentialsApplier) UpdateAlibabaRamRoleChainedSessions(oldAlibabaRamRoleChainedSessions []session.AlibabaRamRoleChainedSession, newAlibabaRamRoleChainedSessions []session.AlibabaRamRoleChainedSession) error {
	for i, oldSess := range oldAlibabaRamRoleChainedSessions {
		if i < len(newAlibabaRamRoleChainedSessions) {
			newSess := newAlibabaRamRoleChainedSessions[i]

			if oldSess.Status == session.NotActive && newSess.Status == session.Pending {

				homeDir, err := alibabaCredentialsApplier.FileSystem.GetHomeDir()
				if err != nil {
					return err
				}

				credentialsFilePath := homeDir + "/" + constant.AlibabaCredentialsFilePath
				// if you got errors check here
				profileName := newSess.Account.Name
				region := newSess.Account.Region

				accessKeyId, secretAccessKey, stsToken, err := alibabaCredentialsApplier.getTrustedCreds(newSess.Id)
				if err != nil {
					return err
				}

				namedProfileSection := NamedProfileSection{Name: profileName, Mode: "StsToken", Access_key_id: accessKeyId, Access_key_secret: secretAccessKey, Sts_token: stsToken, Region_id: region, Output_format: "json", Language: "en"}
				profiles := []NamedProfileSection{namedProfileSection}
				config := CredentialsFile{Current: namedProfileSection.Name, Profiles: profiles}
				out, _ := json.MarshalIndent(config, "", "  ")
				alibabaCredentialsApplier.overwriteFile(out, credentialsFilePath)
			}

			if oldSess.Status == session.Active && newSess.Status == session.NotActive {
				err := alibabaCredentialsApplier.deleteConfig()
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (alibabaCredentialsApplier *AlibabaCredentialsApplier) getPlainCreds(sessionId string) (accessKeyId string, secretAccessKey string, err error) {
	accessKeyId = ""
	secretAccessKey = ""

	accessKeyIdSecretName := sessionId + constant.PlainAlibabaKeyIdSuffix
	accessKeyId, err = alibabaCredentialsApplier.Keychain.GetSecret(accessKeyIdSecretName)
	if err != nil {
		return accessKeyId, secretAccessKey, http_error.NewUnprocessableEntityError(err)
	}

	secretAccessKeySecretName := sessionId + constant.PlainAlibabaSecretAccessKeySuffix
	secretAccessKey, err = alibabaCredentialsApplier.Keychain.GetSecret(secretAccessKeySecretName)
	if err != nil {
		return accessKeyId, secretAccessKey, http_error.NewUnprocessableEntityError(err)
	}

	return accessKeyId, secretAccessKey, nil
}

func (alibabaCredentialsApplier *AlibabaCredentialsApplier) getFederatedCreds(sessionId string) (accessKeyId string, secretAccessKey string, stsToken string, err error) {
	accessKeyId = ""
	secretAccessKey = ""

	accessKeyIdSecretName := sessionId + constant.FederatedAlibabaKeyIdSuffix
	accessKeyId, err = alibabaCredentialsApplier.Keychain.GetSecret(accessKeyIdSecretName)
	if err != nil {
		return accessKeyId, secretAccessKey, stsToken, http_error.NewUnprocessableEntityError(err)
	}

	secretAccessKeySecretName := sessionId + constant.FederatedAlibabaSecretAccessKeySuffix
	secretAccessKey, err = alibabaCredentialsApplier.Keychain.GetSecret(secretAccessKeySecretName)
	if err != nil {
		return accessKeyId, secretAccessKey, stsToken, http_error.NewUnprocessableEntityError(err)
	}

	stsTokenName := sessionId + constant.TrustedAlibabaStsTokenSuffix
	stsToken, err = alibabaCredentialsApplier.Keychain.GetSecret(stsTokenName)
	if err != nil {
		return accessKeyId, secretAccessKey, stsToken, http_error.NewUnprocessableEntityError(err)
	}

	return accessKeyId, secretAccessKey, stsToken, nil
}

func (alibabaCredentialsApplier *AlibabaCredentialsApplier) getTrustedCreds(sessionId string) (accessKeyId string, secretAccessKey string, stsToken string, err error) {
	accessKeyId = ""
	secretAccessKey = ""

	accessKeyIdSecretName := sessionId + constant.TrustedAlibabaKeyIdSuffix
	accessKeyId, err = alibabaCredentialsApplier.Keychain.GetSecret(accessKeyIdSecretName)
	if err != nil {
		return accessKeyId, secretAccessKey, stsToken, http_error.NewUnprocessableEntityError(err)
	}

	secretAccessKeySecretName := sessionId + constant.TrustedAlibabaSecretAccessKeySuffix
	secretAccessKey, err = alibabaCredentialsApplier.Keychain.GetSecret(secretAccessKeySecretName)
	if err != nil {
		return accessKeyId, secretAccessKey, stsToken, http_error.NewUnprocessableEntityError(err)
	}

	stsTokenName := sessionId + constant.TrustedAlibabaStsTokenSuffix
	stsToken, err = alibabaCredentialsApplier.Keychain.GetSecret(stsTokenName)
	if err != nil {
		return accessKeyId, secretAccessKey, stsToken, http_error.NewUnprocessableEntityError(err)
	}

	return accessKeyId, secretAccessKey, stsToken, nil
}

func (alibabaCredentialsApplier *AlibabaCredentialsApplier) overwriteFile(content []byte, path string) error {

	err := ioutil.WriteFile(path, content, 0644)
	if err != nil {
		return http_error.NewUnprocessableEntityError(err)
	}

	return nil
}

func (alibabaCredentialsApplier *AlibabaCredentialsApplier) deleteConfig() error {

	homeDir, err := alibabaCredentialsApplier.FileSystem.GetHomeDir()
	if err != nil {
		return err
	}

	credentialsFilePath := homeDir + "/" + constant.AlibabaCredentialsFilePath
	err = os.Remove(credentialsFilePath)
	if err != nil {
		return err
	}

	return nil
}
