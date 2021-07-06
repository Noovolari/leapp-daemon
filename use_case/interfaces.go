package use_case

import (
	"github.com/aws/aws-sdk-go/service/sts"
	"golang.org/x/oauth2"
	"leapp_daemon/domain"
	"leapp_daemon/domain/domain_aws/aws_iam_user"
	"leapp_daemon/domain/domain_aws/named_profile"
	"leapp_daemon/domain/domain_gcp/gcp_iam_user_account_oauth"
	"leapp_daemon/domain/domain_gcp/named_configuration"
	"leapp_daemon/interface/aws"
)

type FileSystem interface {
	DoesFileExist(path string) bool
	GetHomeDir() (string, error)
}

type Environment interface {
	GenerateUuid() string
	GetTime() string
}

type Keychain interface {
	DoesSecretExist(label string) (bool, error)
	GetSecret(label string) (string, error)
	DeleteSecret(label string) error
	SetSecret(secret string, label string) error
}

type StsApi interface {
	GenerateNewSessionToken(accessKeyId string, secretKey string, region string, mfaDevice string, mfaToken *string) (*sts.Credentials, error)
}

type GcpApi interface {
	GetOauthUrl() (string, error)
	GetOauthToken(authCode string) (*oauth2.Token, error)
	GetCredentials(oauthToken *oauth2.Token) string
}

type ConfigurationRepository interface {
	CreateConfiguration(domain.Configuration) error
	GetConfiguration() (domain.Configuration, error)
	UpdateConfiguration(domain.Configuration) error
}

type AwsConfigurationRepository interface {
	WriteCredentials(credentials []aws.AwsTempCredentials) error
}

type GcpConfigurationRepository interface {
	DoesGcloudConfigFolderExist() (bool, error)
	CreateConfiguration(account string, project string) error
	RemoveConfiguration() error
	ActivateConfiguration() error
	DeactivateConfiguration() error
	WriteDefaultCredentials(credentialsJson string) error
	RemoveDefaultCredentials() error
	WriteCredentialsToDb(accountId string, credentialsJson string) error
	RemoveCredentialsFromDb(accountId string) error
	RemoveAccessTokensFromDb(accountId string) error
}

type NamedProfilesFacade interface {
	GetNamedProfiles() []named_profile.NamedProfile
	GetNamedProfileById(id string) (named_profile.NamedProfile, error)
	GetNamedProfileByName(name string) (named_profile.NamedProfile, error)
	AddNamedProfile(namedProfile named_profile.NamedProfile) error
}

type NamedConfigurationsFacade interface {
	GetNamedConfigurations() []named_configuration.NamedConfiguration
	GetNamedConfigurationById(id string) (named_configuration.NamedConfiguration, error)
	GetNamedConfigurationByName(name string) (named_configuration.NamedConfiguration, error)
	AddNamedConfiguration(namedProfile named_configuration.NamedConfiguration) error
}

type NamedProfilesActionsInterface interface {
	GetOrCreateNamedProfile(profileName string) (named_profile.NamedProfile, error)
}

type NamedConfigurationsActionsInterface interface {
	GetOrCreateNamedConfiguration(configurationName string) (named_configuration.NamedConfiguration, error)
}

type AwsIamUserSessionsFacade interface {
	GetSessions() []aws_iam_user.AwsIamUserSession
	GetSessionById(sessionId string) (aws_iam_user.AwsIamUserSession, error)
	AddSession(session aws_iam_user.AwsIamUserSession) error
	RemoveSession(sessionId string) error
	EditSession(sessionId string, sessionName string, region string, accountNumber string, userName string,
		mfaDevice string, namedProfileId string) error
	SetSessionTokenExpiration(sessionId string, sessionTokenExpiration string) error
	StartingSession(sessionId string) error
	StartSession(sessionId string, startTime string) error
	StopSession(sessionId string, stopTime string) error
}

type GcpIamUserAccountOauthSessionsFacade interface {
	GetSessions() []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession
	GetSessionById(sessionId string) (gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession, error)
	AddSession(session gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession) error
	StartSession(sessionId string, startTime string) error
	StopSession(sessionId string, stopTime string) error
	RemoveSession(sessionId string) error
	EditSession(sessionId string, sessionName string, projectName string, namedConfigurationId string) error
}
