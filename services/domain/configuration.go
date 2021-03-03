package domain

type Configuration struct {
	SsoUrl string
	ProxyConfiguration ProxyConfiguration
	FederatedAwsAccountSessions []FederatedAwsAccountSession
	PlainAwsAccountSessions []PlainAwsAccountSession
}

type ProxyConfiguration struct {
	ProxyProtocol string
	ProxyUrl string
	ProxyPort uint64
	Username string
	Password string
}

type FederatedAwsAccountSession struct {
	Id           string
	Active       bool
	Loading      bool
	LastStopDate string
	Account      FederatedAwsAccount
}

type PlainAwsAccountSession struct {
	Id           string
	Active       bool
	Loading      bool
	LastStopDate string
	Account      PlainAwsAccount
}

type FederatedAwsAccount struct {
	AccountNumber string
	Name          string
	Role          FederatedAwsRole
	IdpArn        string
	Region        string
	SsoUrl        string
	// Type            string
	// ParentSessionId string
	// ParentRole      string
}

type FederatedAwsRole struct {
	Name string
	Arn  string
	// Parent string
	// ParentRole string
}

type PlainAwsAccount struct {
	AccountNumber string
	Name          string
	Region        string
	User          string
	MfaDevice     string
}