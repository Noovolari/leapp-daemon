package alibaba_ram_user

import (
	"leapp_daemon/domain/domain_alibaba"
	"leapp_daemon/test"
	"net/http"
	"reflect"
	"testing"
)

var (
	facade               *AlibabaRamUserSessionsFacade
	sessionsBeforeUpdate []AlibabaRamUserSession
	sessionsAfterUpdate  []AlibabaRamUserSession
)

func facadeSetup() {
	facade = NewAlibabaRamUserSessionsFacade()
	sessionsBeforeUpdate = []AlibabaRamUserSession{}
	sessionsAfterUpdate = []AlibabaRamUserSession{}
}

func TestAlibabaRamUserSessionsFacade_GetSessions(t *testing.T) {
	facadeSetup()

	newSessionessions := []AlibabaRamUserSession{{Id: "id"}}
	facade.alibabaRamUserSessions = newSessionessions

	if !reflect.DeepEqual(facade.GetSessions(), newSessionessions) {
		t.Errorf("unexpected sessions")
	}
}

func TestAlibabaRamUserSessionsFacade_SetSessions(t *testing.T) {
	facadeSetup()

	newSessionessions := []AlibabaRamUserSession{{Id: "id"}}
	facade.SetSessions(newSessionessions)

	if !reflect.DeepEqual(facade.alibabaRamUserSessions, newSessionessions) {
		t.Errorf("unexpected sessions")
	}
}

func TestAlibabaRamUserSessionsFacade_AddSession(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	newSession := AlibabaRamUserSession{Id: "id"}
	facade.AddSession(newSession)

	if !reflect.DeepEqual(sessionsBeforeUpdate, []AlibabaRamUserSession{}) {
		t.Errorf("sessions were not empty")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []AlibabaRamUserSession{newSession}) {
		t.Errorf("unexpected session")
	}
}

func TestAlibabaRamUserSessionsFacade_AddSession_alreadyExistentId(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	newSessionession := AlibabaRamUserSession{Id: "ID"}
	facade.alibabaRamUserSessions = []AlibabaRamUserSession{newSessionession}

	err := facade.AddSession(newSessionession)
	test.ExpectHttpError(t, err, http.StatusConflict, "a session with id ID is already present")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestAlibabaRamUserSessionsFacade_AddSession_alreadyExistentName(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	facade.alibabaRamUserSessions = []AlibabaRamUserSession{{Id: "1", Name: "NAME"}}

	err := facade.AddSession(AlibabaRamUserSession{Id: "2", Name: "NAME"})
	test.ExpectHttpError(t, err, http.StatusConflict, "a session named NAME is already present")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestAlibabaRamUserSessionsFacade_RemoveSession(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	session1 := AlibabaRamUserSession{Id: "ID1"}
	session2 := AlibabaRamUserSession{Id: "ID2"}
	facade.alibabaRamUserSessions = []AlibabaRamUserSession{session1, session2}

	facade.RemoveSession("ID1")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []AlibabaRamUserSession{session1, session2}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []AlibabaRamUserSession{session2}) {
		t.Errorf("sessions were not empty")
	}
}

func TestAlibabaRamUserSessionsFacade_RemoveSession_notFound(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	err := facade.RemoveSession("ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestAlibabaRamUserSessionsFacade_StartingSession(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	newSession := AlibabaRamUserSession{Id: "ID", Status: domain_alibaba.NotActive}
	facade.alibabaRamUserSessions = []AlibabaRamUserSession{newSession}

	facade.StartingSession("ID")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []AlibabaRamUserSession{newSession}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []AlibabaRamUserSession{{Id: "ID", Status: domain_alibaba.Pending}}) {
		t.Errorf("sessions were not updated")
	}
}

func TestAlibabaRamUserSessionsFacade_StartingSession_notFound(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	err := facade.StartingSession("ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestAlibabaRamUserSessionsFacade_StartSession(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	newSession := AlibabaRamUserSession{Id: "ID", Status: domain_alibaba.NotActive}
	facade.alibabaRamUserSessions = []AlibabaRamUserSession{newSession}

	facade.StartSession("ID")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []AlibabaRamUserSession{newSession}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []AlibabaRamUserSession{{Id: "ID", Status: domain_alibaba.Active}}) {
		t.Errorf("sessions were not updated")
	}
}

func TestAlibabaRamUserSessionsFacade_StartSession_notFound(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	err := facade.StartSession("ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestAlibabaRamUserSessionsFacade_StopSession(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	newSessionession := AlibabaRamUserSession{Id: "ID", Status: domain_alibaba.Active}
	facade.alibabaRamUserSessions = []AlibabaRamUserSession{newSessionession}

	facade.StopSession("ID")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []AlibabaRamUserSession{newSessionession}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []AlibabaRamUserSession{{Id: "ID", Status: domain_alibaba.NotActive}}) {
		t.Errorf("sessions were not updated")
	}
}

func TestAlibabaRamUserSessionsFacade_StopSession_notFound(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	err := facade.StopSession("ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestAlibabaRamUserSessionsFacade_EditSession(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	session1 := AlibabaRamUserSession{Id: "ID1", Name: "Name1", Region: "region", NamedProfileId: "ProfileId1"}
	session2 := AlibabaRamUserSession{Id: "ID2", Name: "Name2", Region: "region2", NamedProfileId: "ProfileId2"}
	facade.alibabaRamUserSessions = []AlibabaRamUserSession{session1, session2}

	facade.EditSession("ID1", "newName", "newRegion", "newProfileId")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []AlibabaRamUserSession{session1, session2}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []AlibabaRamUserSession{
		{Id: "ID1", Name: "newName", Region: "newRegion", NamedProfileId: "newProfileId", Status: domain_alibaba.NotActive}, session2}) {
		t.Errorf("sessions were not updated")
	}
}

func TestAlibabaRamUserSessionsFacade_EditSession_DuplicateSessionNameAttempt(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	session1 := AlibabaRamUserSession{Id: "ID1", Name: "Name1", Region: "region", NamedProfileId: "ProfileId1"}
	session2 := AlibabaRamUserSession{Id: "ID2", Name: "Name2", Region: "region2", NamedProfileId: "ProfileId2"}
	facade.alibabaRamUserSessions = []AlibabaRamUserSession{session1, session2}

	err := facade.EditSession("ID1", "Name2", "newRegion", "newProfileId")
	test.ExpectHttpError(t, err, http.StatusConflict, "a session named Name2 is already present")

	err = facade.EditSession("ID2", "Name1", "newRegion", "newProfileId")
	test.ExpectHttpError(t, err, http.StatusConflict, "a session named Name1 is already present")
}

func TestAlibabaRamUserSessionsFacade_EditSession_notFound(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	err := facade.EditSession("ID", "Name2", "NewRegion", "newProfileId")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

type fakeSessionsObserver struct {
}

func (observer fakeSessionsObserver) UpdateAlibabaRamUserSessions(oldSession []AlibabaRamUserSession, newSessions []AlibabaRamUserSession) error {
	sessionsBeforeUpdate = oldSession
	sessionsAfterUpdate = newSessions

	return nil
}
