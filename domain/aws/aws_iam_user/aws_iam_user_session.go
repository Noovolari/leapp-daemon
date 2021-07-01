package aws_iam_user

import (
	"leapp_daemon/domain/aws"
)

type AwsIamUserSessionContainer interface {
	AddAwsIamUserSession(AwsIamUserSession) error
	GetAllAwsIamUserSessions() ([]AwsIamUserSession, error)
	RemoveAwsIamUserSession(session AwsIamUserSession) error
}

type AwsIamUserSession struct {
	Id                     string
	Name                   string
	Region                 string
	AccountNumber          string
	UserName               string
	AccessKeyIdLabel       string
	SecretKeyLabel         string
	SessionTokenLabel      string
	MfaDevice              string
	NamedProfileId         string
	Status                 aws.AwsSessionStatus
	StartTime              string
	LastStopTime           string
	SessionTokenExpiration string
}
