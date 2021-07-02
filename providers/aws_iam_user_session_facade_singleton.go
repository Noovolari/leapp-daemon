package providers

import (
	"leapp_daemon/domain/session"
	"sync"
)

var awsSessionsFacadeSingleton *session.AwsIamUserSessionsFacade
var awsSessionsFacadeLock sync.Mutex

func (prov *Providers) GetAwsIamUserSessionFacade() *session.AwsIamUserSessionsFacade {
	awsSessionsFacadeLock.Lock()
	defer awsSessionsFacadeLock.Unlock()

	if awsSessionsFacadeSingleton == nil {
		awsSessionsFacadeSingleton = session.NewAwsIamUserSessionsFacade()
	}
	return awsSessionsFacadeSingleton
}
