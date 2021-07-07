package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	aws_domain "leapp_daemon/domain/domain_aws"
)

const (
	SessionTokenDurationInSeconds int64 = 3600
)

type StsApi struct {
}

func (stsApi *StsApi) GenerateNewSessionToken(accessKeyId string, secretKey string, region string,
	mfaDevice string, mfaToken *string) (*sts.Credentials, error) {

	staticCredentialsClient, err := stsApi.getStaticCredentialsClient(accessKeyId, secretKey, region)
	if err != nil {
		return nil, err
	}

	sessionTokenInput := stsApi.getSessionTokenInput(mfaDevice, mfaToken)
	sessionTokenOutput, err := staticCredentialsClient.GetSessionToken(&sessionTokenInput)
	if err != nil {
		return nil, err
	}

	return sessionTokenOutput.Credentials, nil
}

func (stsApi *StsApi) getStaticCredentialsClient(accessKeyId string, secretKey string, region string) (*sts.STS, error) {
	endpoint, err := aws_domain.GetRegionalEndpoint(region)
	if err != nil {
		return nil, err
	}

	stsConfig := &aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKeyId, secretKey, ""),
		Region:      &region,
		Endpoint:    &endpoint,
	}

	newSession, err := session.NewSession(stsConfig)
	sts := sts.New(session.Must(newSession, err))
	return sts, nil
}

func (stsApi *StsApi) getSessionTokenInput(mfaDevice string, mfaToken *string) sts.GetSessionTokenInput {
	tokenDuration := SessionTokenDurationInSeconds

	if mfaToken == nil {
		return sts.GetSessionTokenInput{
			DurationSeconds: &tokenDuration,
			SerialNumber:    nil,
			TokenCode:       nil,
		}
	}

	return sts.GetSessionTokenInput{
		DurationSeconds: &tokenDuration,
		SerialNumber:    &mfaDevice,
		TokenCode:       mfaToken,
	}
}
