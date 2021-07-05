package use_case

import (
	"fmt"
	"leapp_daemon/domain/constant"
	"leapp_daemon/domain/region"

	"leapp_daemon/domain/session"
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
		account := sess.Account
		if sess.ParentId == parentId && account.AccountNumber == accountNumber && account.Role.Name == roleName {
			err := http_error.NewConflictError(fmt.Errorf("a session with the same parent, account number and role name already exists"))
			return err
		}
	}

	alibabaRamRoleChainedAccount := session.AlibabaRamRoleChainedAccount{
		AccountNumber: accountNumber,
		Name:          accountName,
		Role: &session.AlibabaRamRole{
			Name: roleName,
			Arn:  fmt.Sprintf("acs:ram::%s:role/%s", accountNumber, roleName),
		},
		Region:         region,
		NamedProfileId: namedProfile.Id,
	}

	sess := session.AlibabaRamRoleChainedSession{
		Id:         actions.Environment.GenerateUuid(),
		Status:     session.NotActive,
		StartTime:  "",
		ParentId:   parentSession.GetId(),
		ParentType: parentSession.GetTypeString(),
		Account:    &alibabaRamRoleChainedAccount,
	}

	actions.AlibabaRamRoleChainedSessionsFacade.SetSessions(append(sessions, sess))

	return nil
}

func (actions *AlibabaRamRoleChainedSessionActions) Get(id string) (*session.AlibabaRamRoleChainedSession, error) {
	return actions.AlibabaRamRoleChainedSessionsFacade.GetSessionById(id)
}

func (actions *AlibabaRamRoleChainedSessionActions) Update(id string, parentId string, accountName string, accountNumber string, roleName string, regionName string, profileName string) error {
	parentSession, err := GetAlibabaParentById(parentId)
	if err != nil {
		return err
	}

	isRegionValid := region.IsAlibabaRegionValid(regionName)
	if !isRegionValid {
		return http_error.NewUnprocessableEntityError(fmt.Errorf("Region " + regionName + " not valid"))
	}

	oldSess, err := actions.AlibabaRamRoleChainedSessionsFacade.GetSessionById(id)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	alibabaRamRoleChainedAccount := session.AlibabaRamRoleChainedAccount{
		AccountNumber: accountNumber,
		Name:          accountName,
		Role: &session.AlibabaRamRole{
			Name: roleName,
			Arn:  fmt.Sprintf("acs:ram::%s:role/%s", accountNumber, roleName),
		},
		Region:         regionName,
		NamedProfileId: oldSess.Account.NamedProfileId,
	}

	sess := session.AlibabaRamRoleChainedSession{
		Id:     id,
		Status: session.NotActive,
		//StartTime string
		ParentId:   parentId,
		ParentType: parentSession.GetTypeString(),
		Account:    &alibabaRamRoleChainedAccount,
		Profile:    profileName,
	}

	oldNamedProfile, err := actions.NamedProfilesActions.GetNamedProfileById(oldSess.Account.NamedProfileId)
	if err != nil {
		return err //TODO: return right error
	}
	oldNamedProfile.Name = profileName
	err = actions.NamedProfilesActions.SetNamedProfileName(oldNamedProfile)
	if err != nil {
		return err //TODO: return right error
	}
	err = actions.AlibabaRamRoleChainedSessionsFacade.SetSessionById(&sess)
	if err != nil {
		return err //TODO: return right error
	}
	return nil
}

func (actions *AlibabaRamRoleChainedSessionActions) Delete(id string) error {
	sess, err := actions.AlibabaRamRoleChainedSessionsFacade.GetSessionById(id)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	if sess.Status != session.NotActive {
		err = actions.Stop(id)
		if err != nil {
			return err
		}
	}

	oldSess, err := actions.AlibabaRamRoleChainedSessionsFacade.GetSessionById(id)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	oldNamedProfile, err := actions.NamedProfilesActions.GetNamedProfileById(oldSess.Account.NamedProfileId)
	if err != nil {
		return err //TODO: return right error
	}
	actions.NamedProfilesActions.DeleteNamedProfile(oldNamedProfile.Id)

	err = actions.AlibabaRamRoleChainedSessionsFacade.RemoveSession(id)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = actions.Keychain.DeleteSecret(id + constant.TrustedAlibabaKeyIdSuffix)
	if err != nil {
		return err
	}

	err = actions.Keychain.DeleteSecret(id + constant.TrustedAlibabaSecretAccessKeySuffix)
	if err != nil {
		return err
	}

	err = actions.Keychain.DeleteSecret(id + constant.TrustedAlibabaStsTokenSuffix)
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
	region := sess.Account.Region
	label := sess.ParentId + "-" + sess.ParentType + "-alibaba-session-access-key-id"
	accessKeyId, err := actions.Keychain.GetSecret(label)
	if err != nil {
		return err
	}
	label = sess.ParentId + "-" + sess.ParentType + "-alibaba-session-secret-access-key"
	accessKeySecret, err := actions.Keychain.GetSecret(label)
	if err != nil {
		return err
	}

	var client *sts.Client
	if sess.ParentType == "plain" {
		client, err = sts.NewClientWithAccessKey(region, accessKeyId, accessKeySecret)
		if err != nil {
			return err
		}
	} else {
		label = sess.ParentId + "-" + sess.ParentType + "-alibaba-session-sts-token"
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
	request.RoleArn = sess.Account.Role.Arn
	request.RoleSessionName = "leapp" // TODO: is this ok?
	response, err := client.AssumeRole(request)
	if err != nil {
		return err
	}

	// saves credentials into keychain
	err = actions.Keychain.SetSecret(response.Credentials.AccessKeyId, sess.Id+"-trusted-alibaba-session-access-key-id")
	if err != nil {
		return http_error.NewInternalServerError(err)
	}
	err = actions.Keychain.SetSecret(response.Credentials.AccessKeySecret, sess.Id+"-trusted-alibaba-session-secret-access-key")
	if err != nil {
		return http_error.NewInternalServerError(err)
	}
	err = actions.Keychain.SetSecret(response.Credentials.SecurityToken, sess.Id+"-trusted-alibaba-session-sts-token")
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	err = session.GetAlibabaRamRoleChainedSessionsFacade().SetSessionStatusToPending(sessionId)
	if err != nil {
		return err
	}

	err = session.GetAlibabaRamRoleChainedSessionsFacade().SetSessionStatusToActive(sessionId)
	if err != nil {
		return err
	}

	return nil
}

func (actions *AlibabaRamRoleChainedSessionActions) Stop(sessionId string) error {
	err := session.GetAlibabaRamRoleChainedSessionsFacade().SetSessionStatusToInactive(sessionId)
	if err != nil {
		return err
	}

	return nil
}

func GetAlibabaParentById(parentId string) (session.AlibabaParentSession, error) {
	plain, err := session.GetAlibabaRamUserSessionsFacade().GetSessionById(parentId)
	if err != nil {
		federated, err := session.GetAlibabaRamRoleFederatedSessionsFacade().GetSessionById(parentId)
		if err != nil {
			return nil, http_error.NewNotFoundError(fmt.Errorf("no plain or federated session with id %s found", parentId))
		}
		return federated, nil
	}
	return plain, nil
}
