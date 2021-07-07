package use_case

import (
	"leapp_daemon/domain/domain_gcp/gcp_iam_user_account_oauth"
	"leapp_daemon/infrastructure/logging"
)

type GcpSessionsWriter struct {
	ConfigurationRepository ConfigurationRepository
}

func (sessionWriter *GcpSessionsWriter) UpdateGcpIamUserAccountOauthSessions(oldSessions []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession, newSessions []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession) {
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
