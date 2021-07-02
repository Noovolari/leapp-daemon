package use_case

import (
	"leapp_daemon/domain/session"
)

type AlibabaSessionsWriter struct {
	ConfigurationRepository ConfigurationRepository
}

func (sessionWriter *AlibabaSessionsWriter) UpdateAlibabaRamUserSessions(oldSessions []session.AlibabaRamUserSession, newSessions []session.AlibabaRamUserSession) error {
	config, err := sessionWriter.ConfigurationRepository.GetConfiguration()
	if err != nil {
		return err
	}

	config.AlibabaRamUserSessions = newSessions
	err = sessionWriter.ConfigurationRepository.UpdateConfiguration(config)

	return err
}

func (sessionWriter *AlibabaSessionsWriter) UpdateAlibabaRamRoleFederatedSessions(oldSessions []session.AlibabaRamRoleFederatedSession, newSessions []session.AlibabaRamRoleFederatedSession) error {
	config, err := sessionWriter.ConfigurationRepository.GetConfiguration()
	if err != nil {
		return err
	}

	config.AlibabaRamRoleFederatedSessions = newSessions
	err = sessionWriter.ConfigurationRepository.UpdateConfiguration(config)

	return err
}

func (sessionWriter *AlibabaSessionsWriter) UpdateAlibabaRamRoleChainedSessions(oldSessions []session.AlibabaRamRoleChainedSession, newSessions []session.AlibabaRamRoleChainedSession) error {
	config, err := sessionWriter.ConfigurationRepository.GetConfiguration()
	if err != nil {
		return err
	}

	config.AlibabaRamRoleChainedSessions = newSessions
	err = sessionWriter.ConfigurationRepository.UpdateConfiguration(config)

	return err
}
