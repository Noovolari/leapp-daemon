package use_case

import (
	"leapp_daemon/domain/domain_gcp/named_configuration"
	"leapp_daemon/test"
	"leapp_daemon/test/mock"
	"net/http"
	"reflect"
	"testing"
)

var (
	namedConfigurations             []named_configuration.NamedConfiguration
	namedConfigurationActionEnvMock mock.EnvironmentMock
	namedConfigurationFacadeMock    mock.NamedConfigurationsFacadeMock
	namedConfigurationActions       *NamedConfigurationsActions
)

func namedConfigurationActionsSetup() {
	namedConfigurations = []named_configuration.NamedConfiguration{
		{Id: "ID1", Name: "configurationName1"},
		{Id: "ID2", Name: "configurationName2"},
	}

	namedConfigurationActionEnvMock = mock.NewEnvironmentMock()
	namedConfigurationFacadeMock = mock.NewNamedConfigurationsFacadeMock()
	gcpIamUserAccountOauthSessionActionsKeychainMock = mock.NewKeychainMock()
	namedConfigurationActions = &NamedConfigurationsActions{
		Environment:               &namedConfigurationActionEnvMock,
		NamedConfigurationsFacade: &namedConfigurationFacadeMock,
	}
}

func namedConfigurationActionsVerifyExpectedCalls(t *testing.T, envMockCalls, namedConfigurationFacadeMockCalls []string) {
	if !reflect.DeepEqual(namedConfigurationActionEnvMock.GetCalls(), envMockCalls) {
		t.Fatalf("envMock expectation violation.\nMock calls: %v", namedConfigurationActionEnvMock.GetCalls())
	}
	if !reflect.DeepEqual(namedConfigurationFacadeMock.GetCalls(), namedConfigurationFacadeMockCalls) {
		t.Fatalf("namedConfigurationFacadeMock expectation violation.\nMock calls: %v", namedConfigurationFacadeMock.GetCalls())
	}
}

func TestGetNamedConfigurations(t *testing.T) {
	namedConfigurationActionsSetup()
	namedConfigurationFacadeMock.ExpNamedConfigurations = namedConfigurations

	actualConfigurations := namedConfigurationActions.GetNamedConfigurations()
	if !reflect.DeepEqual(actualConfigurations, namedConfigurations) {
		t.Fatalf("Returned unexpected named configurations")
	}
	namedConfigurationActionsVerifyExpectedCalls(t, []string{}, []string{"GetNamedConfigurations()"})
}

func TestGetNamedConfigurationById(t *testing.T) {
	namedConfigurationActionsSetup()
	namedConfigurationFacadeMock.ExpNamedConfiguration = namedConfigurations[0]

	actualConfiguration, err := namedConfigurationActions.GetNamedConfigurationById("ID1")
	if err != nil {
		t.Fatalf("Returned unexpected error")

	}
	if !reflect.DeepEqual(actualConfiguration, namedConfigurations[0]) {
		t.Fatalf("Returned unexpected named configuration")
	}
	namedConfigurationActionsVerifyExpectedCalls(t, []string{}, []string{"GetNamedConfigurationById(ID1)"})
}

func TestGetNamedConfigurationById_facadeReturnsError(t *testing.T) {
	namedConfigurationActionsSetup()
	namedConfigurationFacadeMock.ExpErrorOnGetNamedConfigurationById = true

	_, err := namedConfigurationActions.GetNamedConfigurationById("ID1")
	test.ExpectHttpError(t, err, http.StatusNotFound, "named configuration not found")
	namedConfigurationActionsVerifyExpectedCalls(t, []string{}, []string{"GetNamedConfigurationById(ID1)"})
}

func TestGetOrCreateNamedConfiguration(t *testing.T) {
	namedConfigurationActionsSetup()
	namedConfigurationFacadeMock.ExpNamedConfiguration = namedConfigurations[0]

	actualConfiguration, err := namedConfigurationActions.GetOrCreateNamedConfiguration("configurationName1")
	if err != nil {
		return
	}
	if !reflect.DeepEqual(actualConfiguration, namedConfigurations[0]) {
		t.Fatalf("Returned unexpected named configuration")
	}
	namedConfigurationActionsVerifyExpectedCalls(t, []string{}, []string{"GetNamedConfigurationByName(configurationName1)"})
}

func TestGetOrCreateNamedConfiguration_DefaultNamedConfiguration(t *testing.T) {
	namedConfigurationActionsSetup()
	namedConfigurationFacadeMock.ExpNamedConfiguration = namedConfigurations[0]

	actualConfiguration, err := namedConfigurationActions.GetOrCreateNamedConfiguration("")
	if err != nil {
		return
	}
	if !reflect.DeepEqual(actualConfiguration, namedConfigurations[0]) {
		t.Fatalf("Returned unexpected named configuration")
	}
	namedConfigurationActionsVerifyExpectedCalls(t, []string{}, []string{"GetNamedConfigurationByName(leapp-default)"})
}

func TestGetOrCreateNamedConfiguration_FacadeGetByNameReturnsError(t *testing.T) {
	namedConfigurationActionsSetup()
	namedConfigurationFacadeMock.ExpErrorOnGetNamedConfigurationByName = true
	namedConfigurationActionEnvMock.ExpUuid = "NewID"
	expectedConfiguration := named_configuration.NamedConfiguration{Id: "NewID", Name: "newConfigurationName"}

	actualConfiguration, err := namedConfigurationActions.GetOrCreateNamedConfiguration("newConfigurationName")
	if err != nil {
		return
	}
	if !reflect.DeepEqual(actualConfiguration, expectedConfiguration) {
		t.Fatalf("Returned unexpected named configuration")
	}
	namedConfigurationActionsVerifyExpectedCalls(t, []string{"GenerateUuid()"},
		[]string{"GetNamedConfigurationByName(newConfigurationName)", "AddNamedConfiguration({NewID newConfigurationName})"})
}

func TestGetOrCreateNamedConfiguration_FacadeAddNamedConfigurationReturnsError(t *testing.T) {
	namedConfigurationActionsSetup()
	namedConfigurationFacadeMock.ExpErrorOnGetNamedConfigurationByName = true
	namedConfigurationFacadeMock.ExpErrorOnAddNamedConfiguration = true
	namedConfigurationActionEnvMock.ExpUuid = "NewID"

	_, err := namedConfigurationActions.GetOrCreateNamedConfiguration("newConfigurationName")
	test.ExpectHttpError(t, err, http.StatusConflict, "unable to add named configuration")

	namedConfigurationActionsVerifyExpectedCalls(t, []string{"GenerateUuid()"},
		[]string{"GetNamedConfigurationByName(newConfigurationName)", "AddNamedConfiguration({NewID newConfigurationName})"})
}
