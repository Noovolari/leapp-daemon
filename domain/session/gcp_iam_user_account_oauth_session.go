package session

type GcpIamUserAccountOauthSession struct {
	Id               string
	Name             string
	AccountId        string
	ProjectName      string
	CredentialsLabel string
	Status           Status
	StartTime        string
	LastStopTime     string
}
