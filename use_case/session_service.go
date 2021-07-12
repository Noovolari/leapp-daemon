package use_case

import (
	"leapp_daemon/domain/named_profile"
)

func ListAllSessions(gcpIamUserAccountOauthSessionFacade GcpIamUserAccountOauthSessionsFacade,
	awsIamUserSessionFacade AwsIamUserSessionsFacade) (*map[string]interface{}, error) {

	return &map[string]interface{}{
		"AwsSessions": awsIamUserSessionFacade.GetSessions(),
		"GcpSessions": gcpIamUserAccountOauthSessionFacade.GetSessions(),
	}, nil
}

func ListAllNamedProfiles(namedProfileFacade NamedProfilesFacade) ([]named_profile.NamedProfile, error) {
	return namedProfileFacade.GetNamedProfiles(), nil
}
