package session

/*type PlainAlibabaSessionContainer interface {
	AddPlainAlibabaSession(PlainAlibabaSession) error
	GetAllPlainAlibabaSessions() ([]PlainAlibabaSession, error)
	RemovePlainAlibabaSession(session PlainAlibabaSession) error
}*/

type AlibabaRamUserSession struct {
	Id      string
	Alias   string
	Status  Status
	Account *AlibabaRamUserAccount
}

type AlibabaRamUserAccount struct {
	Region         string
	NamedProfileId string
}

func (sess *AlibabaRamUserSession) GetId() string {
	return sess.Id
}

func (sess *AlibabaRamUserSession) GetTypeString() string {
	return "plain"
}