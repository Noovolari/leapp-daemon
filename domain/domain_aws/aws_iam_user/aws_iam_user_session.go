package aws_iam_user

import (
	"leapp_daemon/domain/domain_aws"
)

type AwsIamUserSessionContainer interface {
	AddAwsIamUserSession(AwsIamUserSession) error
	GetAllAwsIamUserSessions() ([]AwsIamUserSession, error)
	RemoveAwsIamUserSession(session AwsIamUserSession) error
}

// swagger:model AwsIamUserSession
type AwsIamUserSession struct {
	Id                     string  `json:"id"`
	Name                   string  `json:"name"`
	Region                 string  `json:"region"`
	AccountNumber          string  `json:"accountNumber"`
	UserName               string  `json:"userName"`
	AccessKeyIdLabel       string  `json:"accessKeyIdLabel"`
	SecretKeyLabel         string  `json:"secretKeyLabel"`
	SessionTokenLabel      string  `json:"sessionTokenLabel"`
	MfaDevice              string  `json:"mfaDevice"`
	NamedProfileId         string  `json:"namedProfileId"`
	Status                 domain_aws.AwsSessionStatus `json:"status"`
	StartTime              string  `json:"startTime"`
	LastStopTime           string  `json:"lastStopTime"`
	SessionTokenExpiration string  `json:"sessionTokenExpiration"`
}
