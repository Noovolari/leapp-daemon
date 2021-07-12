package alibaba_ram_user

import "leapp_daemon/domain/domain_alibaba"

/*type PlainAlibabaSessionContainer interface {
	AddPlainAlibabaSession(PlainAlibabaSession) error
	GetAllPlainAlibabaSessions() ([]PlainAlibabaSession, error)
	RemovePlainAlibabaSession(session PlainAlibabaSession) error
}*/

type AlibabaRamUserSession struct {
	Id             string
	Name           string
	Region         string
	Status         domain_alibaba.AlibabaSessionStatus
	NamedProfileId string
}

func (sess *AlibabaRamUserSession) GetId() string {
	return sess.Id
}

func (sess *AlibabaRamUserSession) GetTypeString() string {
	return "user"
}
