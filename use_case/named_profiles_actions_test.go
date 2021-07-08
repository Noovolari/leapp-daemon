package use_case

import (
	"leapp_daemon/domain/domain_aws/named_profile"
	"leapp_daemon/test"
	"leapp_daemon/test/mock"
	"net/http"
	"reflect"
	"testing"
)

var (
	namedProfiles              []named_profile.NamedProfile
	namedProfilesActionEnvMock mock.EnvironmentMock
	namedProfilesFacadeMock    mock.NamedProfilesFacadeMock
	namedProfilesActions       *NamedProfilesActions
)

func namedProfilesActionsSetup() {
	namedProfiles = []named_profile.NamedProfile{
		{Id: "ID1", Name: "profileName1"},
		{Id: "ID2", Name: "profileName2"},
	}

	namedProfilesActionEnvMock = mock.NewEnvironmentMock()
	namedProfilesFacadeMock = mock.NewNamedProfilesFacadeMock()
	namedProfilesActions = &NamedProfilesActions{
		Environment:         &namedProfilesActionEnvMock,
		NamedProfilesFacade: &namedProfilesFacadeMock,
	}
}

func namedProfilesActionsVerifyExpectedCalls(t *testing.T, envMockCalls, namedProfilesFacadeMockCalls []string) {
	if !reflect.DeepEqual(namedProfilesActionEnvMock.GetCalls(), envMockCalls) {
		t.Fatalf("envMock expectation violation.\nMock calls: %v", namedProfilesActionEnvMock.GetCalls())
	}
	if !reflect.DeepEqual(namedProfilesFacadeMock.GetCalls(), namedProfilesFacadeMockCalls) {
		t.Fatalf("namedProfilesFacadeMock expectation violation.\nMock calls: %v", namedProfilesFacadeMock.GetCalls())
	}
}

func TestGetNamedProfiles(t *testing.T) {
	namedProfilesActionsSetup()
	namedProfilesFacadeMock.ExpNamedProfiles = namedProfiles

	actualProfiles := namedProfilesActions.GetNamedProfiles()
	if !reflect.DeepEqual(actualProfiles, namedProfiles) {
		t.Fatalf("Returned unexpected named profiles")
	}
	namedProfilesActionsVerifyExpectedCalls(t, []string{}, []string{"GetNamedProfiles()"})
}

func TestGetNamedProfileById(t *testing.T) {
	namedProfilesActionsSetup()
	namedProfilesFacadeMock.ExpNamedProfile = namedProfiles[0]

	actualConfiguration, err := namedProfilesActions.GetNamedProfileById("ID1")
	if err != nil {
		t.Fatalf("Returned unexpected error")

	}
	if !reflect.DeepEqual(actualConfiguration, namedProfiles[0]) {
		t.Fatalf("Returned unexpected named profile")
	}
	namedProfilesActionsVerifyExpectedCalls(t, []string{}, []string{"GetNamedProfileById(ID1)"})
}

func TestGetNamedProfileById_facadeReturnsError(t *testing.T) {
	namedProfilesActionsSetup()
	namedProfilesFacadeMock.ExpErrorOnGetNamedProfileById = true

	_, err := namedProfilesActions.GetNamedProfileById("ID1")
	test.ExpectHttpError(t, err, http.StatusNotFound, "named profile not found")
	namedProfilesActionsVerifyExpectedCalls(t, []string{}, []string{"GetNamedProfileById(ID1)"})
}

func TestGetOrCreateNamedProfile(t *testing.T) {
	namedProfilesActionsSetup()
	namedProfilesFacadeMock.ExpNamedProfile = namedProfiles[0]

	actualConfiguration, err := namedProfilesActions.GetOrCreateNamedProfile("profileName1")
	if err != nil {
		return
	}
	if !reflect.DeepEqual(actualConfiguration, namedProfiles[0]) {
		t.Fatalf("Returned unexpected named profile")
	}
	namedProfilesActionsVerifyExpectedCalls(t, []string{}, []string{"GetNamedProfileByName(profileName1)"})
}

func TestGetOrCreateNamedProfile_DefaultNamedProfileName(t *testing.T) {
	namedProfilesActionsSetup()
	namedProfilesFacadeMock.ExpNamedProfile = namedProfiles[0]

	actualConfiguration, err := namedProfilesActions.GetOrCreateNamedProfile("")
	if err != nil {
		return
	}
	if !reflect.DeepEqual(actualConfiguration, namedProfiles[0]) {
		t.Fatalf("Returned unexpected named profile")
	}
	namedProfilesActionsVerifyExpectedCalls(t, []string{}, []string{"GetNamedProfileByName(default)"})
}

func TestGetOrCreateNamedProfile_FacadeGetByNameReturnsError(t *testing.T) {
	namedProfilesActionsSetup()
	namedProfilesFacadeMock.ExpErrorOnGetNamedProfileByName = true
	namedProfilesActionEnvMock.ExpUuid = "NewID"
	expectedConfiguration := named_profile.NamedProfile{Id: "NewID", Name: "newProfileName"}

	actualConfiguration, err := namedProfilesActions.GetOrCreateNamedProfile("newProfileName")
	if err != nil {
		return
	}
	if !reflect.DeepEqual(actualConfiguration, expectedConfiguration) {
		t.Fatalf("Returned unexpected named profile")
	}
	namedProfilesActionsVerifyExpectedCalls(t, []string{"GenerateUuid()"},
		[]string{"GetNamedProfileByName(newProfileName)", "AddNamedProfile({NewID newProfileName})"})
}

func TestGetOrCreateNamedProfile_FacadeAddNamedConfigurationReturnsError(t *testing.T) {
	namedProfilesActionsSetup()
	namedProfilesFacadeMock.ExpErrorOnGetNamedProfileByName = true
	namedProfilesFacadeMock.ExpErrorOnAddNamedProfile = true
	namedProfilesActionEnvMock.ExpUuid = "NewID"

	_, err := namedProfilesActions.GetOrCreateNamedProfile("newProfileName")
	test.ExpectHttpError(t, err, http.StatusInternalServerError, "failed to add named profile")

	namedProfilesActionsVerifyExpectedCalls(t, []string{"GenerateUuid()"},
		[]string{"GetNamedProfileByName(newProfileName)", "AddNamedProfile({NewID newProfileName})"})
}
