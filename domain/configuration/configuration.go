package configuration

import (
	"encoding/json"
	"fmt"
	"leapp_daemon/domain/named_profile"
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/http/http_error"
)

type Configuration struct {
	ProxyConfiguration              ProxyConfiguration
	AwsIamUserSessions              []session.AwsIamUserSession
	AwsIamRoleFederatedSessions     []session.AwsIamRoleFederatedSession
	AwsIamRoleChainedSessions       []session.AwsIamRoleChainedSession
	GcpIamUserAccountOauthSessions  []session.GcpIamUserAccountOauthSession
	AlibabaRamUserSessions          []session.AlibabaRamUserSession
	AlibabaRamRoleFederatedSessions []session.AlibabaRamRoleFederatedSession
	AlibabaRamRoleChainedSessions   []session.AlibabaRamRoleChainedSession
	NamedProfiles                   []named_profile.NamedProfile
}

type ProxyConfiguration struct {
	ProxyProtocol string
	ProxyUrl      string
	ProxyPort     uint64
	Username      string
	Password      string
}

func GetDefaultConfiguration() Configuration {
	return Configuration{
		ProxyConfiguration: ProxyConfiguration{
			ProxyProtocol: "https",
			ProxyUrl:      "",
			ProxyPort:     8080,
			Username:      "",
			Password:      "",
		},
		AwsIamRoleFederatedSessions:    make([]session.AwsIamRoleFederatedSession, 0),
		AwsIamUserSessions:             make([]session.AwsIamUserSession, 0),
		GcpIamUserAccountOauthSessions: make([]session.GcpIamUserAccountOauthSession, 0),
		AlibabaRamUserSessions:         make([]session.AlibabaRamUserSession, 0),
	}
}

func FromJson(configurationJson string) Configuration {
	var config Configuration
	_ = json.Unmarshal([]byte(configurationJson), &config)
	return config
}

func (config *Configuration) AddAwsIamUserSession(awsIamUserSession session.AwsIamUserSession) error {
	sessions, err := config.GetAllAwsIamUserSessions()
	if err != nil {
		return err
	}

	for _, sess := range sessions {
		if awsIamUserSession.Id == sess.Id {
			return http_error.NewConflictError(fmt.Errorf("a AwsIamUserSession with id " + awsIamUserSession.Id +
				" is already present"))
		}
	}

	sessions = append(sessions, awsIamUserSession)
	config.AwsIamUserSessions = sessions

	return nil
}

func (config *Configuration) GetAllAwsIamUserSessions() ([]session.AwsIamUserSession, error) {
	return config.AwsIamUserSessions, nil
}

func (config *Configuration) RemoveAwsIamUserSession(awsIamUserSession session.AwsIamUserSession) error {
	sessions, err := config.GetAllAwsIamUserSessions()
	if err != nil {
		return err
	}

	for i, sess := range sessions {
		if awsIamUserSession.Id == sess.Id {
			config.AwsIamUserSessions = append(config.AwsIamUserSessions[:i], config.AwsIamUserSessions[i+1:]...)
			return nil
		}
	}

	return http_error.NewNotFoundError(fmt.Errorf("AwsIamUserSession with id " + awsIamUserSession.Id +
		" not found"))
}

func (config *Configuration) AddNamedProfile(namedProfile named_profile.NamedProfile) error {
	for _, tmpNamedProfile := range config.NamedProfiles {
		if namedProfile.Name == tmpNamedProfile.Name {
			return http_error.NewConflictError(fmt.Errorf("a NamedProfile with name " + namedProfile.Name +
				" is already present"))
		}
	}
	config.NamedProfiles = append(config.NamedProfiles, namedProfile)
	return nil
}

func (config *Configuration) FindNamedProfileByName(name string) (named_profile.NamedProfile, error) {
	var namedProfile named_profile.NamedProfile

	for _, tmpNamedProfile := range config.NamedProfiles {
		if name == tmpNamedProfile.Name {
			return tmpNamedProfile, nil
		}
	}

	return namedProfile, http_error.NewNotFoundError(fmt.Errorf("NamedProfile with name " + name + " not found"))
}

func (config *Configuration) DoesNamedProfileExist(name string) bool {
	for _, tmpNamedProfile := range config.NamedProfiles {
		if name == tmpNamedProfile.Name {
			return true
		}
	}
	return false
}

/*
// The zero value is an unlocked mutex
var configurationFileMutex sync.Mutex

func CreateConfiguration() error {
	configuration := GetDefaultConfiguration()
	err := UpdateConfiguration(configuration, true)
	if err != nil { return err }
	return nil
}

func ReadConfiguration() (*Configuration, error) {
	homeDir, err := file_system.GetHomeDir()
	if err != nil { return nil, err }

	configurationFilePath := fmt.Sprintf("%s/%s", homeDir, configurationFilePath)

	encryptedText, err := ioutil.ReadFile(configurationFilePath)
	if err != nil {
		return nil, http_error.NewInternalServerError(err)
	}

	plainText, err := encryption.Decrypt(string(encryptedText))
	if err != nil { return nil, err }

	return FromJson(plainText), nil
}

func UpdateConfiguration(configuration *Configuration, deleteExistingFile bool) error {
	configurationFileMutex.Lock()
	defer configurationFileMutex.Unlock()

	homeDir, err := file_system.GetHomeDir()
	if err != nil { return err }

	configurationFilePath := fmt.Sprintf("%s/%s", homeDir, configurationFilePath)

	if deleteExistingFile == true {
		if file_system.DoesFileExist(configurationFilePath) {
			err = os.Remove(configurationFilePath)
			if err != nil {
				return err
			}
		}
	}

	configurationJson, err := json.Marshal(configuration)
	if err != nil { return err }

	encryptedConfigurationJson, err := encryption.Encrypt(string(configurationJson))
	if err != nil { return err }

	err = ioutil.WriteToFile(configurationFilePath, []byte(encryptedConfigurationJson), 0644)
	if err != nil { return err }

	return nil
}

func (config *Configuration) Update() error {
  err := UpdateConfiguration(config, false)
  if err != nil {
    return err
  }
  return nil
}

func (config *Configuration) AwsGetIamRoleFederatedSessions() ([]*session.AwsIamRoleFederatedSession, error) {
	return config.AwsIamRoleFederatedSessions, nil
}

func (config *Configuration) SetFederatedAwsSessions(federatedAwsSessions []*session.AwsIamRoleFederatedSession) error {
	config.AwsIamRoleFederatedSessions = federatedAwsSessions
	return nil
}

func (config *Configuration) GetAwsIamRoleChainedSessions() ([]*session.AwsIamRoleChainedSession, error) {
  return config.AwsIamRoleChainedSessions, nil
}

func (config *Configuration) SetAwsIamRoleChainedSessions(trustedAwsSessions []*session.AwsIamRoleChainedSession) error {
  config.AwsIamRoleChainedSessions = trustedAwsSessions
  return nil
}

func (config *Configuration) GetNamedProfiles() ([]*named_profile.NamedProfile, error) {
	return config.NamedProfiles, nil
}

func (config *Configuration) SetNamedProfiles(namedProfiles []*named_profile.NamedProfile) error {
	config.NamedProfiles = namedProfiles
	return nil
}

func (config *Configuration) GetAllSessions() []session.Rotatable {
  sessions := make([]session.Rotatable, 0)

  for i := range config.awsIamUserSessions {
    sess := config.awsIamUserSessions[i]
    sessions = append(sessions, sess)
  }

  for i := range config.AwsIamRoleFederatedSessions {
    sess := config.AwsIamRoleFederatedSessions[i]
    sessions = append(sessions, sess)
  }

  return sessions
}
*/
