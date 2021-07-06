package session

import (
	"leapp_daemon/domain/constant"
	"time"
)

/*type AlibabaRamRoleFederatedSessionContainer interface {
	AddAlibabaRamRoleFederatedSession(AlibabaRamRoleFederatedSession) error
	GetAllAlibabaRamRoleFederatedSessions() ([]AlibabaRamRoleFederatedSession, error)
	RemoveAlibabaRamRoleFederatedSession(session AlibabaRamRoleFederatedSession) error
}*/

type AlibabaRamRoleFederatedSession struct {
	Id        string
	Status    Status
	StartTime string
	Account   *AlibabaRamRoleFederatedAccount
	Profile   string
}

type AlibabaRamRoleFederatedAccount struct {
	AccountNumber string
	Name          string
	Role          *AlibabaRamRole
	IdpArn        string
	Region        string
	/*SsoUrl        string*/
	NamedProfileId string
}

type AlibabaRamRole struct {
	Name string
	Arn  string
}

func (sess *AlibabaRamRoleFederatedSession) IsRotationIntervalExpired() (bool, error) {
	startTime, _ := time.Parse(time.RFC3339, sess.StartTime)
	secondsPassedFromStart := time.Now().Sub(startTime).Seconds()
	return int64(secondsPassedFromStart) > constant.RotationIntervalInSeconds, nil
}

func (sess *AlibabaRamRoleFederatedSession) GetId() string {
	return sess.Id
}

func (sess *AlibabaRamRoleFederatedSession) GetTypeString() string {
	return "role-federated"
}

/*
func CreateAlibabaRamRoleFederatedSession(sessionContainer Container, name string, accountNumber string, roleName string, roleArn string, idpArn string,
	region string, ssoUrl string, profile string) error {

	sessions, err := sessionContainer.GetAlibabaRamRoleFederatedSessions()
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

	role := AlibabaRamRoleFederatedRole{
		Name: roleName,
		Arn:  roleArn,
	}

	federatedAlibabaAccount := AlibabaRamRoleFederatedAccount{
		AccountNumber: accountNumber,
		Name:          name,
		Role:          &role,
		IdpArn:        idpArn,
		Region:        region,
		SsoUrl:        ssoUrl,
	}

	uuidString := uuid.New().String()
	uuidString = strings.Replace(uuidString, "-", "", -1)

	namedProfileId, err := named_profile.CreateNamedProfile(sessionContainer, profile)
	if err != nil {
		return err
	}

	session := AlibabaRamRoleFederatedSession{
		Id:        uuidString,
		Status:    NotActive,
		StartTime: "",
		Account:   &federatedAlibabaAccount,
		Profile:   namedProfileId,
	}

	err = sessionContainer.SetAlibabaRamRoleFederatedSessions(append(sessions, &session))
	if err != nil { return err }

	return nil
}

func GetAlibabaRamRoleFederatedSession(sessionContainer Container, id string) (*AlibabaRamRoleFederatedSession, error) {
	sessions, err := sessionContainer.GetAlibabaRamRoleFederatedSessions()
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

func ListAlibabaRamRoleFederatedSession(sessionContainer Container, query string) ([]*AlibabaRamRoleFederatedSession, error) {
	sessions, err := sessionContainer.GetAlibabaRamRoleFederatedSessions()
	if err != nil {
		return nil, err
	}

	filteredList := make([]*AlibabaRamRoleFederatedSession, 0)

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

func UpdateAlibabaRamRoleFederatedSession(sessionContainer Container, id string, name string, accountNumber string, roleName string, roleArn string, idpArn string,
	region string, ssoUrl string, profile string) error {

	sessions, err := sessionContainer.GetAlibabaRamRoleFederatedSessions()
	if err != nil { return err }

	found := false
	for index := range sessions {
		if sessions[index].Id == id {
			namedProfileId, err := named_profile.EditNamedProfile(sessionContainer, sessions[index].Profile, profile)
			if err != nil { return err }

			sessions[index].Profile = namedProfileId
			sessions[index].Account = &AlibabaRamRoleFederatedAccount{
				AccountNumber: accountNumber,
				Name:          name,
				Region:        region,
				IdpArn: 	   idpArn,
				SsoUrl:        ssoUrl,
			}

			sessions[index].Account.Role = &AlibabaRamRoleFederatedRole{
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

	err = sessionContainer.SetAlibabaRamRoleFederatedSessions(sessions)
	if err != nil { return err }

	return nil
}

func DeleteAlibabaRamRoleFederatedSession(sessionContainer Container, id string) error {
	sessions, err := sessionContainer.GetAlibabaRamRoleFederatedSessions()
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

	err = sessionContainer.SetAlibabaRamRoleFederatedSessions(sessions)
	if err != nil {
		return err
	}

	return nil
}

func StartAlibabaRamRoleFederatedSession(sessionContainer Container, id string) error {
	sess, err := GetAlibabaRamRoleFederatedSession(sessionContainer, id)
	if err != nil {
		return err
	}

	println("Rotating session with id", sess.Id)
	err = sess.Rotate(nil)
	if err != nil { return err }

	return nil
}

func StopAlibabaRamRoleFederatedSession(sessionContainer Container, id string) error {
	sess, err := GetAlibabaRamRoleFederatedSession(sessionContainer, id)
	if err != nil {
		return err
	}

	sess.Status = NotActive
	return nil
}
*/
