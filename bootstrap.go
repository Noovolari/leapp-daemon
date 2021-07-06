package main

import (
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

func NamedProfilesBootstrap(prov *providers.Providers, config domain.Configuration) {
	namedProfilesFacade := prov.GetNamedProfilesFacade()
	namedProfiles := config.NamedProfiles
	namedProfilesFacade.SetNamedProfiles(namedProfiles)
	namedProfilesFacade.Subscribe(prov.GetNamedProfilesWriter())
}
