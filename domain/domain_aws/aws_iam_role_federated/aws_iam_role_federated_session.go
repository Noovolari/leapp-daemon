package aws_iam_role_federated

import (
	"leapp_daemon/domain/domain_aws"
	"time"
)

type AwsIamRoleFederatedSession struct {
	Id        string
	Status    domain_aws.AwsSessionStatus
	StartTime string
	Account   *AwsIamRoleFederatedAccount
	Profile   string
}

type AwsIamRoleFederatedAccount struct {
	AccountNumber string
	Name          string
	Role          *AwsIamRoleFederatedRole
	IdpArn        string
	Region        string
	SsoUrl        string
}

type AwsIamRoleFederatedRole struct {
	Name string
	Arn  string
}

func (sess *AwsIamRoleFederatedSession) IsRotationIntervalExpired() (bool, error) {
	startTime, _ := time.Parse(time.RFC3339, sess.StartTime)
	secondsPassedFromStart := time.Now().Sub(startTime).Seconds()
	return int64(secondsPassedFromStart) > domain_aws.RotationIntervalInSeconds, nil
}

/*
func CreateAwsIamRoleFederatedSession(sessionContainer Container, name string, accountNumber string, roleName string, roleArn string, idpArn string,
	region string, ssoUrl string, profile string) error {

	sessions, err := sessionContainer.AwsGetIamRoleFederatedSessions()
	if err != nil {
		return err
	}

	for _, session := range sessions {
		account := session.Account
		if account.AccountNumber == accountNumber && account.Role.SessionName == roleName {
			err = http_error2.NewUnprocessableEntityError(fmt.Errorf("an account with the same account number and " +
				"role name is already present"))
			return err
		}
	}

	role := AwsIamRoleFederatedRole{
		SessionName: roleName,
		Arn:  roleArn,
	}

	federatedAwsAccount := AwsIamRoleFederatedAccount{
		AccountNumber: accountNumber,
		SessionName:          name,
		Role:          &role,
		IdpArn:        idpArn,
		Region:        region,
		SsoUrl:        ssoUrl,
	}

	uuidString := uuid.New().String() //use Environment.GenerateUuid()
	uuidString = strings.Replace(uuidString, "-", "", -1)

	namedProfileId, err := named_profile.CreateNamedProfile(sessionContainer, profile)
	if err != nil {
		return err
	}

	session := AwsIamRoleFederatedSession{
		ID:        uuidString,
		AwsSessionStatus:    NotActive,
		StartTime: "",
		Account:   &federatedAwsAccount,
		Profile:   namedProfileId,
	}

	err = sessionContainer.SetFederatedAwsSessions(append(sessions, &session))
	if err != nil { return err }

	return nil
}

func GetFederatedAwsSession(sessionContainer Container, id string) (*AwsIamRoleFederatedSession, error) {
	sessions, err := sessionContainer.AwsGetIamRoleFederatedSessions()
	if err != nil {
		return nil, err
	}

	for index, _ := range sessions {
		if sessions[index].ID == id {
			return sessions[index], nil
		}
	}

	return nil, http_error2.NewNotFoundError(fmt.Errorf("No session found with id:" + id))
}

func ListFederatedAwsSession(sessionContainer Container, query string) ([]*AwsIamRoleFederatedSession, error) {
	sessions, err := sessionContainer.AwsGetIamRoleFederatedSessions()
	if err != nil {
		return nil, err
	}

	filteredList := make([]*AwsIamRoleFederatedSession, 0)

	if query == "" {
		return append(filteredList, sessions...), nil
	} else {
		for _, session := range sessions {
			if  strings.Contains(session.ID, query) ||
				strings.Contains(session.Profile, query) ||
				strings.Contains(session.Account.SessionName, query) ||
				strings.Contains(session.Account.IdpArn, query) ||
				strings.Contains(session.Account.SsoUrl, query) ||
				strings.Contains(session.Account.Region, query) ||
				strings.Contains(session.Account.AccountNumber, query) ||
				strings.Contains(session.Account.Role.SessionName, query) ||
				strings.Contains(session.Account.Role.Arn, query) {

				filteredList = append(filteredList, session)
			}
		}

		return filteredList, nil
	}
}

func UpdateFederatedAwsSession(sessionContainer Container, id string, name string, accountNumber string, roleName string, roleArn string, idpArn string,
	region string, ssoUrl string, profile string) error {

	sessions, err := sessionContainer.AwsGetIamRoleFederatedSessions()
	if err != nil { return err }

	found := false
	for index := range sessions {
		if sessions[index].ID == id {
			namedProfileId, err := named_profile.EditNamedProfile(sessionContainer, sessions[index].Profile, profile)
			if err != nil { return err }

			sessions[index].Profile = namedProfileId
			sessions[index].Account = &AwsIamRoleFederatedAccount{
				AccountNumber: accountNumber,
				SessionName:          name,
				Region:        region,
				IdpArn: 	   idpArn,
				SsoUrl:        ssoUrl,
			}

			sessions[index].Account.Role = &AwsIamRoleFederatedRole{
				SessionName: roleName,
				Arn:  roleArn,
			}

			found = true
		}
	}

	if found == false {
		err = http_error2.NewNotFoundError(fmt.Errorf("federated AWS session with id " + id + " not found"))
		return err
	}

	err = sessionContainer.SetFederatedAwsSessions(sessions)
	if err != nil { return err }

	return nil
}

func DeleteFederatedAwsSession(sessionContainer Container, id string) error {
	sessions, err := sessionContainer.AwsGetIamRoleFederatedSessions()
	if err != nil {
		return err
	}

	found := false
	for index := range sessions {
		if sessions[index].ID == id {
			sessions = append(sessions[:index], sessions[index+1:]...)
			found = true
			break
		}
	}

	if found == false {
		err = http_error2.NewNotFoundError(fmt.Errorf("federated AWS session with id " + id + " not found"))
		return err
	}

	err = sessionContainer.SetFederatedAwsSessions(sessions)
	if err != nil {
		return err
	}

	return nil
}

func StartFederatedAwsSession(sessionContainer Container, id string) error {
	sess, err := GetFederatedAwsSession(sessionContainer, id)
	if err != nil {
		return err
	}

	println("Rotating session with id", sess.ID)
	err = sess.Rotate(nil)
	if err != nil { return err }

	return nil
}

func StopFederatedAwsSession(sessionContainer Container, id string) error {
	sess, err := GetFederatedAwsSession(sessionContainer, id)
	if err != nil {
		return err
	}

	sess.AwsSessionStatus = NotActive
	return nil
}
*/
