package use_case

import (
	"fmt"
	"leapp_daemon/domain/domain_alibaba"
	"leapp_daemon/domain/domain_alibaba/alibaba_ram_role_chained"
	"leapp_daemon/domain/domain_alibaba/alibaba_ram_role_federated"
	"leapp_daemon/domain/domain_alibaba/alibaba_ram_user"
	"leapp_daemon/infrastructure/http/http_error"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
)

type AlibabaRamRoleChainedSessionActions struct {
	Environment                         Environment
	Keychain                            Keychain
	NamedProfilesActions                NamedProfilesActionsInterface
	AlibabaRamRoleChainedSessionsFacade AlibabaRamRoleChainedSessionsFacade
}

func (actions *AlibabaRamRoleChainedSessionActions) Create(parentId string, accountName string, accountNumber string, roleName string, region string, profileName string) error {

	namedProfile, err := actions.NamedProfilesActions.GetOrCreateNamedProfile(profileName)
	if err != nil {
		return err
	}

	parentSession, err := GetAlibabaParentById(parentId)
	if err != nil {
		return err
	}

	sessions := actions.AlibabaRamRoleChainedSessionsFacade.GetSessions()

	for _, sess := range sessions {
		if sess.ParentId == parentId && sess.AccountNumber == accountNumber && sess.RoleName == roleName {
			err := http_error.NewConflictError(fmt.Errorf("a session with the same parent, account number and role name already exists"))
			return err
		}
	}

	sess := alibaba_ram_role_chained.AlibabaRamRoleChainedSession{
		Id:         actions.Environment.GenerateUuid(),
		Status:     domain_alibaba.NotActive,
		StartTime:  "",
		ParentId:   parentSession.GetId(),
		ParentType: parentSession.GetTypeString(),
		AccountNumber: accountNumber,
		Name:          accountName,
		RoleName: roleName,
		RoleArn:  fmt.Sprintf("acs:ram::%s:role/%s", accountNumber, roleName),
		Region:         region,
		NamedProfileId: namedProfile.Id,
	}

	actions.AlibabaRamRoleChainedSessionsFacade.SetSessions(append(sessions, sess))

	return nil
}

func (actions *AlibabaRamRoleChainedSessionActions) Get(id string) (*alibaba_ram_role_chained.AlibabaRamRoleChainedSession, error) {
	return actions.AlibabaRamRoleChainedSessionsFacade.GetSessionById(id)
}

func (actions *AlibabaRamRoleChainedSessionActions) Update(id string, parentId string, accountName string, accountNumber string, roleName string, regionName string, profileName string) error {
	parentSession, err := GetAlibabaParentById(parentId)
	if err != nil {
		return err
	}

	isRegionValid := domain_alibaba.IsAlibabaRegionValid(regionName)
	if !isRegionValid {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("Region " + regionName + " not valid"))
	}

	np, err := actions.NamedProfilesActions.GetOrCreateNamedProfile(profileName)
	if err != nil {
		return err //TODO: return right error
	}

	sess := alibaba_ram_role_chained.AlibabaRamRoleChainedSession{
		Id:     id,
		Status: domain_alibaba.NotActive,
		//StartTime string
		ParentId:       parentId,
		ParentType:     parentSession.GetTypeString(),
		AccountNumber:  accountNumber,
		Name:           accountName,
		RoleName:       roleName,
		RoleArn:        fmt.Sprintf("acs:ram::%s:role/%s", accountNumber, roleName),
		Region:         regionName,
		NamedProfileId: np.Id,
		Profile:        profileName,
	}

	actions.AlibabaRamRoleChainedSessionsFacade.SetSessionById(&sess)
	return nil
}

func (actions *AlibabaRamRoleChainedSessionActions) Delete(id string) error {
	sess, err := actions.AlibabaRamRoleChainedSessionsFacade.GetSessionById(id)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	if sess.Status != domain_alibaba.NotActive {
		err = actions.Stop(id)
		if err != nil {
			return err
		}
	}

	err = actions.AlibabaRamRoleChainedSessionsFacade.RemoveSession(id)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = actions.Keychain.DeleteSecret(id + domain_alibaba.TrustedAlibabaKeyIdSuffix)
	if err != nil {
		return err
	}

	err = actions.Keychain.DeleteSecret(id + domain_alibaba.TrustedAlibabaSecretAccessKeySuffix)
	if err != nil {
		return err
	}

	err = actions.Keychain.DeleteSecret(id + domain_alibaba.TrustedAlibabaStsTokenSuffix)
	if err != nil {
		return err
	}

	return nil
}

func (actions *AlibabaRamRoleChainedSessionActions) Start(sessionId string) error {
	// call AssumeRole API
	sess, err := actions.AlibabaRamRoleChainedSessionsFacade.GetSessionById(sessionId)
	if err != nil {
		return err
	}
	region := sess.Region
	label := sess.ParentId + "-alibaba-ram-" + sess.ParentType + "-session-access-key-id"
	accessKeyId, err := actions.Keychain.GetSecret(label)
	if err != nil {
		return err
	}
	label = sess.ParentId + "-alibaba-ram-" + sess.ParentType + "-session-secret-access-key"
	accessKeySecret, err := actions.Keychain.GetSecret(label)
	if err != nil {
		return err
	}

	var client *sts.Client
	if sess.ParentType == "user" {
		client, err = sts.NewClientWithAccessKey(region, accessKeyId, accessKeySecret)
		if err != nil {
			return err
		}
	} else {
		label = sess.ParentId + "-alibaba-ram-" + sess.ParentType + "-session-sts-token"
		stsToken, err := actions.Keychain.GetSecret(label)
		if err != nil {
			return err
		}
		client, err = sts.NewClientWithStsToken(region, accessKeyId, accessKeySecret, stsToken)
		if err != nil {
			return err
		}
	}

	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"
	request.RoleArn = sess.RoleArn
	request.RoleSessionName = "leapp" // TODO: is this ok?
	response, err := client.AssumeRole(request)
	if err != nil {
		return err
	}

	// saves credentials into keychain
	err = actions.Keychain.SetSecret(response.Credentials.AccessKeyId, sess.Id+domain_alibaba.TrustedAlibabaKeyIdSuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}
	err = actions.Keychain.SetSecret(response.Credentials.AccessKeySecret, sess.Id+domain_alibaba.TrustedAlibabaSecretAccessKeySuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}
	err = actions.Keychain.SetSecret(response.Credentials.SecurityToken, sess.Id+domain_alibaba.TrustedAlibabaStsTokenSuffix)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = alibaba_ram_role_chained.GetAlibabaRamRoleChainedSessionsFacade().StartingSession(sessionId)
	if err != nil {
		return err
	}

	err = alibaba_ram_role_chained.GetAlibabaRamRoleChainedSessionsFacade().StartSession(sessionId)
	if err != nil {
		return err
	}

	return nil
}

func (actions *AlibabaRamRoleChainedSessionActions) Stop(sessionId string) error {
	err := alibaba_ram_role_chained.GetAlibabaRamRoleChainedSessionsFacade().StopSession(sessionId)
	if err != nil {
		return err
	}

	return nil
}

func GetAlibabaParentById(parentId string) (alibaba_ram_role_chained.AlibabaParentSession, error) {
	plain, err := alibaba_ram_user.GetAlibabaRamUserSessionsFacade().GetSessionById(parentId)
	if err != nil {
		federated, err := alibaba_ram_role_federated.GetAlibabaRamRoleFederatedSessionsFacade().GetSessionById(parentId)
		if err != nil {
			return nil, http_error.NewNotFoundError(fmt.Errorf("no user or role session with id %s found", parentId))
		}
		return federated, nil
	}
	return plain, nil
}
