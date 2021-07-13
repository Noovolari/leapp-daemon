package use_case

import (
	"fmt"
	"leapp_daemon/domain/domain_alibaba"
	"leapp_daemon/domain/domain_alibaba/alibaba_ram_role_federated"
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
	idpArn string, regionName string, ssoUrl string, assertion string, profileName string) error {

	namedProfile, err := actions.NamedProfilesActions.GetOrCreateNamedProfile(profileName)
	if err != nil {
		return err
	}

	isRegionValid := domain_alibaba.IsAlibabaRegionValid(regionName)
	if !isRegionValid {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("Region " + regionName + " not valid"))
	}

	sess := alibaba_ram_role_federated.AlibabaRamRoleFederatedSession{
		Id:             actions.Environment.GenerateUuid(),
		Status:         domain_alibaba.NotActive,
		Name:           name,
		RoleName:       roleName,
		RoleArn:        roleArn,
		IdpArn:         idpArn,
		Region:         regionName,
		SsoUrl:         ssoUrl,
		NamedProfileId: namedProfile.Id,
	}

	err = actions.AlibabaRamRoleFederatedSessionsFacade.AddSession(sess)
	if err != nil {
		return err
	}

	alibabaAccessKeyId, alibabaSecretAccessKey, alibabaStsToken, err := SAMLAuth(regionName, idpArn, roleArn, assertion)
	if err != nil {
		return err
	}

	err = actions.Keychain.SetSecret(alibabaAccessKeyId, sess.Id+domain_alibaba.FederatedAlibabaKeyIdSuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = actions.Keychain.SetSecret(alibabaSecretAccessKey, sess.Id+domain_alibaba.FederatedAlibabaSecretAccessKeySuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = actions.Keychain.SetSecret(alibabaStsToken, sess.Id+domain_alibaba.FederatedAlibabaStsTokenSuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func (actions *AlibabaRamRoleFederatedSessionActions) Get(id string) (*alibaba_ram_role_federated.AlibabaRamRoleFederatedSession, error) {
	return actions.AlibabaRamRoleFederatedSessionsFacade.GetSessionById(id)
}

func (actions *AlibabaRamRoleFederatedSessionActions) Update(id string, name string, roleName string, roleArn string,
	idpArn string, regionName string, ssoUrl string, assertion string, profileName string) error {

	isRegionValid := domain_alibaba.IsAlibabaRegionValid(regionName)
	if !isRegionValid {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("Region " + regionName + " not valid"))
	}

	np, err := actions.NamedProfilesActions.GetOrCreateNamedProfile(profileName)
	if err != nil {
		return err //TODO: return right error
	}

	actions.AlibabaRamRoleFederatedSessionsFacade.EditSession(id, name, roleName, roleArn, idpArn, regionName, ssoUrl, np.Id)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	alibabaAccessKeyId, alibabaSecretAccessKey, alibabaStsToken, err := SAMLAuth(regionName, idpArn, roleArn, assertion)
	if err != nil {
		return err
	}

	err = actions.Keychain.SetSecret(alibabaAccessKeyId, id+domain_alibaba.FederatedAlibabaKeyIdSuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = actions.Keychain.SetSecret(alibabaSecretAccessKey, id+domain_alibaba.FederatedAlibabaSecretAccessKeySuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = actions.Keychain.SetSecret(alibabaStsToken, id+domain_alibaba.FederatedAlibabaStsTokenSuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func (actions *AlibabaRamRoleFederatedSessionActions) Delete(id string) error {
	sess, err := alibaba_ram_role_federated.GetAlibabaRamRoleFederatedSessionsFacade().GetSessionById(id)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	if sess.Status != domain_alibaba.NotActive {
		err = actions.Stop(id)
		if err != nil {
			return err
		}
	}

	err = actions.AlibabaRamRoleFederatedSessionsFacade.RemoveSession(id)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = actions.Keychain.DeleteSecret(id + domain_alibaba.FederatedAlibabaKeyIdSuffix)
	if err != nil {
		return err
	}

	err = actions.Keychain.DeleteSecret(id + domain_alibaba.FederatedAlibabaSecretAccessKeySuffix)
	if err != nil {
		return err
	}

	err = actions.Keychain.DeleteSecret(id + domain_alibaba.FederatedAlibabaStsTokenSuffix)
	if err != nil {
		return err
	}
	return nil
}

func (actions *AlibabaRamRoleFederatedSessionActions) Start(sessionId string) error {

	err := actions.AlibabaRamRoleFederatedSessionsFacade.StartingSession(sessionId)
	if err != nil {
		return err
	}

	err = actions.AlibabaRamRoleFederatedSessionsFacade.StartSession(sessionId)
	if err != nil {
		return err
	}

	return nil
}

func (actions *AlibabaRamRoleFederatedSessionActions) Stop(sessionId string) error {

	err := actions.AlibabaRamRoleFederatedSessionsFacade.StopSession(sessionId)
	if err != nil {
		return err
	}

	return nil
}
