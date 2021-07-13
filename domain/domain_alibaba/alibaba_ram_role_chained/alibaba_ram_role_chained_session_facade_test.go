package alibaba_ram_role_chained

import (
	"leapp_daemon/domain/domain_alibaba"
	"leapp_daemon/test"
	"net/http"
	"reflect"
	"testing"
)

var (
	facade               *AlibabaRamRoleChainedSessionsFacade
	sessionsBeforeUpdate []AlibabaRamRoleChainedSession
	sessionsAfterUpdate  []AlibabaRamRoleChainedSession
)

func facadeSetup() {
	facade = NewAlibabaRamRoleChainedSessionsFacade()
	sessionsBeforeUpdate = []AlibabaRamRoleChainedSession{}
	sessionsAfterUpdate = []AlibabaRamRoleChainedSession{}
}

func TestAlibabaRamRoleChainedSessionsFacade_GetSessions(t *testing.T) {
	facadeSetup()

	newSessionessions := []AlibabaRamRoleChainedSession{{Id: "id"}}
	facade.alibabaRamRoleChainedSessions = newSessionessions

	if !reflect.DeepEqual(facade.GetSessions(), newSessionessions) {
		t.Errorf("unexpected sessions")
	}
}

func TestAlibabaRamRoleChainedSessionsFacade_SetSessions(t *testing.T) {
	facadeSetup()

	newSessionessions := []AlibabaRamRoleChainedSession{{Id: "id"}}
	facade.SetSessions(newSessionessions)

	if !reflect.DeepEqual(facade.alibabaRamRoleChainedSessions, newSessionessions) {
		t.Errorf("unexpected sessions")
	}
}

func TestAlibabaRamRoleChainedSessionsFacade_AddSession(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	newSession := AlibabaRamRoleChainedSession{Id: "id"}
	facade.AddSession(newSession)

	if !reflect.DeepEqual(sessionsBeforeUpdate, []AlibabaRamRoleChainedSession{}) {
		t.Errorf("sessions were not empty")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []AlibabaRamRoleChainedSession{newSession}) {
		t.Errorf("unexpected session")
	}
}

func TestAlibabaRamRoleChainedSessionsFacade_AddSession_alreadyExistentId(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	newSessionession := AlibabaRamRoleChainedSession{Id: "ID"}
	facade.alibabaRamRoleChainedSessions = []AlibabaRamRoleChainedSession{newSessionession}

	err := facade.AddSession(newSessionession)
	test.ExpectHttpError(t, err, http.StatusConflict, "a session with id ID is already present")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestAlibabaRamRoleChainedSessionsFacade_AddSession_alreadyExistentName(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	facade.alibabaRamRoleChainedSessions = []AlibabaRamRoleChainedSession{{Id: "1", Name: "NAME"}}

	err := facade.AddSession(AlibabaRamRoleChainedSession{Id: "2", Name: "NAME"})
	test.ExpectHttpError(t, err, http.StatusConflict, "a session named NAME is already present")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestAlibabaRamRoleChainedSessionsFacade_RemoveSession(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	session1 := AlibabaRamRoleChainedSession{Id: "ID1"}
	session2 := AlibabaRamRoleChainedSession{Id: "ID2"}
	facade.alibabaRamRoleChainedSessions = []AlibabaRamRoleChainedSession{session1, session2}

	facade.RemoveSession("ID1")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []AlibabaRamRoleChainedSession{session1, session2}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []AlibabaRamRoleChainedSession{session2}) {
		t.Errorf("sessions were not empty")
	}
}

func TestAlibabaRamRoleChainedSessionsFacade_RemoveSession_notFound(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	err := facade.RemoveSession("ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestAlibabaRamRoleChainedSessionsFacade_StartingSession(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	newSession := AlibabaRamRoleChainedSession{Id: "ID", Status: domain_alibaba.NotActive}
	facade.alibabaRamRoleChainedSessions = []AlibabaRamRoleChainedSession{newSession}

	facade.StartingSession("ID")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []AlibabaRamRoleChainedSession{newSession}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []AlibabaRamRoleChainedSession{{Id: "ID", Status: domain_alibaba.Pending}}) {
		t.Errorf("sessions were not updated")
	}
}

func TestAlibabaRamRoleChainedSessionsFacade_StartingSession_notFound(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	err := facade.StartingSession("ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestAlibabaRamRoleChainedSessionsFacade_StartSession(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	newSession := AlibabaRamRoleChainedSession{Id: "ID", Status: domain_alibaba.NotActive}
	facade.alibabaRamRoleChainedSessions = []AlibabaRamRoleChainedSession{newSession}

	facade.StartSession("ID")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []AlibabaRamRoleChainedSession{newSession}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []AlibabaRamRoleChainedSession{{Id: "ID", Status: domain_alibaba.Active}}) {
		t.Errorf("sessions were not updated")
	}
}

func TestAlibabaRamRoleChainedSessionsFacade_StartSession_notFound(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	err := facade.StartSession("ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestAlibabaRamRoleChainedSessionsFacade_StopSession(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	newSessionession := AlibabaRamRoleChainedSession{Id: "ID", Status: domain_alibaba.Active}
	facade.alibabaRamRoleChainedSessions = []AlibabaRamRoleChainedSession{newSessionession}

	facade.StopSession("ID")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []AlibabaRamRoleChainedSession{newSessionession}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []AlibabaRamRoleChainedSession{{Id: "ID", Status: domain_alibaba.NotActive}}) {
		t.Errorf("sessions were not updated")
	}
}

func TestAlibabaRamRoleChainedSessionsFacade_StopSession_notFound(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	err := facade.StopSession("ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestAlibabaRamRoleChainedSessionsFacade_EditSession(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	session1 := AlibabaRamRoleChainedSession{Id: "ID1", Name: "Name1", Region: "region", NamedProfileId: "ProfileId1"}
	session2 := AlibabaRamRoleChainedSession{Id: "ID2", Name: "Name2", Region: "region2", NamedProfileId: "ProfileId2"}
	facade.alibabaRamRoleChainedSessions = []AlibabaRamRoleChainedSession{session1, session2}

	facade.EditSession("ID1", "newName", "newRoleName", "newAccountNumber", "newRoleArn", "newRegion", "newParentId", "newParentType", "newProfileId")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []AlibabaRamRoleChainedSession{session1, session2}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []AlibabaRamRoleChainedSession{
		{Id: "ID1", Name: "newName", RoleName: "newRoleName", RoleArn: "newRoleArn", AccountNumber: "newAccountNumber", Region: "newRegion", ParentId: "newParentId", ParentType: "newParentType", NamedProfileId: "newProfileId", Status: domain_alibaba.NotActive}, session2}) {
		t.Errorf("sessions were not updated")
	}
}

func TestAlibabaRamRoleChainedSessionsFacade_EditSession_DuplicateSessionNameAttempt(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	session1 := AlibabaRamRoleChainedSession{Id: "ID1", Name: "Name1", Region: "region", NamedProfileId: "ProfileId1"}
	session2 := AlibabaRamRoleChainedSession{Id: "ID2", Name: "Name2", Region: "region2", NamedProfileId: "ProfileId2"}
	facade.alibabaRamRoleChainedSessions = []AlibabaRamRoleChainedSession{session1, session2}

	err := facade.EditSession("ID1", "Name2", "newRoleName", "newAccountNumber", "newRoleArn", "newRegion", "newParentId", "newParentType", "newProfileId")
	test.ExpectHttpError(t, err, http.StatusConflict, "a session named Name2 is already present")

	err = facade.EditSession("ID2", "Name1", "newRoleName", "newAccountNumber", "newRoleArn", "newRegion", "newParentId", "newParentType", "newProfileId")
	test.ExpectHttpError(t, err, http.StatusConflict, "a session named Name1 is already present")
}

func TestAlibabaRamRoleChainedSessionsFacade_EditSession_notFound(t *testing.T) {
	facadeSetup()
	facade.Subscribe(fakeSessionsObserver{})

	err := facade.EditSession("ID", "Name2", "newRoleName", "newAccountNumber", "newRoleArn", "newRegion", "newParentId", "newParentType", "newProfileId")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

type fakeSessionsObserver struct {
}

func (observer fakeSessionsObserver) UpdateAlibabaRamRoleChainedSessions(oldSession []AlibabaRamRoleChainedSession, newSessions []AlibabaRamRoleChainedSession) error {
	sessionsBeforeUpdate = oldSession
	sessionsAfterUpdate = newSessions

	return nil
}
