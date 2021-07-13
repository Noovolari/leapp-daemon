package alibaba_ram_role_federated

import (
	"leapp_daemon/domain/domain_alibaba"
	"leapp_daemon/test"
	"net/http"
	"reflect"
	"testing"
)

var (
	facade               *AlibabaRamRoleFederatedSessionsFacade
	sessionsBeforeUpdate []AlibabaRamRoleFederatedSession
	sessionsAfterUpdate  []AlibabaRamRoleFederatedSession
)

func facadeSetup() {
	facade = NewAlibabaRamRoleFederatedSessionsFacade()
	sessionsBeforeUpdate = []AlibabaRamRoleFederatedSession{}
	sessionsAfterUpdate = []AlibabaRamRoleFederatedSession{}
}

func TestAlibabaRamRoleFederatedSessionsFacade_GetSessions(t *testing.T) {
	facadeSetup()

	newSessionessions := []AlibabaRamRoleFederatedSession{{Id: "id"}}
	facade.alibabaRamRoleFederatedSessions = newSessionessions

	if !reflect.DeepEqual(facade.GetSessions(), newSessionessions) {
		t.Errorf("unexpected sessions")
	}
}

func TestAlibabaRamRoleFederatedSessionsFacade_SetSessions(t *testing.T) {
	facadeSetup()

	newSessionessions := []AlibabaRamRoleFederatedSession{{Id: "id"}}
	facade.SetSessions(newSessionessions)

	if !reflect.DeepEqual(facade.alibabaRamRoleFederatedSessions, newSessionessions) {
		t.Errorf("unexpected sessions")
	}
}

func TestAlibabaRamRoleFederatedSessionsFacade_AddSession(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	newSession := AlibabaRamRoleFederatedSession{Id: "id"}
	facade.AddSession(newSession)

	if !reflect.DeepEqual(sessionsBeforeUpdate, []AlibabaRamRoleFederatedSession{}) {
		t.Errorf("sessions were not empty")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []AlibabaRamRoleFederatedSession{newSession}) {
		t.Errorf("unexpected session")
	}
}

func TestAlibabaRamRoleFederatedSessionsFacade_AddSession_alreadyExistentId(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	newSessionession := AlibabaRamRoleFederatedSession{Id: "ID"}
	facade.alibabaRamRoleFederatedSessions = []AlibabaRamRoleFederatedSession{newSessionession}

	err := facade.AddSession(newSessionession)
	test.ExpectHttpError(t, err, http.StatusConflict, "a session with id ID is already present")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestAlibabaRamRoleFederatedSessionsFacade_AddSession_alreadyExistentName(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	facade.alibabaRamRoleFederatedSessions = []AlibabaRamRoleFederatedSession{{Id: "1", Name: "NAME"}}

	err := facade.AddSession(AlibabaRamRoleFederatedSession{Id: "2", Name: "NAME"})
	test.ExpectHttpError(t, err, http.StatusConflict, "a session named NAME is already present")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestAlibabaRamRoleFederatedSessionsFacade_RemoveSession(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	session1 := AlibabaRamRoleFederatedSession{Id: "ID1"}
	session2 := AlibabaRamRoleFederatedSession{Id: "ID2"}
	facade.alibabaRamRoleFederatedSessions = []AlibabaRamRoleFederatedSession{session1, session2}

	facade.RemoveSession("ID1")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []AlibabaRamRoleFederatedSession{session1, session2}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []AlibabaRamRoleFederatedSession{session2}) {
		t.Errorf("sessions were not empty")
	}
}

func TestAlibabaRamRoleFederatedSessionsFacade_RemoveSession_notFound(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	err := facade.RemoveSession("ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestAlibabaRamRoleFederatedSessionsFacade_StartingSession(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	newSession := AlibabaRamRoleFederatedSession{Id: "ID", Status: domain_alibaba.NotActive}
	facade.alibabaRamRoleFederatedSessions = []AlibabaRamRoleFederatedSession{newSession}

	facade.StartingSession("ID")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []AlibabaRamRoleFederatedSession{newSession}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []AlibabaRamRoleFederatedSession{{Id: "ID", Status: domain_alibaba.Pending}}) {
		t.Errorf("sessions were not updated")
	}
}

func TestAlibabaRamRoleFederatedSessionsFacade_StartingSession_notFound(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	err := facade.StartingSession("ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestAlibabaRamRoleFederatedSessionsFacade_StartSession(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	newSession := AlibabaRamRoleFederatedSession{Id: "ID", Status: domain_alibaba.NotActive}
	facade.alibabaRamRoleFederatedSessions = []AlibabaRamRoleFederatedSession{newSession}

	facade.StartSession("ID")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []AlibabaRamRoleFederatedSession{newSession}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []AlibabaRamRoleFederatedSession{{Id: "ID", Status: domain_alibaba.Active}}) {
		t.Errorf("sessions were not updated")
	}
}

func TestAlibabaRamRoleFederatedSessionsFacade_StartSession_notFound(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	err := facade.StartSession("ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestAlibabaRamRoleFederatedSessionsFacade_StopSession(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	newSessionession := AlibabaRamRoleFederatedSession{Id: "ID", Status: domain_alibaba.Active}
	facade.alibabaRamRoleFederatedSessions = []AlibabaRamRoleFederatedSession{newSessionession}

	facade.StopSession("ID")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []AlibabaRamRoleFederatedSession{newSessionession}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []AlibabaRamRoleFederatedSession{{Id: "ID", Status: domain_alibaba.NotActive}}) {
		t.Errorf("sessions were not updated")
	}
}

func TestAlibabaRamRoleFederatedSessionsFacade_StopSession_notFound(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	err := facade.StopSession("ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestAlibabaRamRoleFederatedSessionsFacade_EditSession(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	session1 := AlibabaRamRoleFederatedSession{Id: "ID1", Name: "Name1", Region: "region", NamedProfileId: "ProfileId1"}
	session2 := AlibabaRamRoleFederatedSession{Id: "ID2", Name: "Name2", Region: "region2", NamedProfileId: "ProfileId2"}
	facade.alibabaRamRoleFederatedSessions = []AlibabaRamRoleFederatedSession{session1, session2}

	facade.EditSession("ID1", "newName", "newRoleName", "newRoleArn", "newIdpArn", "newRegion", "newSsoUrl", "newProfileId")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []AlibabaRamRoleFederatedSession{session1, session2}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []AlibabaRamRoleFederatedSession{
		{Id: "ID1", Name: "newName", RoleName: "newRoleName", RoleArn: "newRoleArn", IdpArn: "newIdpArn", Region: "newRegion", SsoUrl: "newSsoUrl", NamedProfileId: "newProfileId", Status: domain_alibaba.NotActive}, session2}) {
		t.Errorf("sessions were not updated")
	}
}

func TestAlibabaRamRoleFederatedSessionsFacade_EditSession_DuplicateSessionNameAttempt(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	session1 := AlibabaRamRoleFederatedSession{Id: "ID1", Name: "Name1", Region: "region", NamedProfileId: "ProfileId1"}
	session2 := AlibabaRamRoleFederatedSession{Id: "ID2", Name: "Name2", Region: "region2", NamedProfileId: "ProfileId2"}
	facade.alibabaRamRoleFederatedSessions = []AlibabaRamRoleFederatedSession{session1, session2}

	err := facade.EditSession("ID1", "Name2", "newRoleName", "newRoleArn", "newIdpArn", "newRegion", "newSsoUrl", "newProfileId")
	test.ExpectHttpError(t, err, http.StatusConflict, "a session named Name2 is already present")

	err = facade.EditSession("ID2", "Name1", "newRoleName", "newRoleArn", "newIdpArn", "newRegion", "newSsoUrl", "newProfileId")
	test.ExpectHttpError(t, err, http.StatusConflict, "a session named Name1 is already present")
}

func TestAlibabaRamRoleFederatedSessionsFacade_EditSession_notFound(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	err := facade.EditSession("ID", "Name2", "newRoleName", "newRoleArn", "newIdpArn", "newRegion", "newSsoUrl", "newProfileId")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

type fakeSessionsObserver struct {
}

func (observer fakeSessionsObserver) UpdateAlibabaRamRoleFederatedSessions(oldSession []AlibabaRamRoleFederatedSession, newSessions []AlibabaRamRoleFederatedSession) error {
	sessionsBeforeUpdate = oldSession
	sessionsAfterUpdate = newSessions

	return nil
}
