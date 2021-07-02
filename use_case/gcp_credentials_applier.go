package use_case

import (
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/logging"
	"reflect"
)

type GcpCredentialsApplier struct {
	Keychain   Keychain
	Repository GcpConfigurationRepository
}

func (applier *GcpCredentialsApplier) UpdateGcpIamUserAccountOauthSessions(oldSessions []session.GcpIamUserAccountOauthSession, newSessions []session.GcpIamUserAccountOauthSession) {
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

func (applier *GcpCredentialsApplier) getActiveSessions(oldSessions []session.GcpIamUserAccountOauthSession, newSessions []session.GcpIamUserAccountOauthSession) (*session.GcpIamUserAccountOauthSession, *session.GcpIamUserAccountOauthSession) {
	var oldActiveSession *session.GcpIamUserAccountOauthSession
	for _, oldSession := range oldSessions {
		if oldSession.Status == session.Active {
			oldActiveSession = &oldSession
			break
		}
	}
	var newActiveSession *session.GcpIamUserAccountOauthSession
	for _, newSession := range newSessions {
		if newSession.Status == session.Active {
			newActiveSession = &newSession
			break
		}
	}
	return oldActiveSession, newActiveSession
}

func (applier *GcpCredentialsApplier) activateSession(session *session.GcpIamUserAccountOauthSession) {
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

func (applier *GcpCredentialsApplier) deactivateSession(session *session.GcpIamUserAccountOauthSession) {
	err := applier.Repository.RemoveDefaultCredentials()
	if err != nil {
		logging.Entry().Error(err)
		return
	}

	err = applier.Repository.DeactivateConfiguration()
	if err != nil {
		logging.Entry().Error(err)
		return
	}

	err = applier.Repository.RemoveCredentialsFromDb(session.AccountId)
	if err != nil {
		logging.Entry().Error(err)
		return
	}

	err = applier.Repository.RemoveAccessTokensFromDb(session.AccountId)
	if err != nil {
		logging.Entry().Error(err)
		return
	}

	err = applier.Repository.RemoveConfiguration()
	if err != nil {
		logging.Entry().Error(err)
		return
	}
}
