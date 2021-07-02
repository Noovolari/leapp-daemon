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
	Environment          Environment
	Keychain             Keychain
	NamedProfilesActions NamedProfilesActionsInterface
	/*AlibabaRamRoleFederatedSessionsFacade AlibabaRamUserSessionsFacade*/
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

func (actions *AlibabaRamRoleFederatedSessionActions) Create(name string, accountNumber string, roleName string, roleArn string,
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
		AccountNumber: accountNumber,
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
		Profile: profileName,
	}

	err = session.GetAlibabaRamRoleFederatedSessionsFacade().AddSession(sess)
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
	return session.GetAlibabaRamRoleFederatedSessionsFacade().GetSessionById(id)
}

func (actions *AlibabaRamRoleFederatedSessionActions) Update(id string, name string, accountNumber string, roleName string, roleArn string,
	idpArn string, regionName string, ssoUrl string, profileName string) error {

	isRegionValid := region.IsAlibabaRegionValid(regionName)
	if !isRegionValid {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("Region " + regionName + " not valid"))
	}

	oldSess, err := session.GetAlibabaRamRoleFederatedSessionsFacade().GetSessionById(id)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	alibabaRole := session.AlibabaRamRole{
		Name: roleName,
		Arn:  roleArn,
	}

	federatedAlibabaAccount := session.AlibabaRamRoleFederatedAccount{
		AccountNumber: accountNumber,
		Name:          name,
		Role:          &alibabaRole,
		IdpArn:        idpArn,
		Region:        regionName,
		/*SsoUrl:        ssoUrl,*/
		NamedProfileId: oldSess.Account.NamedProfileId,
	}

	sess := session.AlibabaRamRoleFederatedSession{
		Id:      id,
		Status:  session.NotActive,
		Account: &federatedAlibabaAccount,
		Profile: profileName,
	}

	err = session.GetAlibabaRamRoleFederatedSessionsFacade().UpdateSession(sess)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	oldNamedProfile, err := actions.NamedProfilesActions.GetNamedProfileById(oldSess.Account.NamedProfileId)
	if err != nil {
		return err //TODO: return right error
	}
	oldNamedProfile.Name = profileName
	actions.NamedProfilesActions.SetNamedProfileName(oldNamedProfile)

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

	oldSess, err := session.GetAlibabaRamRoleFederatedSessionsFacade().GetSessionById(id)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	oldNamedProfile, err := actions.NamedProfilesActions.GetNamedProfileById(oldSess.Account.NamedProfileId)
	if err != nil {
		return err //TODO: return right error
	}	
	actions.NamedProfilesActions.DeleteNamedProfile(oldNamedProfile.Id)

	err = session.GetAlibabaRamRoleFederatedSessionsFacade().RemoveSession(id)
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

	err := session.GetAlibabaRamRoleFederatedSessionsFacade().SetStatusToPending(sessionId)
	if err != nil {
		return err
	}

	err = session.GetAlibabaRamRoleFederatedSessionsFacade().SetStatusToActive(sessionId)
	if err != nil {
		return err
	}

	return nil
}

func (actions *AlibabaRamRoleFederatedSessionActions) Stop(sessionId string) error {

	err := session.GetAlibabaRamRoleFederatedSessionsFacade().SetStatusToInactive(sessionId)
	if err != nil {
		return err
	}

	return nil
}
