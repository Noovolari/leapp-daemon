package gcp_iam_user_account_oauth

import (
	"leapp_daemon/domain/domain_gcp"
	"leapp_daemon/test"
	"net/http"
	"reflect"
	"testing"
)

var (
	sessionFacade        *GcpIamUserAccountOauthSessionsFacade
	sessionsBeforeUpdate []GcpIamUserAccountOauthSession
	sessionsAfterUpdate  []GcpIamUserAccountOauthSession
)

func facadeSetup() {
	sessionFacade = NewGcpIamUserAccountOauthSessionsFacade()
	sessionsBeforeUpdate = []GcpIamUserAccountOauthSession{}
	sessionsAfterUpdate = []GcpIamUserAccountOauthSession{}
}

func TestGcpIamUserSessionFacade_GetSessions(t *testing.T) {
	facadeSetup()

	newSessions := []GcpIamUserAccountOauthSession{{Id: "id"}}
	sessionFacade.gcpIamUserAccountOauthSessions = newSessions

	if !reflect.DeepEqual(sessionFacade.GetSessions(), newSessions) {
		t.Errorf("unexpected sessions")
	}
}

func TestGcpIamUserSessionFacade_SetSessions(t *testing.T) {
	facadeSetup()

	newSessions := []GcpIamUserAccountOauthSession{{Id: "id"}}
	sessionFacade.SetSessions(newSessions)

	if !reflect.DeepEqual(sessionFacade.gcpIamUserAccountOauthSessions, newSessions) {
		t.Errorf("unexpected sessions")
	}
}

func TestGcpIamUserSessionFacade_AddSession(t *testing.T) {
	facadeSetup()
	sessionFacade.Subscribe(fakeSessionObserver{})

	newSession := GcpIamUserAccountOauthSession{Id: "id"}
	sessionFacade.AddSession(newSession)

	if !reflect.DeepEqual(sessionsBeforeUpdate, []GcpIamUserAccountOauthSession{}) {
		t.Errorf("sessions were not empty")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []GcpIamUserAccountOauthSession{newSession}) {
		t.Errorf("unexpected session")
	}
}

func TestGcpIamUserSessionFacade_AddSession_alreadyExistentId(t *testing.T) {
	facadeSetup()
	sessionFacade.Subscribe(fakeSessionObserver{})

	newSession := GcpIamUserAccountOauthSession{Id: "ID"}
	sessionFacade.gcpIamUserAccountOauthSessions = []GcpIamUserAccountOauthSession{newSession}

	err := sessionFacade.AddSession(newSession)
	test.ExpectHttpError(t, err, http.StatusConflict, "a session with id ID is already present")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestGcpIamUserSessionFacade_AddSession_alreadyExistentName(t *testing.T) {
	facadeSetup()
	sessionFacade.Subscribe(fakeSessionObserver{})

	sessionFacade.gcpIamUserAccountOauthSessions = []GcpIamUserAccountOauthSession{{Id: "1", Name: "NAME"}}

	err := sessionFacade.AddSession(GcpIamUserAccountOauthSession{Id: "2", Name: "NAME"})
	test.ExpectHttpError(t, err, http.StatusConflict, "a session named NAME is already present")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestGcpIamUserSessionFacade_RemoveSession(t *testing.T) {
	facadeSetup()
	sessionFacade.Subscribe(fakeSessionObserver{})

	session1 := GcpIamUserAccountOauthSession{Id: "ID1"}
	session2 := GcpIamUserAccountOauthSession{Id: "ID2"}
	sessionFacade.gcpIamUserAccountOauthSessions = []GcpIamUserAccountOauthSession{session1, session2}

	sessionFacade.RemoveSession("ID1")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []GcpIamUserAccountOauthSession{session1, session2}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []GcpIamUserAccountOauthSession{session2}) {
		t.Errorf("sessions were not empty")
	}
}

func TestGcpIamUserSessionFacade_RemoveSession_notFound(t *testing.T) {
	facadeSetup()
	sessionFacade.Subscribe(fakeSessionObserver{})

	err := sessionFacade.RemoveSession("ID")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestGcpIamUserSessionFacade_StartSession(t *testing.T) {
	facadeSetup()
	sessionFacade.Subscribe(fakeSessionObserver{})

	newSession := GcpIamUserAccountOauthSession{Id: "ID", Status: domain_gcp.NotActive}
	sessionFacade.gcpIamUserAccountOauthSessions = []GcpIamUserAccountOauthSession{newSession}

	sessionFacade.StartSession("ID", "start-time")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []GcpIamUserAccountOauthSession{newSession}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []GcpIamUserAccountOauthSession{{Id: "ID", Status: domain_gcp.Active, StartTime: "start-time"}}) {
		t.Errorf("sessions were not updated")
	}
}

func TestGcpIamUserSessionFacade_StartSession_notFound(t *testing.T) {
	facadeSetup()
	sessionFacade.Subscribe(fakeSessionObserver{})

	err := sessionFacade.StartSession("ID", "start-time")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestGcpIamUserSessionFacade_StopSession(t *testing.T) {
	facadeSetup()
	sessionFacade.Subscribe(fakeSessionObserver{})

	newSession := GcpIamUserAccountOauthSession{Id: "ID", Status: domain_gcp.Active}
	sessionFacade.gcpIamUserAccountOauthSessions = []GcpIamUserAccountOauthSession{newSession}

	sessionFacade.StopSession("ID", "stop-time")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []GcpIamUserAccountOauthSession{newSession}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []GcpIamUserAccountOauthSession{{Id: "ID", Status: domain_gcp.NotActive, LastStopTime: "stop-time"}}) {
		t.Errorf("sessions were not updated")
	}
}

func TestGcpIamUserSessionFacade_StopSession_notFound(t *testing.T) {
	facadeSetup()
	sessionFacade.Subscribe(fakeSessionObserver{})

	err := sessionFacade.StopSession("ID", "stop-time")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

func TestGcpIamUserSessionFacade_EditSession(t *testing.T) {
	facadeSetup()
	sessionFacade.Subscribe(fakeSessionObserver{})

	session1 := GcpIamUserAccountOauthSession{Id: "ID1", Name: "Name1", ProjectName: "Project1"}
	session2 := GcpIamUserAccountOauthSession{Id: "ID2", Name: "Name2", ProjectName: "Project2"}
	sessionFacade.gcpIamUserAccountOauthSessions = []GcpIamUserAccountOauthSession{session1, session2}

	sessionFacade.EditSession("ID1", "NewName", "NewProject", "")

	if !reflect.DeepEqual(sessionsBeforeUpdate, []GcpIamUserAccountOauthSession{session1, session2}) {
		t.Errorf("unexpected session")
	}

	if !reflect.DeepEqual(sessionsAfterUpdate, []GcpIamUserAccountOauthSession{
		{Id: "ID1", Name: "NewName", ProjectName: "NewProject"}, session2}) {
		t.Errorf("sessions were not updated")
	}
}

func TestGcpIamUserSessionFacade_EditSession_DuplicateSessionNameAttempt(t *testing.T) {
	facadeSetup()
	sessionFacade.Subscribe(fakeSessionObserver{})

	session1 := GcpIamUserAccountOauthSession{Id: "ID1", Name: "Name1", ProjectName: "Project1"}
	session2 := GcpIamUserAccountOauthSession{Id: "ID2", Name: "Name2", ProjectName: "Project2"}
	sessionFacade.gcpIamUserAccountOauthSessions = []GcpIamUserAccountOauthSession{session1, session2}

	err := sessionFacade.EditSession("ID1", "Name2", "NewProject", "")
	test.ExpectHttpError(t, err, http.StatusConflict, "a session named Name2 is already present")

	err = sessionFacade.EditSession("ID2", "Name1", "NewProject", "")
	test.ExpectHttpError(t, err, http.StatusConflict, "a session named Name1 is already present")
}

func TestGcpIamUserSessionFacade_EditSession_notFound(t *testing.T) {
	facadeSetup()
	sessionFacade.Subscribe(fakeSessionObserver{})

	err := sessionFacade.EditSession("ID", "", "", "")
	test.ExpectHttpError(t, err, http.StatusNotFound, "session with id ID not found")

	if len(sessionsBeforeUpdate) > 0 || len(sessionsAfterUpdate) > 0 {
		t.Errorf("sessions was unexpectedly changed")
	}
}

type fakeSessionObserver struct {
}

func (observer fakeSessionObserver) UpdateGcpIamUserAccountOauthSessions(oldSessions []GcpIamUserAccountOauthSession, newSessions []GcpIamUserAccountOauthSession) {
	sessionsBeforeUpdate = oldSessions
	sessionsAfterUpdate = newSessions
}
