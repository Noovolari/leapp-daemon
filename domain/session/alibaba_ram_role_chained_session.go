package session

/*type AlibabaRamRoleChainedSessionContainer interface {
	AddAlibabaRamRoleChainedSession(AlibabaRamRoleChainedSession) error
	GetAllAlibabaRamRoleChainedSessions() ([]AlibabaRamRoleChainedSession, error)
	RemoveAlibabaRamRoleChainedSession(session AlibabaRamRoleChainedSession) error
}*/

type AlibabaParentSession interface {
	GetId() string
	GetTypeString() string
}

type AlibabaRamRoleChainedSession struct {
	Id         string
	Status     Status
	StartTime  string
	ParentId   string
	ParentType string
	Account    *AlibabaRamRoleChainedAccount
	Profile    string
}

type AlibabaRamRoleChainedAccount struct {
	AccountNumber string
	Name          string
	Role          *AlibabaRamRole//AlibabaRamRoleChainedRole
	Region        string
	// Type            string
	// ParentSessionId string
	// ParentRole      string
	NamedProfileId string
}

/*type AlibabaRamRoleChainedRole struct {
	Name string
	Arn  string
	// Parent string
	// ParentRole string
}*/

/*
func CreateTrusterAlibabaSession(AccountName string, AccountNumber string, RoleName string, Region string) error {

  sessions, err := sessionContainer.GetPlainAlibabaSessions()
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

  plainAlibabaAccount := PlainAlibabaAccount{
    AccountNumber: accountNumber,
    Name:          name,
    Region:        region,
    User:          user,
    AlibabaAccessKeyId: alibabaAccessKeyId,
    AlibabaSecretAccessKey: alibabaSecretAccessKey,
    MfaDevice:     mfaDevice,

  }

  uuidString := uuid.New().String()
  uuidString = strings.Replace(uuidString, "-", "", -1)

  namedProfileId, err := CreateNamedProfile(sessionContainer, profile)
  if err != nil {
    return err
  }


  sess := PlainAlibabaSession{
    Id:        uuidString,
    Status:    NotActive,
    StartTime: "",
    Account:   &plainAlibabaAccount,
    Profile: namedProfileId,
  }

  err = sessionContainer.SetPlainAlibabaSessions(append(sessions, &sess))
  if err != nil {
    return err
  }

  return nil
}
*/