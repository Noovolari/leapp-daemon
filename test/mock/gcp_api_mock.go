package mock

import (
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"leapp_daemon/infrastructure/http/http_error"
)

type GcpApiMock struct {
	calls                 []string
	ExpErrorOnGetOauthUrl bool
	ExpErrorOnGetOauth    bool
	ExpOauthUrl           string
	ExpOauthToken         oauth2.Token
	ExpCredentials        string
}

func NewGcpApiMock() GcpApiMock {
	return GcpApiMock{calls: []string{}}
}

func (gcpApi *GcpApiMock) GetCalls() []string {
	return gcpApi.calls
}

func (gcpApi *GcpApiMock) GetOauthUrl() (string, error) {
	gcpApi.calls = append(gcpApi.calls, "GetOauthUrl()")
	if gcpApi.ExpErrorOnGetOauthUrl {
		return "", http_error.NewNotFoundError(errors.New("error getting oauth url"))
	}

	return gcpApi.ExpOauthUrl, nil
}

func (gcpApi *GcpApiMock) GetOauthToken(authCode string) (*oauth2.Token, error) {
	gcpApi.calls = append(gcpApi.calls, fmt.Sprintf("GetOauthUrl(%v)", authCode))
	if gcpApi.ExpErrorOnGetOauth {
		return &oauth2.Token{}, http_error.NewBadRequestError(errors.New("error getting oauth token"))
	}

	return &gcpApi.ExpOauthToken, nil
}

func (gcpApi *GcpApiMock) GetCredentials(oauthToken *oauth2.Token) string {
	gcpApi.calls = append(gcpApi.calls, "GetCredentials()")
	return gcpApi.ExpCredentials
}
