package mock

import (
	"errors"
	"fmt"
	"leapp_daemon/infrastructure/http/http_error"
)

type KeychainMock struct {
	calls                  []string
	ExpErrorOnSetSecret    bool
	ExpErrorOnGetSecret    bool
	ExpErrorOnDeleteSecret bool
	ExpErrorOnSecretExist  bool
	ExpGetSecret           string
	ExpSecretExist         bool
}

func NewKeychainMock() KeychainMock {
	return KeychainMock{calls: []string{}}
}

func (chainMock *KeychainMock) GetCalls() []string {
	return chainMock.calls
}

func (chainMock *KeychainMock) SetSecret(secret string, label string) error {
	chainMock.calls = append(chainMock.calls, fmt.Sprintf("SetSecret(%v, %v)", secret, label))
	if chainMock.ExpErrorOnSetSecret {
		return http_error.NewUnprocessableEntityError(errors.New("unable to set secret"))
	}

	return nil
}

func (chainMock *KeychainMock) GetSecret(label string) (string, error) {
	chainMock.calls = append(chainMock.calls, fmt.Sprintf("GetSecret(%v)", label))
	if chainMock.ExpErrorOnGetSecret {
		return "", http_error.NewUnprocessableEntityError(nil)
	}

	return chainMock.ExpGetSecret, nil
}

func (chainMock *KeychainMock) DeleteSecret(label string) error {
	chainMock.calls = append(chainMock.calls, fmt.Sprintf("DeleteSecret(%v)", label))
	if chainMock.ExpErrorOnDeleteSecret {
		return http_error.NewNotFoundError(nil)
	}

	return nil
}

func (chainMock *KeychainMock) DoesSecretExist(label string) (bool, error) {
	chainMock.calls = append(chainMock.calls, fmt.Sprintf("DoesSecretExist(%v)", label))
	if chainMock.ExpErrorOnSecretExist {
		return false, http_error.NewUnprocessableEntityError(nil)
	}

	return chainMock.ExpSecretExist, nil
}
