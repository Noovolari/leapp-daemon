package use_case

import (
	"fmt"
	"leapp_daemon/domain/constant"
	"leapp_daemon/domain/region"
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/http/http_error"
)

type AlibabaRamUserSessionActions struct {
	Environment                  Environment
	Keychain                     Keychain
	NamedProfilesActions         NamedProfilesActionsInterface
	AlibabaRamUserSessionsFacade AlibabaRamUserSessionsFacade
}

func (actions *AlibabaRamUserSessionActions) Create(alias string, alibabaAccessKeyId string, alibabaSecretAccessKey string, regionName string, profileName string) error {
	namedProfile, err := actions.NamedProfilesActions.GetOrCreateNamedProfile(profileName)
	if err != nil {
		return err
	}

	isRegionValid := region.IsAlibabaRegionValid(regionName)
	if !isRegionValid {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("Region " + regionName + " not valid"))
	}

	alibabaRamUserAccount := session.AlibabaRamUserAccount{
		Region:         regionName,
		NamedProfileId: namedProfile.Id,
	}

	sess := session.AlibabaRamUserSession{
		Id:      actions.Environment.GenerateUuid(),
		Alias:   alias,
		Status:  session.NotActive,
		Account: &alibabaRamUserAccount,
	}

	err = actions.AlibabaRamUserSessionsFacade.AddSession(sess)
	if err != nil {
		return err
	}

	err = actions.Keychain.SetSecret(alibabaAccessKeyId, sess.Id+constant.PlainAlibabaKeyIdSuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = actions.Keychain.SetSecret(alibabaSecretAccessKey, sess.Id+constant.PlainAlibabaSecretAccessKeySuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func (actions *AlibabaRamUserSessionActions) Get(id string) (*session.AlibabaRamUserSession, error) {
	var sess *session.AlibabaRamUserSession
	sess, err := actions.AlibabaRamUserSessionsFacade.GetSessionById(id)
	return sess, err
}

func (actions *AlibabaRamUserSessionActions) Update(id string, alias string, regionName string,
	alibabaAccessKeyId string, alibabaSecretAccessKey string, profileName string) error {

	isRegionValid := region.IsAlibabaRegionValid(regionName)
	if !isRegionValid {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("Region " + regionName + " not valid"))
	}

	oldSess, err := session.GetAlibabaRamUserSessionsFacade().GetSessionById(id)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	plainAlibabaAccount := session.AlibabaRamUserAccount{
		Region:         regionName,
		NamedProfileId: oldSess.Account.NamedProfileId,
	}

	sess := session.AlibabaRamUserSession{
		Id:      id,
		Alias:   alias,
		Status:  session.NotActive,
		Account: &plainAlibabaAccount,
	}

	err = actions.AlibabaRamUserSessionsFacade.UpdateSession(sess)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = actions.NamedProfilesActions.UpdateNamedProfileName(oldSess.Account.NamedProfileId, profileName)
	if err != nil {
		return err //TODO: return right error
	}

	err = actions.Keychain.SetSecret(alibabaAccessKeyId, sess.Id+constant.PlainAlibabaKeyIdSuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = actions.Keychain.SetSecret(alibabaSecretAccessKey, sess.Id+constant.PlainAlibabaSecretAccessKeySuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func (actions *AlibabaRamUserSessionActions) Delete(id string) error {
	sess, err := actions.AlibabaRamUserSessionsFacade.GetSessionById(id)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	if sess.Status != session.NotActive {
		err = actions.Stop(id)
		if err != nil {
			return err
		}
	}

	oldSess, err := actions.AlibabaRamUserSessionsFacade.GetSessionById(id)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	oldNamedProfile, err := actions.NamedProfilesActions.GetNamedProfileById(oldSess.Account.NamedProfileId)
	if err != nil {
		return err //TODO: return right error
	}
	actions.NamedProfilesActions.DeleteNamedProfile(oldNamedProfile.Id)

	err = actions.AlibabaRamUserSessionsFacade.RemoveSession(id)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = actions.Keychain.DeleteSecret(id + constant.PlainAlibabaKeyIdSuffix)
	if err != nil {
		return err
	}

	err = actions.Keychain.DeleteSecret(id + constant.PlainAlibabaSecretAccessKeySuffix)
	if err != nil {
		return err
	}
	return nil
}

func (actions *AlibabaRamUserSessionActions) Start(sessionId string) error {

	err := actions.AlibabaRamUserSessionsFacade.SetSessionStatusToPending(sessionId)
	if err != nil {
		return err
	}

	err = actions.AlibabaRamUserSessionsFacade.SetSessionStatusToActive(sessionId)
	if err != nil {
		return err
	}

	return nil
}

func (actions *AlibabaRamUserSessionActions) Stop(sessionId string) error {

	err := actions.AlibabaRamUserSessionsFacade.SetSessionStatusToInactive(sessionId)
	if err != nil {
		return err
	}

	return nil
}
