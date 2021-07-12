package main

import (
	"fmt"
	"leapp_daemon/domain"
	"leapp_daemon/infrastructure/logging"
	"leapp_daemon/providers"
)

func ConfigurationBootstrap(prov *providers.Providers) domain.Configuration {
	config, err := prov.GetFileConfigurationRepository().GetConfiguration()
	if err != nil {
		logging.Entry().Error(err)
		panic(err)
	}
	return config
}

func AwsIamUserBootstrap(prov *providers.Providers, config domain.Configuration) {
	awsIamUserSessionFacade := prov.GetAwsIamUserSessionFacade()
	awsIamUserSessions := config.AwsIamUserSessions
	awsIamUserSessionFacade.SetSessions(awsIamUserSessions)
	awsIamUserSessionFacade.Subscribe(prov.GetAwsSessionWriter())
	awsIamUserSessionFacade.Subscribe(prov.GetAwsCredentialsApplier())
	prov.GetTimerCollection().AddTimer(1000,
		prov.GetAwsIamUserSessionActions().RotateSessionTokens)
}

func GcpIamUserAccountOauthBootstrap(prov *providers.Providers, config domain.Configuration) {
	gcpIamUserAccountOauthSessionFacade := prov.GetGcpIamUserAccountOauthSessionFacade()
	gcpIamUserAccountOauthSessions := config.GcpIamUserAccountOauthSessions
	gcpIamUserAccountOauthSessionFacade.SetSessions(gcpIamUserAccountOauthSessions)
	gcpIamUserAccountOauthSessionFacade.Subscribe(prov.GetGcpSessionWriter())
	gcpIamUserAccountOauthSessionFacade.Subscribe(prov.GetGcpCredentialsApplier())
}

func AlibabaRamUserBootstrap(prov *providers.Providers, config domain.Configuration) {
	alibabaRamUserSessionFacade := prov.GetAlibabaRamUserSessionFacade()
	alibabaRamUserSessions := config.AlibabaRamUserSessions
	alibabaRamUserSessionFacade.SetSessions(alibabaRamUserSessions)
	alibabaRamUserSessionFacade.Subscribe(prov.GetAlibabaSessionWriter())
	alibabaRamUserSessionFacade.Subscribe(prov.GetAlibabaCredentialsApplier())
	logging.Info(fmt.Sprintf("%+v", alibabaRamUserSessions))
}

func AlibabaRamRoleFederatedBootstrap(prov *providers.Providers, config domain.Configuration) {
	alibabaRamRoleFederatedSessionFacade := prov.GetAlibabaRamRoleFederatedSessionFacade()
	alibabaRamRoleFederatedSessions := config.AlibabaRamRoleFederatedSessions
	alibabaRamRoleFederatedSessionFacade.SetSessions(alibabaRamRoleFederatedSessions)
	alibabaRamRoleFederatedSessionFacade.Subscribe(prov.GetAlibabaSessionWriter())
	alibabaRamRoleFederatedSessionFacade.Subscribe(prov.GetAlibabaCredentialsApplier())
	logging.Info(fmt.Sprintf("%+v", alibabaRamRoleFederatedSessions))
}

func AlibabaRamRoleChainedBootstrap(prov *providers.Providers, config domain.Configuration) {
	alibabaRamRoleChainedSessionFacade := prov.GetAlibabaRamRoleChainedSessionFacade()
	alibabaRamRoleChainedSessions := config.AlibabaRamRoleChainedSessions
	alibabaRamRoleChainedSessionFacade.SetSessions(alibabaRamRoleChainedSessions)
	alibabaRamRoleChainedSessionFacade.Subscribe(prov.GetAlibabaSessionWriter())
	alibabaRamRoleChainedSessionFacade.Subscribe(prov.GetAlibabaCredentialsApplier())
	logging.Info(fmt.Sprintf("%+v", alibabaRamRoleChainedSessions))
}

func NamedProfilesBootstrap(prov *providers.Providers, config domain.Configuration) {
	namedProfilesFacade := prov.GetNamedProfilesFacade()
	namedProfiles := config.NamedProfiles
	namedProfilesFacade.SetNamedProfiles(namedProfiles)
	namedProfilesFacade.Subscribe(prov.GetNamedProfilesWriter())
}
