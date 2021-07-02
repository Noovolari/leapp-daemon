package session

import (
	"leapp_daemon/test"
	"net/http"
	"reflect"
	"testing"
)

var (
	facade               *GcpIamUserAccountOauthSessionsFacade
	sessionsBeforeUpdate []GcpIamUserAccountOauthSession
	sessionsAfterUpdate  []GcpIamUserAccountOauthSession
)

func gcpIamUserAccountOauthSessionFacadeSetup() {
	facade = NewGcpIamUserAccountOauthSessionsFacade()
	sessionsBeforeUpdate = []GcpIamUserAccountOauthSession{}
	sessionsAfterUpdate = []GcpIamUserAccountOauthSession{}
}

func TestGetSessions(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()

	newSessions := []GcpIamUserAccountOauthSession{{Id: "id"}}
	facade.gcpIamUserAccountOauthSessions = newSessions

	if !reflect.DeepEqual(facade.GetSessions(), newSessions) {
		t.Errorf("unexpected sessions")
	}
}

func TestSetSessions(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()

	newSessions := []GcpIamUserAccountOauthSession{{Id: "id"}}
	facade.SetSessions(newSessions)

	if !reflect.DeepEqual(facade.gcpIamUserAccountOauthSessions, newSessions) {
		t.Errorf("unexpected sessions")
	}
}

func TestAddSession(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()
	facade.Subscribe(FakeObserver{})

	newSession := GcpIamUserAccountOauthSession{Id: "id"}
	facade.AddSession(newSession)

	if !reflect.DeepEqual(sessionsBeforeUpdate, []GcpIamUserAccountOauthSession{}) {
		t.Errorf("sessions were not empty")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []GcpIamUserAccountOauthSession{newSession}) {
		t.Errorf("unexpected session")
	}
}

func TestAddSession_alreadyExistentId(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()
	facade.Subscribe(FakeObserver{})

	newSession := GcpIamUserAccountOauthSession{Id: "ID"}
	facade.gcpIamUserAccountOauthSessions = []GcpIamUserAccountOauthSession{newSession}

	err := facade.AddSession(newSession)
	test.ExpectHttpError(t, err, http.StatusConflict, "a session with id ID is already present")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestAddSession_alreadyExistentName(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()
	facade.Subscribe(FakeObserver{})

	facade.gcpIamUserAccountOauthSessions = []GcpIamUserAccountOauthSession{{Id: "1", Name: "NAME"}}

	err := facade.AddSession(GcpIamUserAccountOauthSession{Id: "2", Name: "NAME"})
	test.ExpectHttpError(t, err, http.StatusConflict, "a session named NAME is already present")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestRemoveSession(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()
	facade.Subscribe(FakeObserver{})

	session1 := GcpIamUserAccountOauthSession{Id: "ID1"}
	session2 := GcpIamUserAccountOauthSession{Id: "ID2"}
	facade.gcpIamUserAccountOauthSessions = []GcpIamUserAccountOauthSession{session1, session2}

	facade.RemoveSession("ID1")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []GcpIamUserAccountOauthSession{session1, session2}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []GcpIamUserAccountOauthSession{session2}) {
		t.Errorf("sessions were not empty")
	}
}

func TestRemoveSession_notFound(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()
	facade.Subscribe(FakeObserver{})

	err := facade.RemoveSession("ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestStartSession(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()
	facade.Subscribe(FakeObserver{})

	newSession := GcpIamUserAccountOauthSession{Id: "ID", Status: NotActive}
	facade.gcpIamUserAccountOauthSessions = []GcpIamUserAccountOauthSession{newSession}

	facade.StartSession("ID", "start-time")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []GcpIamUserAccountOauthSession{newSession}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []GcpIamUserAccountOauthSession{{Id: "ID", Status: Active, StartTime: "start-time"}}) {
		t.Errorf("sessions were not updated")
	}
}

func TestStartSession_notFound(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()
	facade.Subscribe(FakeObserver{})

	err := facade.StartSession("ID", "start-time")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestStopSession(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()
	facade.Subscribe(FakeObserver{})

	newSession := GcpIamUserAccountOauthSession{Id: "ID", Status: Active}
	facade.gcpIamUserAccountOauthSessions = []GcpIamUserAccountOauthSession{newSession}

	facade.StopSession("ID", "stop-time")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []GcpIamUserAccountOauthSession{newSession}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []GcpIamUserAccountOauthSession{{Id: "ID", Status: NotActive, LastStopTime: "stop-time"}}) {
		t.Errorf("sessions were not updated")
	}
}

func TestStopSession_notFound(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()
	facade.Subscribe(FakeObserver{})

	err := facade.StopSession("ID", "stop-time")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestEditSession(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()
	facade.Subscribe(FakeObserver{})

	session1 := GcpIamUserAccountOauthSession{Id: "ID1", Name: "Name1", ProjectName: "Project1"}
	session2 := GcpIamUserAccountOauthSession{Id: "ID2", Name: "Name2", ProjectName: "Project2"}
	facade.gcpIamUserAccountOauthSessions = []GcpIamUserAccountOauthSession{session1, session2}

	facade.EditSession("ID1", "NewName", "NewProject")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []GcpIamUserAccountOauthSession{session1, session2}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []GcpIamUserAccountOauthSession{
		{Id: "ID1", Name: "NewName", ProjectName: "NewProject"}, session2}) {
		t.Errorf("sessions were not updated")
	}
}

func TestEditSession_DuplicateSessionNameAttempt(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()
	facade.Subscribe(FakeObserver{})

	session1 := GcpIamUserAccountOauthSession{Id: "ID1", Name: "Name1", ProjectName: "Project1"}
	session2 := GcpIamUserAccountOauthSession{Id: "ID2", Name: "Name2", ProjectName: "Project2"}
	facade.gcpIamUserAccountOauthSessions = []GcpIamUserAccountOauthSession{session1, session2}

	err := facade.EditSession("ID1", "Name2", "NewProject")
	test.ExpectHttpError(t, err, http.StatusConflict, "a session named Name2 is already present")

	err = facade.EditSession("ID2", "Name1", "NewProject")
	test.ExpectHttpError(t, err, http.StatusConflict, "a session named Name1 is already present")
}

func TestEditSession_notFound(t *testing.T) {
	gcpIamUserAccountOauthSessionFacadeSetup()
	facade.Subscribe(FakeObserver{})

	err := facade.EditSession("ID", "", "")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

type FakeObserver struct {
}

func (observer FakeObserver) UpdateGcpIamUserAccountOauthSessions(oldSessions []GcpIamUserAccountOauthSession, newSessions []GcpIamUserAccountOauthSession) {
	sessionsBeforeUpdate = oldSessions
	sessionsAfterUpdate = newSessions
}
