package use_case

import (
	"leapp_daemon/domain/domain_aws/aws_iam_user"
	"leapp_daemon/infrastructure/logging"
)

type AwsSessionsWriter struct {
	ConfigurationRepository ConfigurationRepository
}

func (sessionWriter *AwsSessionsWriter) UpdateAwsIamUserSessions(oldSessions []aws_iam_user.AwsIamUserSession, newSessions []aws_iam_user.AwsIamUserSession) {
	config, err := sessionWriter.ConfigurationRepository.GetConfiguration()
	if err != nil {
		logging.Entry().Error(err)
		return
	}

	config.AwsIamUserSessions = newSessions
	err = sessionWriter.ConfigurationRepository.UpdateConfiguration(config)
	if err != nil {
		logging.Entry().Error(err)
		return
	}
}
