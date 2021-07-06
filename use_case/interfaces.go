package use_case

import (
	"leapp_daemon/domain/configuration"
	"leapp_daemon/domain/named_profile"
	"leapp_daemon/domain/session"
	"time"

	"golang.org/x/oauth2"
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

type GcpApi interface {
	GetOauthUrl() (string, error)
	GetOauthToken(authCode string) (*oauth2.Token, error)
	GetCredentials(oauthToken *oauth2.Token) string
}

type ConfigurationRepository interface {
	CreateConfiguration(configuration.Configuration) error
	GetConfiguration() (configuration.Configuration, error)
	UpdateConfiguration(configuration.Configuration) error
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
	UpdateNamedProfileName(id string, name string) error
	DeleteNamedProfile(profileId string) error
}

type NamedProfilesActionsInterface interface {
	GetOrCreateNamedProfile(profileName string) (named_profile.NamedProfile, error)
	GetNamedProfileById(profileId string) (named_profile.NamedProfile, error)
	UpdateNamedProfileName(id string, name string) error
	DeleteNamedProfile(profileId string) error
}

type AwsIamUserSessionsFacade interface {
	Subscribe(observer session.AwsIamUserSessionsObserver)
	GetSessions() []session.AwsIamUserSession
	SetSessions(sessions []session.AwsIamUserSession)
	AddSession(session session.AwsIamUserSession) error
	RemoveSession(id string) error
	GetSessionById(id string) (*session.AwsIamUserSession, error)
	SetSessionStatusToPending(id string) error
	SetSessionStatusToActive(id string) error
	SetSessionTokenExpiration(sessionId string, sessionTokenExpiration time.Time) error
}

type GcpIamUserAccountOauthSessionsFacade interface {
	GetSessions() []session.GcpIamUserAccountOauthSession
	GetSessionById(sessionId string) (session.GcpIamUserAccountOauthSession, error)
	AddSession(session session.GcpIamUserAccountOauthSession) error
	StartSession(sessionId string, startTime string) error
	StopSession(sessionId string, stopTime string) error
	RemoveSession(sessionId string) error
	EditSession(sessionId string, name string, projectName string) error
}

type AlibabaRamUserSessionsFacade interface {
	Subscribe(observer session.AlibabaRamUserSessionsObserver)
	GetSessions() []session.AlibabaRamUserSession
	SetSessions(sessions []session.AlibabaRamUserSession) error
	AddSession(session session.AlibabaRamUserSession) error
	RemoveSession(id string) error
	UpdateSession(session session.AlibabaRamUserSession) error
	GetSessionById(id string) (*session.AlibabaRamUserSession, error)
	SetSessionStatusToPending(id string) error
	SetSessionStatusToActive(id string) error
	SetSessionStatusToInactive(id string) error
}

type AlibabaRamRoleFederatedSessionsFacade interface {
	Subscribe(observer session.AlibabaRamRoleFederatedSessionsObserver)
	GetSessions() []session.AlibabaRamRoleFederatedSession
	SetSessions(sessions []session.AlibabaRamRoleFederatedSession) error
	AddSession(session session.AlibabaRamRoleFederatedSession) error
	RemoveSession(id string) error
	UpdateSession(session session.AlibabaRamRoleFederatedSession) error
	GetSessionById(id string) (*session.AlibabaRamRoleFederatedSession, error)
	SetSessionStatusToPending(id string) error
	SetSessionStatusToActive(id string) error
	SetSessionStatusToInactive(id string) error
}

type AlibabaRamRoleChainedSessionsFacade interface {
	Subscribe(observer session.AlibabaRamRoleChainedSessionsObserver)
	GetSessions() []session.AlibabaRamRoleChainedSession
	SetSessions(sessions []session.AlibabaRamRoleChainedSession) error
	AddSession(session session.AlibabaRamRoleChainedSession) error
	RemoveSession(id string) error
	GetSessionById(id string) (*session.AlibabaRamRoleChainedSession, error)
	SetSessionById(newSession *session.AlibabaRamRoleChainedSession) error
	SetSessionStatusToPending(id string) error
	SetSessionStatusToActive(id string) error
	SetSessionStatusToInactive(id string) error
}
