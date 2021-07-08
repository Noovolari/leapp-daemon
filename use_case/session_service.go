package use_case

import (
	"leapp_daemon/domain/domain_aws/named_profile"
)

//TODO remove this class and enventually move this methods under the proper actions class
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
