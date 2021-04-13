package session

import (
  "fmt"
  "github.com/google/uuid"
  "leapp_daemon/domain/constant"
  http_error2 "leapp_daemon/infrastructure/http/http_error"
  "strings"
  "time"
)

type FederatedAwsSession struct {
	Id        string
	Status    Status
	StartTime string
	Account   *FederatedAwsAccount
	Profile   string
}

type FederatedAwsAccount struct {
	AccountNumber string
	Name          string
	Role          *FederatedAwsRole
	IdpArn        string
	Region        string
	SsoUrl        string
}

type FederatedAwsRole struct {
	Name string
	Arn  string
}

func(sess *FederatedAwsSession) Rotate(rotateConfiguration *RotateConfiguration) error {
	// TODO: implement rotate method for federated
	return nil
}

func(sess *FederatedAwsSession) IsRotationIntervalExpired() (bool, error) {
	startTime, _ := time.Parse(time.RFC3339, sess.StartTime)
	secondsPassedFromStart := time.Now().Sub(startTime).Seconds()
	return int64(secondsPassedFromStart) > constant.RotationIntervalInSeconds, nil
}

func CreateFederatedAwsSession(sessionContainer Container, name string, accountNumber string, roleName string, roleArn string, idpArn string,
	region string, ssoUrl string, profile string) error {

	sessions, err := sessionContainer.GetFederatedAwsSessions()
	if err != nil {
		return err
	}

	for _, session := range sessions {
		account := session.Account
		if account.AccountNumber == accountNumber && account.Role.Name == roleName {
			err = http_error2.NewUnprocessableEntityError(fmt.Errorf("an account with the same account number and " +
				"role name is already present"))
			return err
		}
	}

	role := FederatedAwsRole{
		Name: roleName,
		Arn:  roleArn,
	}

	federatedAwsAccount := FederatedAwsAccount{
		AccountNumber: accountNumber,
		Name:          name,
		Role:          &role,
		IdpArn:        idpArn,
		Region:        region,
		SsoUrl:        ssoUrl,
	}

	uuidString := uuid.New().String()
	uuidString = strings.Replace(uuidString, "-", "", -1)

	namedProfileId, err := CreateNamedProfile(sessionContainer, profile)
	if err != nil {
		return err
	}

	session := FederatedAwsSession{
		Id:        uuidString,
		Status:    NotActive,
		StartTime: "",
		Account:   &federatedAwsAccount,
		Profile:   namedProfileId,
	}

	err = sessionContainer.SetFederatedAwsSessions(append(sessions, &session))
	if err != nil { return err }

	return nil
}

func GetFederatedAwsSession(sessionContainer Container, id string) (*FederatedAwsSession, error) {
	sessions, err := sessionContainer.GetFederatedAwsSessions()
	if err != nil {
		return nil, err
	}

	for index, _ := range sessions {
		if sessions[index].Id == id {
			return sessions[index], nil
		}
	}

	return nil, http_error2.NewNotFoundError(fmt.Errorf("No session found with id:" + id))
}

func ListFederatedAwsSession(sessionContainer Container, query string) ([]*FederatedAwsSession, error) {
	sessions, err := sessionContainer.GetFederatedAwsSessions()
	if err != nil {
		return nil, err
	}

	filteredList := make([]*FederatedAwsSession, 0)

	if query == "" {
		return append(filteredList, sessions...), nil
	} else {
		for _, session := range sessions {
			if  strings.Contains(session.Id, query) ||
				strings.Contains(session.Profile, query) ||
				strings.Contains(session.Account.Name, query) ||
				strings.Contains(session.Account.IdpArn, query) ||
				strings.Contains(session.Account.SsoUrl, query) ||
				strings.Contains(session.Account.Region, query) ||
				strings.Contains(session.Account.AccountNumber, query) ||
				strings.Contains(session.Account.Role.Name, query) ||
				strings.Contains(session.Account.Role.Arn, query) {

				filteredList = append(filteredList, session)
			}
		}

		return filteredList, nil
	}
}

func UpdateFederatedAwsSession(sessionContainer Container, id string, name string, accountNumber string, roleName string, roleArn string, idpArn string,
	region string, ssoUrl string, profile string) error {

	sessions, err := sessionContainer.GetFederatedAwsSessions()
	if err != nil { return err }

	found := false
	for index := range sessions {
		if sessions[index].Id == id {
			namedProfileId, err := EditNamedProfile(sessionContainer, sessions[index].Profile, profile)
			if err != nil { return err }

			sessions[index].Profile = namedProfileId
			sessions[index].Account = &FederatedAwsAccount{
				AccountNumber: accountNumber,
				Name:          name,
				Region:        region,
				IdpArn: 	   idpArn,
				SsoUrl:        ssoUrl,
			}

			sessions[index].Account.Role = &FederatedAwsRole{
				Name: roleName,
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
	sessions, err := sessionContainer.GetFederatedAwsSessions()
	if err != nil {
		return err
	}

	found := false
	for index := range sessions {
		if sessions[index].Id == id {
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

	println("Rotating session with id", sess.Id)
	err = sess.Rotate(nil)
	if err != nil { return err }

	return nil
}

func StopFederatedAwsSession(sessionContainer Container, id string) error {
	sess, err := GetFederatedAwsSession(sessionContainer, id)
	if err != nil {
		return err
	}

	sess.Status = NotActive
	return nil
}
