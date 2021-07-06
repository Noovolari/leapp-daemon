// TODO: probably, it should be moved in interface layer, since it acts as an interface between a domain concept,
//  e.g. access keys, and an external service.

package keychain

import (
	"leapp_daemon/infrastructure/http/http_error"

	"github.com/zalando/go-keyring"
)

const ServiceName = "Leapp"

type Keychain struct{}

func (keychain *Keychain) SetSecret(secret string, label string) error {
	err := keyring.Set(ServiceName, label, secret)
	if err != nil {
		return http_error.NewUnprocessableEntityError(err)
	}
	return nil
}

func (keychain *Keychain) GetSecret(label string) (string, error) {
	secret, err := keyring.Get(ServiceName, label)
	if err != nil {
		return "", http_error.NewNotFoundError(err)
	}
	return secret, nil
}

func (keychain *Keychain) DeleteSecret(label string) error {
	err := keyring.Delete(ServiceName, label)
	if err != nil {
		return http_error.NewNotFoundError(err)
	}
	return nil
}

func (keychain *Keychain) DoesSecretExist(label string) (bool, error) {
	_, err := keyring.Get(ServiceName, label)
	if err != nil {
		if err.Error() == "secret not found in keyring" {
			return false, nil
		}
		return false, http_error.NewUnprocessableEntityError(err)
	}
	return true, nil
}
