package aws_iam_user

import (
	"leapp_daemon/domain/domain_aws"
)

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
	Status                 domain_aws.AwsSessionStatus
	StartTime              string
	LastStopTime           string
	SessionTokenExpiration string
}
