package aws

type AwsTempCredentials struct {
	ProfileName  string
	AccessKeyId  string
	SecretKey    string
	SessionToken string
	Expiration   string
	Region       string
}
