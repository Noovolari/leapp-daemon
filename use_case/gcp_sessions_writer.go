package use_case

import (
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/logging"
)

type GcpSessionsWriter struct {
	ConfigurationRepository ConfigurationRepository
}

func (sessionWriter *GcpSessionsWriter) UpdateGcpIamUserAccountOauthSessions(oldSessions []session.GcpIamUserAccountOauthSession, newSessions []session.GcpIamUserAccountOauthSession) {
	config, err := sessionWriter.ConfigurationRepository.GetConfiguration()
	if err != nil {
		logging.Entry().Error(err)
		return
	}

	config.GcpIamUserAccountOauthSessions = newSessions
	err = sessionWriter.ConfigurationRepository.UpdateConfiguration(config)
	if err != nil {
		logging.Entry().Error(err)
		return
	}
}
