package use_case

import (
	"encoding/json"
	"leapp_daemon/adapter/aws"
	"leapp_daemon/domain/domain_aws"
	"leapp_daemon/domain/domain_aws/aws_iam_user"
	"leapp_daemon/infrastructure/logging"
)

type AwsCredentialsApplier struct {
	Keychain            Keychain
	NamedProfilesFacade NamedProfilesFacade
	Repository          AwsConfigurationRepository
}

type AwsSessionToken struct {
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
	Expiration      string
}

func (applier *AwsCredentialsApplier) UpdateAwsIamUserSessions(oldSessions []aws_iam_user.AwsIamUserSession, newSessions []aws_iam_user.AwsIamUserSession) {

	activeCredentials := make([]aws.AwsTempCredentials, 0)
	for _, newSession := range newSessions {
		if newSession.Status != domain_aws.Active {
			continue
		}

		namedProfile, err := applier.NamedProfilesFacade.GetNamedProfileById(newSession.NamedProfileId)
		if err != nil {
			logging.Entry().Error(err)
			return
		}

		sessionTokenJson, err := applier.Keychain.GetSecret(newSession.SessionTokenLabel)
		if err != nil {
			logging.Entry().Error(err)
			return
		}

		sessionToken := AwsSessionToken{}
		err = json.Unmarshal([]byte(sessionTokenJson), &sessionToken)
		if err != nil {
			logging.Entry().Error(err)
			return
		}

		tempCredentials := aws.AwsTempCredentials{
			ProfileName:  namedProfile.Name,
			AccessKeyId:  sessionToken.AccessKeyId,
			SecretKey:    sessionToken.SecretAccessKey,
			SessionToken: sessionToken.SessionToken,
			Expiration:   sessionToken.Expiration,
			Region:       newSession.Region,
		}
		activeCredentials = append(activeCredentials, tempCredentials)
	}
	err := applier.Repository.WriteCredentials(activeCredentials)
	if err != nil {
		logging.Entry().Error(err)
	}
}
