package session

import (
	"leapp_daemon/domain/constant"
	"time"
)

type AwsIamUserSessionContainer interface {
	AddAwsIamUserSession(AwsIamUserSession) error
	GetAllAwsIamUserSessions() ([]AwsIamUserSession, error)
	RemoveAwsIamUserSession(session AwsIamUserSession) error
}

type AwsIamUserSession struct {
	Id           string
	Alias        string
	Status       Status
	StartTime    string
	LastStopTime string
	Account      *AwsIamUserAccount
}

type AwsIamUserAccount struct {
	MfaDevice              string
	Region                 string
	NamedProfileId         string
	SessionTokenExpiration string
}

type AwsSessionToken struct {
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
}

func (sess *AwsIamUserSession) IsMfaRequired() (bool, error) {
	return sess.Account.MfaDevice != "", nil
}

func (sess *AwsIamUserSession) IsRotationIntervalExpired() (bool, error) {
	startTime, _ := time.Parse(time.RFC3339, sess.StartTime)
	secondsPassedFromStart := time.Now().Sub(startTime).Seconds()
	return int64(secondsPassedFromStart) > constant.RotationIntervalInSeconds, nil
}

/*
func(sess *AwsIamUserSession) Rotate(rotateConfiguration *RotateConfiguration) error {
	if sess.Status == Active {
		isRotationIntervalExpired, err := sess.IsRotationIntervalExpired()
		if err != nil {
			return err
		}

		if isRotationIntervalExpired {

			err = sess.RotateAwsIamUserSessionCredentials(&rotateConfiguration.MfaToken)
			if err != nil {
				return nil
			}
		}
	}

	return nil
}

func(sess *AwsIamUserSession) RotateAwsIamUserSessionCredentials(mfaToken *string) error {
	doSessionTokenExist, err := session_token.DoExist(sess.Account.Name)
	if err != nil {
		return err
	}

	if doSessionTokenExist {
		isSessionTokenExpired, err := session_token.IsExpired(sess.Account.Name)
		if err != nil {
			return err
		}

		if isSessionTokenExpired {
			logging.Entry().Error("AWS Iam User session token no more valid")

			isMfaTokenRequired, err := sess.IsMfaRequired()
			if err != nil {
				return err
			}

			if isMfaTokenRequired && mfaToken == nil {
				sess.Status = Pending

				err = sendMfaRequestMessage(sess)
				if err != nil {
					return err
				}

				return nil
			}

			credentials, err := session_token.Generate(sess.Account.Name, sess.Account.Region, sess.Account.MfaDevice, mfaToken)
			if err != nil {
				return err
			}

			err = session_token.SaveInKeychain(sess.Account.Name, credentials)
			if err != nil {
				return err
			}

			err = session_token.SaveInIniFile(*credentials.AccessKeyId, *credentials.SecretAccessKey,
				*credentials.SessionToken, sess.Account.Region, "default")
			if err != nil {
				return err
			}

			sess.Status = Active
			sess.StartTime = time.Now().Format(time.RFC3339)
		} else {
			logging.Entry().Error("AWS Iam User session token still valid")

			data, err := sess.unmarshallSessionToken()
			if err != nil {
				return err
			}

			err = session_token.SaveInIniFile(data.AccessKeyId,
				data.SecretAccessKey,
				data.SessionToken,
				sess.Account.Region,
				"default")
			if err != nil {
				return err
			}

			sess.Status = Active
			sess.StartTime = time.Now().Format(time.RFC3339)
		}

		return nil
	} else {
		credentials, err := session_token.Generate(sess.Account.Name, sess.Account.Region, sess.Account.MfaDevice, mfaToken)
		if err != nil {
			return err
		}

		err = session_token.SaveInKeychain(sess.Account.Name, credentials)
		if err != nil {
			return err
		}

		err = session_token.SaveInIniFile(*credentials.AccessKeyId, *credentials.SecretAccessKey,
			*credentials.SessionToken, sess.Account.Region, "default")
		if err != nil {
			return err
		}

		sess.Status = Active
		sess.StartTime = time.Now().Format(time.RFC3339)

		return nil
	}
}

func (sess *AwsIamUserSession) unmarshallSessionToken() (AwsSessionToken, error) {
	sessionTokenJson, _, err := access_keys.Get(sess.Account.Name)

	var data AwsSessionToken

	err = json.Unmarshal([]byte(sessionTokenJson), &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func getById(sessionContainer Container, id string) (*AwsIamUserSession, error) {
	var sess *AwsIamUserSession

	sessions, err := sessionContainer.GetSessions()
	if err != nil { return sess, err }

	for index := range sessions {
		if sessions[index].Id == id {
			sess = sessions[index]
			return sess, nil
		}
	}

	err = http_error.NewNotFoundError(fmt.Errorf("AWS Iam User session with id " + id + " not found"))
	return sess, err
}

func sendMfaRequestMessage(sess *AwsIamUserSession) error {
	messageData := websocket.MfaTokenRequestData{
		SessionId: sess.Id,
	}

	messageDataJson, err := json.Marshal(messageData)
	if err != nil {
		return http_error.NewUnprocessableEntityError(err)
	}

	msg := websocket.Message{
		MessageType: websocket.MfaTokenRequest,
		Data:        string(messageDataJson),
	}

	err = websocket.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}

func CreateAwsIamUserSession(sessionContainer Container, name string, accountNumber string,
	region string, user string, awsAccessKeyId string, awsSecretAccessKey string, mfaDevice string, profile string) error {

	sessions, err := sessionContainer.GetSessions()
	if err != nil {
		return err
	}

	for _, sess := range sessions {
		account := sess.Account
		if account.AccountNumber == accountNumber && account.User == user {
			err := http_error.NewUnprocessableEntityError(fmt.Errorf("a session with the same account number and user is already present"))
			return err
		}
	}

	awsIamUserAccount := AwsIamUserAccount{
		AccountNumber: accountNumber,
		Name:          name,
		Region:        region,
		User:          user,
		AwsAccessKeyId: awsAccessKeyId,
		AwsSecretAccessKey: awsSecretAccessKey,
		MfaDevice:     mfaDevice,
	}

	uuidString := uuid.New().String() //use Environment.GenerateUuid()
	uuidString = strings.Replace(uuidString, "-", "", -1)

	namedProfileId, err := named_profile.CreateNamedProfile(sessionContainer, profile)
	if err != nil {
		return err
	}

	sess := AwsIamUserSession{
		Id:        uuidString,
		Status:    NotActive,
		StartTime: "",
		Account:   &awsIamUserAccount,
		Profile:   namedProfileId,
	}

	err = sessionContainer.SetSessions(append(sessions, &sess))
	if err != nil {
		return err
	}

	return nil
}

func GetAwsIamUserSession(sessionContainer Container, id string) (*AwsIamUserSession, error) {
	sess, err := getById(sessionContainer, id)
	if err != nil {
		return sess, err
	}

	return sess, nil
}

func ListAwsIamUserSession(sessionContainer Container, query string) ([]*AwsIamUserSession, error) {
	filteredList := make([]*AwsIamUserSession, 0)

	allSessions, err := sessionContainer.GetSessions()
	if err != nil { return filteredList, nil }

	if query == "" {
		return append(filteredList, allSessions...), nil
	} else {
		for _, sess := range allSessions {
			if  strings.Contains(sess.Id, query) ||
				strings.Contains(sess.Profile, query) ||
				strings.Contains(sess.Account.Name, query) ||
				strings.Contains(sess.Account.MfaDevice, query) ||
				strings.Contains(sess.Account.User, query) ||
				strings.Contains(sess.Account.Region, query) ||
				strings.Contains(sess.Account.AccountNumber, query) {

				filteredList = append(filteredList, sess)
			}
		}

		return filteredList, nil
	}
}

func UpdateAwsIamUserSession(sessionContainer Container, id string, name string, accountNumber string, region string,
	user string, awsAccessKeyId string, awsSecretAccessKey string, mfaDevice string, profile string) error {

	sessions, err := sessionContainer.GetSessions()
	if err != nil { return err }

	found := false
	for index := range sessions {
		if sessions[index].Id == id {
			namedProfileId, err := named_profile.EditNamedProfile(sessionContainer, sessions[index].Profile, profile)
			if err != nil { return err }

			sessions[index].Profile = namedProfileId
			sessions[index].Account = &AwsIamUserAccount{
				AccountNumber: accountNumber,
				Name:          name,
				Region:        region,
				User:          user,
				AwsAccessKeyId: awsAccessKeyId,
				AwsSecretAccessKey: awsSecretAccessKey,
				MfaDevice:     mfaDevice,
			}
			found = true
		}
	}

	if found == false {
		err = http_error.NewNotFoundError(fmt.Errorf("AWS Iam User session with id " + id + " not found"))
		return err
	}

	err = sessionContainer.SetSessions(sessions)
	if err != nil { return err }

	return nil
}

func DeleteAwsIamUserSession(sessionContainer Container, id string) error {
	sessions, err := sessionContainer.GetSessions()
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
		err = http_error.NewNotFoundError(fmt.Errorf("AWS Iam User session with id " + id + " not found"))
		return err
	}

	err = sessionContainer.SetSessions(sessions)
	if err != nil {
		return err
	}

	return nil
}

func IsMfaRequiredForAwsIamUserSession(sessionContainer Container, id string) (bool, error) {
	sess, err := getById(sessionContainer, id)

	if err != nil {
		return false, err
	}

	return sess.IsMfaRequired()
}

func StartAwsIamUserSession(sessionContainer Container, id string, mfaToken *string) error {
	sess, err := getById(sessionContainer, id)
	if err != nil {
		return err
	}

	err = sess.RotateAwsIamUserSessionCredentials(mfaToken)
	if err != nil { return err }

	return nil
}

func StopAwsIamUserSession(sessionContainer Container, id string) error {
	sess, err := GetAwsIamUserSession(sessionContainer, id)
	if err != nil {
		return err
	}

	sess.Status = NotActive
	return nil
}
*/
