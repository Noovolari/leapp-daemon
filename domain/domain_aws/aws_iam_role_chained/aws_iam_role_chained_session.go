package aws_iam_role_chained

import (
	"leapp_daemon/domain/domain_aws"
)

type AwsIamRoleChainedSession struct {
	Id        string
	Status    domain_aws.AwsSessionStatus
	StartTime string
	ParentId  string
	Account   *AwsIamRoleChainedAccount
}

type AwsIamRoleChainedAccount struct {
	AccountNumber string
	Name          string
	Role          *AwsIamRole
	Region        string
	// Type            string
	// ParentSessionId string
	// ParentRole      string
}

type AwsIamRole struct {
	Name string
	Arn  string
	// Parent string
	// ParentRole string
}

/*
func CreateTrusterAwsSession(AccountName string, AccountNumber string, RoleName string, Region string) error {

  sessions, err := sessionContainer.GetSessions()
  if err != nil {
    return err
  }

  for _, sess := range sessions {
    account := sess.Account
    if account.AccountNumber == accountNumber && account.UserName == user {
      err := http_error.NewUnprocessableEntityError(fmt.Errorf("a session with the same account number and user is already present"))
      return err
    }
  }

  awsIamUserAccount := AwsIamUserAccount{
    AccountNumber: accountNumber,
    SessionName:          name,
    Region:        region,
    UserName:          user,
    AwsAccessKeyId: awsAccessKeyId,
    AwsSecretKey: awsSecretAccessKey,
    MfaDevice:     mfaDevice,

  }

  uuidString := uuid.New().String() //use Environment.GenerateUuid()
  uuidString = strings.Replace(uuidString, "-", "", -1)

  namedProfileId, err := CreateNamedProfile(sessionContainer, profile)
  if err != nil {
    return err
  }


  sess := AwsIamUserSession{
    Id:        uuidString,
    AwsSessionStatus:    NotActive,
    StartTime: "",
    Account:   &awsIamUserAccount,
    Profile: namedProfileId,
  }

  err = sessionContainer.SetSessions(append(sessions, &sess))
  if err != nil {
    return err
  }

  return nil
}
*/
