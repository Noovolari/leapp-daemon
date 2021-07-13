package mock

import (
	"errors"
	"fmt"
	"leapp_daemon/domain/named_profile"
	"leapp_daemon/infrastructure/http/http_error"
)

type NamedProfilesFacadeMock struct {
	calls                           []string
	ExpNamedProfiles                []named_profile.NamedProfile
	ExpNamedProfile                 named_profile.NamedProfile
	ExpErrorOnAddNamedProfile       bool
	ExpErrorOnGetNamedProfileByName bool
	ExpErrorOnGetNamedProfileById   bool
}

func NewNamedProfilesFacadeMock() NamedProfilesFacadeMock {
	return NamedProfilesFacadeMock{calls: []string{}}
}

func (facade *NamedProfilesFacadeMock) GetCalls() []string {
	return facade.calls
}

func (facade *NamedProfilesFacadeMock) Subscribe(observer named_profile.NamedProfilesObserver) {
	facade.calls = append(facade.calls, "Subscribe()")
}

func (facade *NamedProfilesFacadeMock) GetNamedProfiles() []named_profile.NamedProfile {
	facade.calls = append(facade.calls, "GetNamedProfiles()")
	return facade.ExpNamedProfiles
}

func (facade *NamedProfilesFacadeMock) SetNamedProfiles(namedProfiles []named_profile.NamedProfile) {
	facade.calls = append(facade.calls, fmt.Sprintf("SetNamedProfiles(%v)", namedProfiles))
}

func (facade *NamedProfilesFacadeMock) AddNamedProfile(namedProfile named_profile.NamedProfile) error {
	facade.calls = append(facade.calls, fmt.Sprintf("AddNamedProfile(%v)", namedProfile.Id))
	if facade.ExpErrorOnAddNamedProfile {
		return http_error.NewInternalServerError(errors.New("failed to add named profile"))
	}

	return nil
}

func (facade *NamedProfilesFacadeMock) GetNamedProfileByName(profileName string) (named_profile.NamedProfile, error) {
	facade.calls = append(facade.calls, fmt.Sprintf("GetNamedProfileByName(%v)", profileName))
	if facade.ExpErrorOnAddNamedProfile {
		return named_profile.NamedProfile{}, http_error.NewNotFoundError(errors.New("named profile not found"))
	}

	return facade.ExpNamedProfile, nil
}

func (facade *NamedProfilesFacadeMock) GetNamedProfileById(id string) (named_profile.NamedProfile, error) {
	facade.calls = append(facade.calls, fmt.Sprintf("GetNamedProfileById(%v)", id))
	if facade.ExpErrorOnGetNamedProfileById {
		return named_profile.NamedProfile{}, http_error.NewNotFoundError(errors.New("named profile not found"))
	}

	return facade.ExpNamedProfile, nil
}
