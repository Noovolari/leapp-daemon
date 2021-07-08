package named_profile

import (
	"errors"
	"leapp_daemon/infrastructure/http/http_error"
	"leapp_daemon/test"
	"net/http"
	"reflect"
	"testing"
)

var (
	facadeNamedProfiles       []NamedProfile
	namedProfilesBeforeUpdate []NamedProfile
	namedProfilesAfterUpdate  []NamedProfile
	facade                    *NamedProfilesFacade
)

func facadeSetup() {
	facade = NewNamedProfilesFacade()
	facadeNamedProfiles = []NamedProfile{
		{Id: "ID1", Name: "profile-name1"},
		{Id: "ID2", Name: "profile-name2"},
	}
	namedProfilesBeforeUpdate = []NamedProfile{}
	namedProfilesAfterUpdate = []NamedProfile{}
}

func TestSetAndGetNamedProfiles(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeNamedProfileObserver{})
	facade.SetNamedProfiles(facadeNamedProfiles)

	actualProfiles := facade.GetNamedProfiles()
	if !reflect.DeepEqual(facadeNamedProfiles, actualProfiles) {
		t.Errorf("unexpected profiles")
	}

	if !reflect.DeepEqual(namedProfilesBeforeUpdate, []NamedProfile{}) ||
		!reflect.DeepEqual(namedProfilesAfterUpdate, []NamedProfile{}) {
		t.Errorf("unexpected call to updateState")
	}
}

func TestNamedProfilesFacade_AddNamedProfile(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeNamedProfileObserver{})
	facade.SetNamedProfiles(facadeNamedProfiles)
	namedConfigurationToAdd := NamedProfile{Id: "ID3", Name: "profile-name3"}

	err := facade.AddNamedProfile(namedConfigurationToAdd)
	if err != nil {
		t.Errorf("unexpected error")
	}

	if !reflect.DeepEqual(facade.namedProfiles, append(facadeNamedProfiles, namedConfigurationToAdd)) {
		t.Errorf("unexpected profiles")
	}

	if !reflect.DeepEqual(namedProfilesBeforeUpdate, append(facadeNamedProfiles)) {
		t.Errorf("unexpected profiles before update")
	}
	if !reflect.DeepEqual(namedProfilesAfterUpdate, append(facadeNamedProfiles, namedConfigurationToAdd)) {
		t.Errorf("unexpected profiles after update")
	}
}

func TestNamedProfilesFacade_AddNamedProfile_IDAlreadyPresent(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeNamedProfileObserver{})
	facade.SetNamedProfiles(facadeNamedProfiles)
	namedConfigurationToAdd := NamedProfile{Id: "ID1", Name: "profile-name3"}

	err := facade.AddNamedProfile(namedConfigurationToAdd)
	test.ExpectHttpError(t, err, http.StatusConflict, "a NamedProfile with id ID1 is already present")
	if !reflect.DeepEqual(namedProfilesBeforeUpdate, []NamedProfile{}) ||
		!reflect.DeepEqual(namedProfilesAfterUpdate, []NamedProfile{}) {
		t.Errorf("unexpected call to updateState")
	}
}

func TestNamedProfilesFacade_AddNamedProfile_NameAlreadyPresent(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeNamedProfileObserver{})
	facade.SetNamedProfiles(facadeNamedProfiles)
	namedConfigurationToAdd := NamedProfile{Id: "ID3", Name: "profile-name1"}

	err := facade.AddNamedProfile(namedConfigurationToAdd)
	test.ExpectHttpError(t, err, http.StatusConflict, "a NamedProfile with name profile-name1 is already present")
	if !reflect.DeepEqual(namedProfilesBeforeUpdate, []NamedProfile{}) ||
		!reflect.DeepEqual(namedProfilesAfterUpdate, []NamedProfile{}) {
		t.Errorf("unexpected call to updateState")
	}
}

func TestNamedProfilesFacade_AddNamedProfile_UpdateStateReturnsError(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeNamedProfileObserver{ExpErrorOnUpdateState: true})
	facade.SetNamedProfiles(facadeNamedProfiles)
	namedConfigurationToAdd := NamedProfile{Id: "ID3", Name: "profile-name3"}

	err := facade.AddNamedProfile(namedConfigurationToAdd)
	test.ExpectHttpError(t, err, http.StatusInternalServerError, "unable to update named profiles")

	if !reflect.DeepEqual(facade.namedProfiles, append(facadeNamedProfiles, namedConfigurationToAdd)) {
		t.Errorf("unexpected profiles")
	}

	if !reflect.DeepEqual(namedProfilesBeforeUpdate, append(facadeNamedProfiles)) {
		t.Errorf("unexpected profiles before update")
	}
	if !reflect.DeepEqual(namedProfilesAfterUpdate, append(facadeNamedProfiles, namedConfigurationToAdd)) {
		t.Errorf("unexpected profiles after update")
	}
}

func TestNamedProfilesFacade_GetNamedProfileByName(t *testing.T) {
	facadeSetup()
	facade.SetNamedProfiles(facadeNamedProfiles)

	namedConfiguration, err := facade.GetNamedProfileByName("profile-name1")
	if err != nil {
		t.Errorf("unexpected error")
	}

	if !reflect.DeepEqual(namedConfiguration, facadeNamedProfiles[0]) {
		t.Errorf("unexpected profile")
	}
}

func TestNamedProfilesFacade_GetNamedProfileByName_NamedProfileNotFound(t *testing.T) {
	facadeSetup()
	facade.SetNamedProfiles(facadeNamedProfiles)

	_, err := facade.GetNamedProfileByName("profile-other-name")
	test.ExpectHttpError(t, err, http.StatusNotFound, "named profile with name profile-other-name not found")
}

func TestNamedProfilesFacade_GetNamedProfileById(t *testing.T) {
	facadeSetup()
	facade.SetNamedProfiles(facadeNamedProfiles)

	namedConfiguration, err := facade.GetNamedProfileById("ID1")
	if err != nil {
		t.Errorf("unexpected error")
	}

	if !reflect.DeepEqual(namedConfiguration, facadeNamedProfiles[0]) {
		t.Errorf("unexpected profile")
	}
}

func TestNamedProfilesFacade_GetNamedProfileById_NamedProfileNotFound(t *testing.T) {
	facadeSetup()
	facade.SetNamedProfiles(facadeNamedProfiles)

	_, err := facade.GetNamedProfileById("OTHER-ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "named profile with id OTHER-ID not found")
}

type fakeNamedProfileObserver struct {
	ExpErrorOnUpdateState bool
}

func (fakeObserver fakeNamedProfileObserver) UpdateNamedProfiles(oldNamedProfiles []NamedProfile, newNamedProfiles []NamedProfile) error {
	namedProfilesBeforeUpdate = oldNamedProfiles
	namedProfilesAfterUpdate = newNamedProfiles

	if fakeObserver.ExpErrorOnUpdateState {
		return http_error.NewInternalServerError(errors.New("unable to update named profiles"))
	}

	return nil
}
