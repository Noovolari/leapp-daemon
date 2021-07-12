package providers

import (
	"leapp_daemon/domain/domain_gcp/gcp_iam_user_account_oauth"
	"sync"
)

var gcpSessionsFacadeSingleton *gcp_iam_user_account_oauth.GcpIamUserAccountOauthSessionsFacade
var gcpSessionsFacadeLock sync.Mutex

func (prov *Providers) GetGcpIamUserAccountOauthSessionFacade() *gcp_iam_user_account_oauth.GcpIamUserAccountOauthSessionsFacade {
	gcpSessionsFacadeLock.Lock()
	defer gcpSessionsFacadeLock.Unlock()

	if gcpSessionsFacadeSingleton == nil {
		gcpSessionsFacadeSingleton = gcp_iam_user_account_oauth.NewGcpIamUserAccountOauthSessionsFacade()
	}
	return gcpSessionsFacadeSingleton
}
