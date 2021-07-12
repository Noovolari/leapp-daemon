package use_case

import (
	"leapp_daemon/domain/domain_aws/aws_iam_role_federated"
)

func CreateAwsIamRoleFederatedSession(name string, accountNumber string, roleName string, roleArn string,
	idpArn string, region string, ssoUrl string, profile string) error {

	/*
		config, err := configuration.ReadConfiguration()
		if err != nil { return err }

		err = session2.CreateAwsIamRoleFederatedSession(config, name, accountNumber, roleName, roleArn, idpArn, region, ssoUrl, profile)
		if err != nil { return err }

		err = config.Update()
		if err != nil { return err }
	*/

	return nil
}

func GetAwsIamRoleFederatedSession(id string) (*aws_iam_role_federated.AwsIamRoleFederatedSession, error) {
	var sess *aws_iam_role_federated.AwsIamRoleFederatedSession

	/*
		config, err := configuration.ReadConfiguration()
		if err != nil {
			return sess, err
		}

		sess, err = session2.GetAwsIamRoleFederatedSession(config, id)
		if err != nil {
			return sess, err
		}
	*/

	return sess, nil
}

func UpdateAwsIamRoleFederatedSession(sessionId string, name string, accountNumber string, roleName string, roleArn string,
	idpArn string, region string, ssoUrl string, profile string) error {

	/*
		config, err := configuration.ReadConfiguration()
		if err != nil {
			return err
		}

		err = session2.UpdateAwsIamRoleFederatedSession(config, sessionId, name, accountNumber, roleName, roleArn, idpArn, region, ssoUrl, profile)
		if err != nil {
			return err
		}

		err = config.Update()
		if err != nil {
			return err
		}
	*/

	return nil
}

func DeleteAwsIamRoleFederatedSession(sessionId string) error {
	/*
		config, err := configuration.ReadConfiguration()
		if err != nil {
			return err
		}

		err = session2.DeleteAwsIamRoleFederatedSession(config, sessionId)
		if err != nil {
			return err
		}

		err = config.Update()
		if err != nil {
			return err
		}
	*/

	return nil
}

func StartAwsIamRoleFederatedSession(sessionId string) error {
	/*
		config, err := configuration.ReadConfiguration()
		if err != nil {
			return err
		}

		// Passing nil because, it will be the rotate method to check if we need the mfaToken or not
		err = session2.StartAwsIamRoleFederatedSession(config, sessionId)
		if err != nil {
			return err
		}

		err = config.Update()
		if err != nil {
			return err
		}
	*/

	return nil
}

func StopAwsIamRoleFederatedSession(sessionId string) error {
	/*
		config, err := configuration.ReadConfiguration()
		if err != nil {
			return err
		}

		// Passing nil because, it will be the rotate method to check if we need the mfaToken or not
		err = session2.StopAwsIamRoleFederatedSession(config, sessionId)
		if err != nil {
			return err
		}

		err = config.Update()
		if err != nil {
			return err
		}

		// sess, err := session.GetAwsIamUserSession(config, sessionId)
		err = session_token.RemoveFromIniFile("default")
		if err != nil {
			return err
		}
	*/

	return nil
}
