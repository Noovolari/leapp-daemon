package use_case

import (
	"leapp_daemon/domain/session"
)

type AwsSessionsWriter struct {
	ConfigurationRepository ConfigurationRepository
}

func (sessionWriter *AwsSessionsWriter) UpdateAwsIamUserSessions(oldSessions []session.AwsIamUserSession, newSessions []session.AwsIamUserSession) error {
	config, err := sessionWriter.ConfigurationRepository.GetConfiguration()
	if err != nil {
		return err
	}

	config.AwsIamUserSessions = newSessions
	err = sessionWriter.ConfigurationRepository.UpdateConfiguration(config)

	return err
}
