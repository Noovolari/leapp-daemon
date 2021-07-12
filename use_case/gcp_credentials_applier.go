package use_case

import (
	"leapp_daemon/domain/domain_gcp"
	"leapp_daemon/domain/domain_gcp/gcp_iam_user_account_oauth"
	"leapp_daemon/infrastructure/logging"
	"reflect"
)

type GcpCredentialsApplier struct {
	Keychain   Keychain
	Repository GcpConfigurationRepository
}

func (applier *GcpCredentialsApplier) UpdateGcpIamUserAccountOauthSessions(oldSessions []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession, newSessions []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession) {
	oldActiveSession, newActiveSession := applier.getActiveSessions(oldSessions, newSessions)

	if oldActiveSession != nil {
		if newActiveSession == nil {
			applier.deactivateSession(oldActiveSession)
		} else if oldActiveSession.Id != newActiveSession.Id {
			applier.deactivateSession(oldActiveSession)
			applier.activateSession(newActiveSession)
		} else {
			if !reflect.DeepEqual(oldActiveSession, newActiveSession) {
				applier.activateSession(newActiveSession)
			}
		}
	} else {
		if newActiveSession != nil {
			applier.activateSession(newActiveSession)
		}
	}
}

func (applier *GcpCredentialsApplier) getActiveSessions(oldSessions []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession, newSessions []gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession) (*gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession, *gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession) {
	var oldActiveSession *gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession
	for _, oldSession := range oldSessions {
		if oldSession.Status == domain_gcp.Active {
			oldActiveSession = &oldSession
			break
		}
	}
	var newActiveSession *gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession
	for _, newSession := range newSessions {
		if newSession.Status == domain_gcp.Active {
			newActiveSession = &newSession
			break
		}
	}
	return oldActiveSession, newActiveSession
}

func (applier *GcpCredentialsApplier) activateSession(session *gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession) {
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

	err = applier.Repository.CreateConfiguration(session.AccountId, session.ProjectName)
	if err != nil {
		logging.Entry().Error(err)
		return
	}

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

func (applier *GcpCredentialsApplier) deactivateSession(session *gcp_iam_user_account_oauth.GcpIamUserAccountOauthSession) {
	err := applier.Repository.RemoveDefaultCredentials()
	if err != nil {
		logging.Entry().Error(err)
	}

	err = applier.Repository.DeactivateConfiguration()
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

	err = applier.Repository.RemoveConfiguration()
	if err != nil {
		logging.Entry().Error(err)
	}
}
