package named_configuration

import (
	"errors"
	"leapp_daemon/infrastructure/http/http_error"
	"leapp_daemon/test"
	"net/http"
	"reflect"
	"testing"
)

var (
	facade                         *NamedConfigurationsFacade
	facadeNamedConfigurations      []NamedConfiguration
	namedConfigurationBeforeUpdate []NamedConfiguration
	namedConfigurationAfterUpdate  []NamedConfiguration
)

func facadeSetup() {
	facade = NewNamedConfigurationsFacade()
	facadeNamedConfigurations = []NamedConfiguration{
		{Id: "ID1", Name: "configuration-name1"},
		{Id: "ID2", Name: "configuration-name2"},
	}
	namedConfigurationBeforeUpdate = []NamedConfiguration{}
	namedConfigurationAfterUpdate = []NamedConfiguration{}
}

func TestNamedConfigurationsFacade_SetAndGetNamedConfigurations(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeNamedConfigurationObserver{})
	facade.SetNamedConfigurations(facadeNamedConfigurations)

	actualConfigurations := facade.GetNamedConfigurations()
	if !reflect.DeepEqual(facadeNamedConfigurations, actualConfigurations) {
		t.Errorf("unexpected configurations")
	}

	if !reflect.DeepEqual(namedConfigurationBeforeUpdate, []NamedConfiguration{}) ||
		!reflect.DeepEqual(namedConfigurationAfterUpdate, []NamedConfiguration{}) {
		t.Errorf("unexpected call to updateState")
	}
}

func TestNamedConfigurationsFacade_AddNamedConfiguration(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeNamedConfigurationObserver{})
	facade.SetNamedConfigurations(facadeNamedConfigurations)
	namedConfigurationToAdd := NamedConfiguration{Id: "ID3", Name: "configuration-name3"}

	err := facade.AddNamedConfiguration(namedConfigurationToAdd)
	if err != nil {
		t.Errorf("unexpected error")
	}

	if !reflect.DeepEqual(facade.namedConfigurations, append(facadeNamedConfigurations, namedConfigurationToAdd)) {
		t.Errorf("unexpected configurations")
	}

	if !reflect.DeepEqual(namedConfigurationBeforeUpdate, append(facadeNamedConfigurations)) {
		t.Errorf("unexpected configurations before update")
	}
	if !reflect.DeepEqual(namedConfigurationAfterUpdate, append(facadeNamedConfigurations, namedConfigurationToAdd)) {
		t.Errorf("unexpected configurations after update")
	}
}

func TestNamedConfigurationsFacade_AddNamedConfiguration_InvalidConfigurationName(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeNamedConfigurationObserver{})
	facade.SetNamedConfigurations(facadeNamedConfigurations)
	namedConfigurationToAdd := NamedConfiguration{Id: "ID3", Name: "1InvalidNamE1"}

	err := facade.AddNamedConfiguration(namedConfigurationToAdd)
	test.ExpectHttpError(t, err, http.StatusBadRequest, "configuration names must start with a lower case letter and contain only lower case letters a-z, digits 0-9, and hyphens '-'")
	if !reflect.DeepEqual(namedConfigurationBeforeUpdate, []NamedConfiguration{}) ||
		!reflect.DeepEqual(namedConfigurationAfterUpdate, []NamedConfiguration{}) {
		t.Errorf("unexpected call to updateState")
	}
}

func TestNamedConfigurationsFacade_AddNamedConfiguration_UnivocityCheckFailed(t *testing.T) {
	facadeSetup()
	facade.SetNamedConfigurations(facadeNamedConfigurations)
	facade.Subscribe(fakeNamedConfigurationObserver{})
	namedConfigurationToAdd1 := NamedConfiguration{Id: "ID3", Name: "configuration-name1"}
	namedConfigurationToAdd2 := NamedConfiguration{Id: "ID1", Name: "configuration-name3"}

	err := facade.AddNamedConfiguration(namedConfigurationToAdd1)
	test.ExpectHttpError(t, err, http.StatusConflict, "a ConfigurationName with name configuration-name1 is already present")
	if !reflect.DeepEqual(facade.namedConfigurations, facadeNamedConfigurations) {
		t.Errorf("unexpected configurations")
	}

	err = facade.AddNamedConfiguration(namedConfigurationToAdd2)
	test.ExpectHttpError(t, err, http.StatusConflict, "a ConfigurationName with id ID1 is already present")
	if !reflect.DeepEqual(facade.namedConfigurations, facadeNamedConfigurations) {
		t.Errorf("unexpected configurations")
	}

	if !reflect.DeepEqual(namedConfigurationBeforeUpdate, []NamedConfiguration{}) ||
		!reflect.DeepEqual(namedConfigurationAfterUpdate, []NamedConfiguration{}) {
		t.Errorf("unexpected call to updateState")
	}
}

func TestNamedConfigurationsFacade_AddNamedConfiguration_UpdateStateReturnsError(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeNamedConfigurationObserver{ExpErrorOnUpdateState: true})
	facade.SetNamedConfigurations(facadeNamedConfigurations)
	namedConfigurationToAdd := NamedConfiguration{Id: "ID3", Name: "configuration-name3"}

	err := facade.AddNamedConfiguration(namedConfigurationToAdd)
	test.ExpectHttpError(t, err, http.StatusInternalServerError, "unable to update named configurations")

	if !reflect.DeepEqual(facade.namedConfigurations, append(facadeNamedConfigurations, namedConfigurationToAdd)) {
		t.Errorf("unexpected configurations")
	}

	if !reflect.DeepEqual(namedConfigurationBeforeUpdate, append(facadeNamedConfigurations)) {
		t.Errorf("unexpected configurations before update")
	}
	if !reflect.DeepEqual(namedConfigurationAfterUpdate, append(facadeNamedConfigurations, namedConfigurationToAdd)) {
		t.Errorf("unexpected configurations after update")
	}
}

func TestNamedConfigurationsFacade_GetNamedConfigurationByName(t *testing.T) {
	facadeSetup()
	facade.SetNamedConfigurations(facadeNamedConfigurations)

	namedConfiguration, err := facade.GetNamedConfigurationByName("configuration-name1")
	if err != nil {
		t.Errorf("unexpected error")
	}

	if !reflect.DeepEqual(namedConfiguration, facadeNamedConfigurations[0]) {
		t.Errorf("unexpected configuration")
	}
}

func TestNamedConfigurationsFacade_GetNamedConfigurationByName_NamedConfigurationNotFound(t *testing.T) {
	facadeSetup()
	facade.SetNamedConfigurations(facadeNamedConfigurations)

	_, err := facade.GetNamedConfigurationByName("configuration-other-name")
	test.ExpectHttpError(t, err, http.StatusNotFound, "named configuration with name configuration-other-name not found")
}

func TestNamedConfigurationsFacade_GetNamedConfigurationById(t *testing.T) {
	facadeSetup()
	facade.SetNamedConfigurations(facadeNamedConfigurations)

	namedConfiguration, err := facade.GetNamedConfigurationById("ID1")
	if err != nil {
		t.Errorf("unexpected error")
	}

	if !reflect.DeepEqual(namedConfiguration, facadeNamedConfigurations[0]) {
		t.Errorf("unexpected configuration")
	}
}

func TestNamedConfigurationsFacade_GetNamedConfigurationById_NamedConfigurationNotFound(t *testing.T) {
	facadeSetup()
	facade.SetNamedConfigurations(facadeNamedConfigurations)

	_, err := facade.GetNamedConfigurationById("OTHER-ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "named configuration with id OTHER-ID not found")
}

func TestCheckConfigurationNameValidity(t *testing.T) {
	facadeSetup()

	invalidNames := []string{
		"",
		"1start-whit-a-number",
		"#contains#special#chars",
		"UPPER-CASE-chars"}

	for _, invalidName := range invalidNames {
		err := facade.checkConfigurationNameValidity(invalidName)
		test.ExpectHttpError(t, err, http.StatusBadRequest, "configuration names must start with a lower case letter and contain only lower case letters a-z, digits 0-9, and hyphens '-'")
	}
}

type fakeNamedConfigurationObserver struct {
	ExpErrorOnUpdateState bool
}

func (fakeObserver fakeNamedConfigurationObserver) UpdateNamedConfigurations(oldNamedConfigurations []NamedConfiguration,
	newNamedConfigurations []NamedConfiguration) error {
	namedConfigurationBeforeUpdate = oldNamedConfigurations
	namedConfigurationAfterUpdate = newNamedConfigurations

	if fakeObserver.ExpErrorOnUpdateState {
		return http_error.NewInternalServerError(errors.New("unable to update named configurations"))
	}

	return nil
}
