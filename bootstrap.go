package main

import (
	"fmt"
	"leapp_daemon/domain/configuration"
	"leapp_daemon/infrastructure/logging"
	"leapp_daemon/providers"
)

func ConfigurationBootstrap(prov *providers.Providers) configuration.Configuration {
	config, err := prov.GetFileConfigurationRepository().GetConfiguration()
	if err != nil {
		logging.Entry().Error(err)
		panic(err)
	}
	return config
}

func AwsIamUserBootstrap(prov *providers.Providers, config configuration.Configuration) {
	fileSystem := prov.GetFileSystem()
	homeDir, err := fileSystem.GetHomeDir()
	if err != nil {
		logging.Entry().Error(err)
		panic(err)
	}

	configurationFileBackupPath := fmt.Sprintf("%s/%s", homeDir, ".Leapp/Leapp-lock.json")
	doesConfigurationFileExist := fileSystem.DoesFileExist(configurationFileBackupPath)

	//AWS Stuff
	credentialsFilePath := fmt.Sprintf("%s/%s", homeDir, ".aws/credentials")
	doesCredentialsFileExist := fileSystem.DoesFileExist(credentialsFilePath)

	credentialsFileBackupPath := fmt.Sprintf("%s/%s", homeDir, ".aws/credentials.leapp.bkp")
	doesCredentialsFileBackupExist := fileSystem.DoesFileExist(credentialsFileBackupPath)

	// TODO: move this logic to aws_credentials_applier and don't use doesConfigurationFileExist
	if !doesConfigurationFileExist && doesCredentialsFileExist && !doesCredentialsFileBackupExist {
		err = fileSystem.RenameFile(credentialsFilePath, credentialsFileBackupPath)
		if err != nil {
			logging.Entry().Error(err)
			panic(err)
		}
	}

	awsIamUserSessionFacade := prov.GetAwsIamUserSessionFacade()
	awsIamUserSessions := config.AwsIamUserSessions
	awsIamUserSessionFacade.SetSessions(awsIamUserSessions)
	awsIamUserSessionFacade.Subscribe(prov.GetAwsSessionWriter())
	awsIamUserSessionFacade.Subscribe(prov.GetAwsCredentialsApplier())
	logging.Info(fmt.Sprintf("%+v", awsIamUserSessions))
}

func GcpIamUserAccountOauthBootstrap(prov *providers.Providers, config configuration.Configuration) {
	gcpIamUserAccountOauthSessionFacade := prov.GetGcpIamUserAccountOauthSessionFacade()
	gcpIamUserAccountOauthSessions := config.GcpIamUserAccountOauthSessions
	gcpIamUserAccountOauthSessionFacade.SetSessions(gcpIamUserAccountOauthSessions)
	gcpIamUserAccountOauthSessionFacade.Subscribe(prov.GetGcpSessionWriter())
	gcpIamUserAccountOauthSessionFacade.Subscribe(prov.GetGcpCredentialsApplier())
	logging.Info(fmt.Sprintf("%+v", gcpIamUserAccountOauthSessions))
}

func AlibabaRamUserBootstrap(prov *providers.Providers, config configuration.Configuration) {
	alibabaRamUserSessionFacade := prov.GetAlibabaRamUserSessionFacade()
	alibabaRamUserSessions := config.AlibabaRamUserSessions
	alibabaRamUserSessionFacade.SetSessions(alibabaRamUserSessions)
	alibabaRamUserSessionFacade.Subscribe(prov.GetAlibabaSessionWriter())
	alibabaRamUserSessionFacade.Subscribe(prov.GetAlibabaCredentialsApplier())
	logging.Info(fmt.Sprintf("%+v", alibabaRamUserSessions))
}

func AlibabaRamRoleFederatedBootstrap(prov *providers.Providers, config configuration.Configuration) {
	alibabaRamRoleFederatedSessionFacade := prov.GetAlibabaRamRoleFederatedSessionFacade()
	alibabaRamRoleFederatedSessions := config.AlibabaRamRoleFederatedSessions
	alibabaRamRoleFederatedSessionFacade.SetSessions(alibabaRamRoleFederatedSessions)
	alibabaRamRoleFederatedSessionFacade.Subscribe(prov.GetAlibabaSessionWriter())
	alibabaRamRoleFederatedSessionFacade.Subscribe(prov.GetAlibabaCredentialsApplier())
	logging.Info(fmt.Sprintf("%+v", alibabaRamRoleFederatedSessions))
}

func AlibabaRamRoleChainedBootstrap(prov *providers.Providers, config configuration.Configuration) {
	alibabaRamRoleChainedSessionFacade := prov.GetAlibabaRamRoleChainedSessionFacade()
	alibabaRamRoleChainedSessions := config.AlibabaRamRoleChainedSessions
	alibabaRamRoleChainedSessionFacade.SetSessions(alibabaRamRoleChainedSessions)
	alibabaRamRoleChainedSessionFacade.Subscribe(prov.GetAlibabaSessionWriter())
	alibabaRamRoleChainedSessionFacade.Subscribe(prov.GetAlibabaCredentialsApplier())
	logging.Info(fmt.Sprintf("%+v", alibabaRamRoleChainedSessions))
}

func NamedProfilesBootstrap(prov *providers.Providers, config configuration.Configuration) {
	namedProfilesFacade := prov.GetNamedProfilesFacade()
	namedProfiles := config.NamedProfiles
	namedProfilesFacade.SetNamedProfiles(namedProfiles)
	namedProfilesFacade.Subscribe(prov.GetNamedProfilesWriter())
	logging.Info(fmt.Sprintf("%+v", namedProfiles))
}
