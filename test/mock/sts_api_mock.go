package mock

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/sts"
	"leapp_daemon/infrastructure/http/http_error"
)

type StsApiMock struct {
	calls                             []string
	ExpErrorOnGenerateNewSessionToken bool
	ExpCredentials                    sts.Credentials
}

func NewStsApiMock() StsApiMock {
	return StsApiMock{calls: []string{}}
}

func (stsApi *StsApiMock) GetCalls() []string {
	return stsApi.calls
}

func (stsApi *StsApiMock) GenerateNewSessionToken(accessKeyId string, secretKey string, region string,
	mfaDevice string, mfaToken *string) (*sts.Credentials, error) {
	stsApi.calls = append(stsApi.calls, fmt.Sprintf("GenerateNewSessionToken(%v, %v, %v, %v, %v)",
		accessKeyId, secretKey, region, mfaDevice, mfaToken))

	if stsApi.ExpErrorOnGenerateNewSessionToken {
		return &sts.Credentials{}, http_error.NewUnprocessableEntityError(errors.New("unable to generate new session token"))
	}

	return &stsApi.ExpCredentials, nil
}
