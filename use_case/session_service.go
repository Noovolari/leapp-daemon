package use_case

import (
	"leapp_daemon/domain/named_profile"
)

func ListAllSessions(gcpIamUserAccountOauthSessionFacade GcpIamUserAccountOauthSessionsFacade) (*map[string]interface{}, error) {

	sessions := gcpIamUserAccountOauthSessionFacade.GetSessions()

	return &map[string]interface{}{
		"GcpSessions": sessions,
	}, nil
}

func ListAllNamedProfiles(namedProfileFacade NamedProfilesFacade) ([]named_profile.NamedProfile, error) {
	return namedProfileFacade.GetNamedProfiles(), nil
}
