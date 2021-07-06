package use_case

import (
	"fmt"
	"leapp_daemon/domain/constant"
	"leapp_daemon/domain/region"
	"leapp_daemon/domain/session"
	"leapp_daemon/infrastructure/http/http_error"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
)

type AlibabaRamRoleFederatedSessionActions struct {
	Environment                           Environment
	Keychain                              Keychain
	NamedProfilesActions                  NamedProfilesActionsInterface
	AlibabaRamRoleFederatedSessionsFacade AlibabaRamRoleFederatedSessionsFacade
}

// TODO: mettere da qualche parte questa funzione
func SAMLAuth(region string, idpArn string, roleArn string, assertion string) (key string, secret string, token string, err error) {
	// I'm using this since NewClient() method returns a panic saying literally "not support yet"
	// This method actually never use the credentials so I placed 2 placeholders
	client, _ := sts.NewClientWithAccessKey(region, "", "")

	request := sts.CreateAssumeRoleWithSAMLRequest()
	request.Scheme = "https"
	request.SAMLProviderArn = idpArn
	request.RoleArn = roleArn
	request.SAMLAssertion = assertion
	response, err := client.AssumeRoleWithSAML(request)
	if err != nil {
		return "", "", "", err
	}
	key = response.Credentials.AccessKeyId
	secret = response.Credentials.AccessKeySecret
	token = response.Credentials.SecurityToken
	return
}

func (actions *AlibabaRamRoleFederatedSessionActions) Create(name string, roleName string, roleArn string,
	idpArn string, regionName string, ssoUrl string, profileName string) error {

	namedProfile, err := actions.NamedProfilesActions.GetOrCreateNamedProfile(profileName)
	if err != nil {
		return err
	}

	isRegionValid := region.IsAlibabaRegionValid(regionName)
	if !isRegionValid {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("Region " + regionName + " not valid"))
	}

	alibabaRole := session.AlibabaRamRole{
		Name: roleName,
		Arn:  roleArn,
	}

	federatedAlibabaAccount := session.AlibabaRamRoleFederatedAccount{
		Name:          name,
		Role:          &alibabaRole,
		IdpArn:        idpArn,
		Region:        regionName,
		/*SsoUrl:        ssoUrl,*/
		NamedProfileId: namedProfile.Id,
	}

	sess := session.AlibabaRamRoleFederatedSession{
		Id:      actions.Environment.GenerateUuid(),
		Status:  session.NotActive,
		Account: &federatedAlibabaAccount,
	}

	err = actions.AlibabaRamRoleFederatedSessionsFacade.AddSession(sess)
	if err != nil {
		return err
	}

	alibabaAccessKeyId, alibabaSecretAccessKey, alibabaStsToken, err := SAMLAuth(regionName, idpArn, roleArn, ssoUrl)
	if err != nil {
		return err
	}

	err = actions.Keychain.SetSecret(alibabaAccessKeyId, sess.Id+constant.FederatedAlibabaKeyIdSuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = actions.Keychain.SetSecret(alibabaSecretAccessKey, sess.Id+constant.FederatedAlibabaSecretAccessKeySuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = actions.Keychain.SetSecret(alibabaStsToken, sess.Id+constant.FederatedAlibabaStsTokenSuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func (actions *AlibabaRamRoleFederatedSessionActions) Get(id string) (*session.AlibabaRamRoleFederatedSession, error) {
	return actions.AlibabaRamRoleFederatedSessionsFacade.GetSessionById(id)
}

func (actions *AlibabaRamRoleFederatedSessionActions) Update(id string, name string, roleName string, roleArn string,
	idpArn string, regionName string, ssoUrl string, profileName string) error {

	isRegionValid := region.IsAlibabaRegionValid(regionName)
	if !isRegionValid {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("Region " + regionName + " not valid"))
	}

	np, err := actions.NamedProfilesActions.GetOrCreateNamedProfile(profileName)
	if err != nil {
		return err //TODO: return right error
	}

	alibabaRole := session.AlibabaRamRole{
		Name: roleName,
		Arn:  roleArn,
	}

	federatedAlibabaAccount := session.AlibabaRamRoleFederatedAccount{
		Name:          name,
		Role:          &alibabaRole,
		IdpArn:        idpArn,
		Region:        regionName,
		/*SsoUrl:        ssoUrl,*/
		NamedProfileId: np.Id,
	}

	sess := session.AlibabaRamRoleFederatedSession{
		Id:      id,
		Status:  session.NotActive,
		Account: &federatedAlibabaAccount,
	}

	err = actions.AlibabaRamRoleFederatedSessionsFacade.UpdateSession(sess)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	alibabaAccessKeyId, alibabaSecretAccessKey, alibabaStsToken, err := SAMLAuth(regionName, idpArn, roleArn, ssoUrl)
	if err != nil {
		return err
	}

	err = actions.Keychain.SetSecret(alibabaAccessKeyId, id+constant.FederatedAlibabaKeyIdSuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = actions.Keychain.SetSecret(alibabaSecretAccessKey, id+constant.FederatedAlibabaSecretAccessKeySuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = actions.Keychain.SetSecret(alibabaStsToken, id+constant.FederatedAlibabaStsTokenSuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func (actions *AlibabaRamRoleFederatedSessionActions) Delete(id string) error {
	sess, err := session.GetAlibabaRamRoleFederatedSessionsFacade().GetSessionById(id)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	if sess.Status != session.NotActive {
		err = actions.Stop(id)
		if err != nil {
			return err
		}
	}

	err = actions.AlibabaRamRoleFederatedSessionsFacade.RemoveSession(id)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = actions.Keychain.DeleteSecret(id + constant.FederatedAlibabaKeyIdSuffix)
	if err != nil {
		return err
	}

	err = actions.Keychain.DeleteSecret(id + constant.FederatedAlibabaSecretAccessKeySuffix)
	if err != nil {
		return err
	}

	err = actions.Keychain.DeleteSecret(id + constant.FederatedAlibabaStsTokenSuffix)
	if err != nil {
		return err
	}
	return nil
}

func (actions *AlibabaRamRoleFederatedSessionActions) Start(sessionId string) error {

	err := actions.AlibabaRamRoleFederatedSessionsFacade.SetSessionStatusToPending(sessionId)
	if err != nil {
		return err
	}

	err = actions.AlibabaRamRoleFederatedSessionsFacade.SetSessionStatusToActive(sessionId)
	if err != nil {
		return err
	}

	return nil
}

func (actions *AlibabaRamRoleFederatedSessionActions) Stop(sessionId string) error {

	err := actions.AlibabaRamRoleFederatedSessionsFacade.SetSessionStatusToInactive(sessionId)
	if err != nil {
		return err
	}

	return nil
}
