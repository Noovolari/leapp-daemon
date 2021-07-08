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
	namedConfigurations              []named_configuration.NamedConfiguration
	namedConfigurationsActionEnvMock mock.EnvironmentMock
	namedConfigurationsFacadeMock    mock.NamedConfigurationsFacadeMock
	namedConfigurationsActions       *NamedConfigurationsActions
)

func namedConfigurationActionsSetup() {
	namedConfigurations = []named_configuration.NamedConfiguration{
		{Id: "ID1", Name: "configurationName1"},
		{Id: "ID2", Name: "configurationName2"},
	}

	namedConfigurationsActionEnvMock = mock.NewEnvironmentMock()
	namedConfigurationsFacadeMock = mock.NewNamedConfigurationsFacadeMock()
	gcpIamUserAccountOauthSessionActionsKeychainMock = mock.NewKeychainMock()
	namedConfigurationsActions = &NamedConfigurationsActions{
		Environment:               &namedConfigurationsActionEnvMock,
		NamedConfigurationsFacade: &namedConfigurationsFacadeMock,
	}
}

func namedConfigurationActionsVerifyExpectedCalls(t *testing.T, envMockCalls, namedConfigurationsFacadeMockCalls []string) {
	if !reflect.DeepEqual(namedConfigurationsActionEnvMock.GetCalls(), envMockCalls) {
		t.Fatalf("envMock expectation violation.\nMock calls: %v", namedConfigurationsActionEnvMock.GetCalls())
	}
	if !reflect.DeepEqual(namedConfigurationsFacadeMock.GetCalls(), namedConfigurationsFacadeMockCalls) {
		t.Fatalf("namedConfigurationsFacadeMock expectation violation.\nMock calls: %v", namedConfigurationsFacadeMock.GetCalls())
	}
}

func TestGetNamedConfigurations(t *testing.T) {
	namedConfigurationActionsSetup()
	namedConfigurationsFacadeMock.ExpNamedConfigurations = namedConfigurations

	actualConfigurations := namedConfigurationsActions.GetNamedConfigurations()
	if !reflect.DeepEqual(actualConfigurations, namedConfigurations) {
		t.Fatalf("Returned unexpected named configurations")
	}
	namedConfigurationActionsVerifyExpectedCalls(t, []string{}, []string{"GetNamedConfigurations()"})
}

func TestGetNamedConfigurationById(t *testing.T) {
	namedConfigurationActionsSetup()
	namedConfigurationsFacadeMock.ExpNamedConfiguration = namedConfigurations[0]

	actualConfiguration, err := namedConfigurationsActions.GetNamedConfigurationById("ID1")
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
	namedConfigurationsFacadeMock.ExpErrorOnGetNamedConfigurationById = true

	_, err := namedConfigurationsActions.GetNamedConfigurationById("ID1")
	test.ExpectHttpError(t, err, http.StatusNotFound, "named configuration not found")
	namedConfigurationActionsVerifyExpectedCalls(t, []string{}, []string{"GetNamedConfigurationById(ID1)"})
}

func TestGetOrCreateNamedConfiguration(t *testing.T) {
	namedConfigurationActionsSetup()
	namedConfigurationsFacadeMock.ExpNamedConfiguration = namedConfigurations[0]

	actualConfiguration, err := namedConfigurationsActions.GetOrCreateNamedConfiguration("configurationName1")
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
	namedConfigurationsFacadeMock.ExpNamedConfiguration = namedConfigurations[0]

	actualConfiguration, err := namedConfigurationsActions.GetOrCreateNamedConfiguration("")
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
	namedConfigurationsFacadeMock.ExpErrorOnGetNamedConfigurationByName = true
	namedConfigurationsActionEnvMock.ExpUuid = "NewID"
	expectedConfiguration := named_configuration.NamedConfiguration{Id: "NewID", Name: "newConfigurationName"}

	actualConfiguration, err := namedConfigurationsActions.GetOrCreateNamedConfiguration("newConfigurationName")
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
	namedConfigurationsFacadeMock.ExpErrorOnGetNamedConfigurationByName = true
	namedConfigurationsFacadeMock.ExpErrorOnAddNamedConfiguration = true
	namedConfigurationsActionEnvMock.ExpUuid = "NewID"

	_, err := namedConfigurationsActions.GetOrCreateNamedConfiguration("newConfigurationName")
	test.ExpectHttpError(t, err, http.StatusConflict, "unable to add named configuration")

	namedConfigurationActionsVerifyExpectedCalls(t, []string{"GenerateUuid()"},
		[]string{"GetNamedConfigurationByName(newConfigurationName)", "AddNamedConfiguration({NewID newConfigurationName})"})
}
