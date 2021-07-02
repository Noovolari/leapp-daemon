package providers

import (
	"leapp_daemon/use_case"
	"sync"
)

var awsIamUserSessionActionsSingleton *use_case.AwsIamUserSessionActions
var awsIamUserSessionActionsMutex sync.Mutex

func (prov *Providers) GetAwsIamUserSessionActions() *use_case.AwsIamUserSessionActions {
	awsIamUserSessionActionsMutex.Lock()
	defer awsIamUserSessionActionsMutex.Unlock()

	if awsIamUserSessionActionsSingleton == nil {
		awsIamUserSessionActionsSingleton = &use_case.AwsIamUserSessionActions{
			NamedProfilesActions:     prov.GetNamedProfilesActions(),
			Environment:              prov.GetEnvironment(),
			Keychain:                 prov.GetKeychain(),
			AwsIamUserSessionsFacade: prov.GetAwsIamUserSessionFacade(),
		}
	}
	return awsIamUserSessionActionsSingleton
}
