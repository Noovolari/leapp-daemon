package gcp_iam_user_account_oauth

import "leapp_daemon/domain/domain_gcp"

type GcpIamUserAccountOauthSession struct {
	Id               string
	Name             string
	AccountId        string
	ProjectName      string
	CredentialsLabel string
	Status           domain_gcp.GcpSessionStatus
	StartTime        string
	LastStopTime     string
}
