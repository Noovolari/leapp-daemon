package providers

import (
	"leapp_daemon/domain/domain_aws/aws_iam_user"
	"sync"
)

var awsSessionsFacadeSingleton *aws_iam_user.AwsIamUserSessionsFacade
var awsSessionsFacadeLock sync.Mutex

func (prov *Providers) GetAwsIamUserSessionFacade() *aws_iam_user.AwsIamUserSessionsFacade {
	awsSessionsFacadeLock.Lock()
	defer awsSessionsFacadeLock.Unlock()

	if awsSessionsFacadeSingleton == nil {
		awsSessionsFacadeSingleton = aws_iam_user.NewAwsIamUserSessionsFacade()
	}
	return awsSessionsFacadeSingleton
}
