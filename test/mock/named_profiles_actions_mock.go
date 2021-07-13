package mock

import (
	"errors"
	"fmt"
	"leapp_daemon/domain/named_profile"
	"leapp_daemon/infrastructure/http/http_error"
)

type NamedProfilesActionsMock struct {
	calls                             []string
	ExpErrorOnGetOrCreateNamedProfile bool
	ExpNamedProfile                   named_profile.NamedProfile
}

func NewNamedProfilesActionsMock() NamedProfilesActionsMock {
	return NamedProfilesActionsMock{calls: []string{}}
}

func (actions *NamedProfilesActionsMock) GetCalls() []string {
	return actions.calls
}

func (actions *NamedProfilesActionsMock) GetOrCreateNamedProfile(profileName string) (named_profile.NamedProfile, error) {
	actions.calls = append(actions.calls, fmt.Sprintf("GetOrCreateNamedProfile(%v)", profileName))
	if actions.ExpErrorOnGetOrCreateNamedProfile {
		return named_profile.NamedProfile{}, http_error.NewNotFoundError(errors.New("named profile not found"))
	}

	return actions.ExpNamedProfile, nil
}
