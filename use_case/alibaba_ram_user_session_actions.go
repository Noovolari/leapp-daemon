package use_case

import (
	"fmt"
	"leapp_daemon/domain/domain_alibaba"
	"leapp_daemon/domain/domain_alibaba/alibaba_ram_user"
	"leapp_daemon/infrastructure/http/http_error"
)

type AlibabaRamUserSessionActions struct {
	Environment                  Environment
	Keychain                     Keychain
	NamedProfilesActions         NamedProfilesActionsInterface
	AlibabaRamUserSessionsFacade AlibabaRamUserSessionsFacade
}

func (actions *AlibabaRamUserSessionActions) Create(name string, alibabaAccessKeyId string, alibabaSecretAccessKey string, regionName string, profileName string) error {
	namedProfile, err := actions.NamedProfilesActions.GetOrCreateNamedProfile(profileName)
	if err != nil {
		return err
	}

	isRegionValid := domain_alibaba.IsAlibabaRegionValid(regionName)
	if !isRegionValid {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("Region " + regionName + " not valid"))
	}

	sess := alibaba_ram_user.AlibabaRamUserSession{
		Id:             actions.Environment.GenerateUuid(),
		Name:           name,
		Status:         domain_alibaba.NotActive,
		Region:         regionName,
		NamedProfileId: namedProfile.Id,
	}

	err = actions.AlibabaRamUserSessionsFacade.AddSession(sess)
	if err != nil {
		return err
	}

	err = actions.Keychain.SetSecret(alibabaAccessKeyId, sess.Id+domain_alibaba.PlainAlibabaKeyIdSuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = actions.Keychain.SetSecret(alibabaSecretAccessKey, sess.Id+domain_alibaba.PlainAlibabaSecretAccessKeySuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func (actions *AlibabaRamUserSessionActions) Get(id string) (*alibaba_ram_user.AlibabaRamUserSession, error) {
	var sess *alibaba_ram_user.AlibabaRamUserSession
	sess, err := actions.AlibabaRamUserSessionsFacade.GetSessionById(id)
	return sess, err
}

func (actions *AlibabaRamUserSessionActions) Update(id string, name string, regionName string,
	alibabaAccessKeyId string, alibabaSecretAccessKey string, profileName string) error {

	isRegionValid := domain_alibaba.IsAlibabaRegionValid(regionName)
	if !isRegionValid {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("Region " + regionName + " not valid"))
	}

	np, err := actions.NamedProfilesActions.GetOrCreateNamedProfile(profileName)
	if err != nil {
		return err //TODO: return right error
	}

	/*sess := alibaba_ram_user.AlibabaRamUserSession{
		Id:             id,
		Name:           name,
		Status:         domain_alibaba.NotActive,
		Region:         regionName,
		NamedProfileId: np.Id,
	}*/

	err = actions.AlibabaRamUserSessionsFacade.EditSession(id, name, regionName, np.Id)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = actions.Keychain.SetSecret(alibabaAccessKeyId, id+domain_alibaba.PlainAlibabaKeyIdSuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = actions.Keychain.SetSecret(alibabaSecretAccessKey, id+domain_alibaba.PlainAlibabaSecretAccessKeySuffix)
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

	if sess.Status != domain_alibaba.NotActive {
		err = actions.Stop(id)
		if err != nil {
			return err
		}
	}

	err = actions.AlibabaRamUserSessionsFacade.RemoveSession(id)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = actions.Keychain.DeleteSecret(id + domain_alibaba.PlainAlibabaKeyIdSuffix)
	if err != nil {
		return err
	}

	err = actions.Keychain.DeleteSecret(id + domain_alibaba.PlainAlibabaSecretAccessKeySuffix)
	if err != nil {
		return err
	}
	return nil
}

func (actions *AlibabaRamUserSessionActions) Start(sessionId string) error {

	err := actions.AlibabaRamUserSessionsFacade.StartingSession(sessionId)
	if err != nil {
		return err
	}

	err = actions.AlibabaRamUserSessionsFacade.StartSession(sessionId)
	if err != nil {
		return err
	}

	return nil
}

func (actions *AlibabaRamUserSessionActions) Stop(sessionId string) error {

	err := actions.AlibabaRamUserSessionsFacade.StopSession(sessionId)
	if err != nil {
		return err
	}

	return nil
}
