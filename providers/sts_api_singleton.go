package providers

import (
	"leapp_daemon/adapter/aws"
	"sync"
)

var stsApiSingleton *aws.StsApi
var stsApiMutex sync.Mutex

func (prov *Providers) GetStsApi() *aws.StsApi {
	stsApiMutex.Lock()
	defer stsApiMutex.Unlock()

	if stsApiSingleton == nil {
		stsApiSingleton = &aws.StsApi{}
	}
	return stsApiSingleton
}
