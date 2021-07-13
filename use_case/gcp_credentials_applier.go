package use_case

import (
	"leapp_daemon/domain/domain_gcp"
	"leapp_daemon/domain/domain_gcp/gcp_iam_user_account_oauth"
	"leapp_daemon/infrastructure/logging"
)

type GcpCredentialsApplier struct {
	Keychain                  Keychain
	Repository                GcpConfigurationRepository
	NamedConfigurationsFacade NamedConfigurationsFacade
}

func (applier *GcpCredentialsApplier) UpdateGcpIamUserAccountOauthSessions(oldSessions []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession, newSessions []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession) {
	for _, oldSession := range oldSessions {
		applier.deactivateSession(&oldSession)
	}

	for _, newSession := range newSessions {
		applier.activateSession(&newSession)
	}
}

func (applier *GcpCredentialsApplier) activateSession(session *gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession) {
	if session.Status != domain_gcp.Active {
		return
	}

	credentials, err := applier.Keychain.GetSecret(session.CredentialsLabel)
	if err != nil {
		logging.Entry().Error(err)
		return
	}

	err = applier.Repository.WriteCredentialsToDb(session.AccountId, credentials)
	if err != nil {
		logging.Entry().Error(err)
		return
	}

	namedConfiguration, err := applier.NamedConfigurationsFacade.GetNamedConfigurationById(session.NamedConfigurationId)
	if err != nil {
		logging.Entry().Error(err)
	}

	err = applier.Repository.CreateConfiguration(session.AccountId, session.ProjectName, namedConfiguration.Name)
	if err != nil {
		logging.Entry().Error(err)
		return
	}

	if namedConfiguration.IsDefault() {
		err = applier.Repository.ActivateConfiguration()
		if err != nil {
			logging.Entry().Error(err)
			return
		}

		err = applier.Repository.WriteDefaultCredentials(credentials)
		if err != nil {
			logging.Entry().Error(err)
			return
		}
	}
}

func (applier *GcpCredentialsApplier) deactivateSession(session *gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession) {
	if session.Status != domain_gcp.Active {
		return
	}

	namedConfiguration, err := applier.NamedConfigurationsFacade.GetNamedConfigurationById(session.NamedConfigurationId)
	if err != nil {
		logging.Entry().Error(err)
	}

	err = applier.Repository.RemoveConfiguration(namedConfiguration.Name)
	if err != nil {
		logging.Entry().Error(err)
	}

	err = applier.Repository.RemoveCredentialsFromDb(session.AccountId)
	if err != nil {
		logging.Entry().Error(err)
	}

	err = applier.Repository.RemoveAccessTokensFromDb(session.AccountId)
	if err != nil {
		logging.Entry().Error(err)
	}

	if namedConfiguration.IsDefault() {
		err := applier.Repository.RemoveDefaultCredentials()
		if err != nil {
			logging.Entry().Error(err)
		}

		err = applier.Repository.DeactivateConfiguration()
		if err != nil {
			logging.Entry().Error(err)
		}
	}
}
