package use_case

import (
	"leapp_daemon/domain/domain_alibaba/alibaba_ram_role_chained"
	"leapp_daemon/domain/domain_alibaba/alibaba_ram_role_federated"
	"leapp_daemon/domain/domain_alibaba/alibaba_ram_user"
)

type AlibabaSessionsWriter struct {
	ConfigurationRepository ConfigurationRepository
}

func (sessionWriter *AlibabaSessionsWriter) UpdateAlibabaRamUserSessions(oldSessions []alibaba_ram_user.AlibabaRamUserSession, newSessions []alibaba_ram_user.AlibabaRamUserSession) error {
	config, err := sessionWriter.ConfigurationRepository.GetConfiguration()
	if err != nil {
		return err
	}

	config.AlibabaRamUserSessions = newSessions
	err = sessionWriter.ConfigurationRepository.UpdateConfiguration(config)

	return err
}

func (sessionWriter *AlibabaSessionsWriter) UpdateAlibabaRamRoleFederatedSessions(oldSessions []alibaba_ram_role_federated.AlibabaRamRoleFederatedSession, newSessions []alibaba_ram_role_federated.AlibabaRamRoleFederatedSession) error {
	config, err := sessionWriter.ConfigurationRepository.GetConfiguration()
	if err != nil {
		return err
	}

	config.AlibabaRamRoleFederatedSessions = newSessions
	err = sessionWriter.ConfigurationRepository.UpdateConfiguration(config)

	return err
}

func (sessionWriter *AlibabaSessionsWriter) UpdateAlibabaRamRoleChainedSessions(oldSessions []alibaba_ram_role_chained.AlibabaRamRoleChainedSession, newSessions []alibaba_ram_role_chained.AlibabaRamRoleChainedSession) error {
	config, err := sessionWriter.ConfigurationRepository.GetConfiguration()
	if err != nil {
		return err
	}

	config.AlibabaRamRoleChainedSessions = newSessions
	err = sessionWriter.ConfigurationRepository.UpdateConfiguration(config)

	return err
}
