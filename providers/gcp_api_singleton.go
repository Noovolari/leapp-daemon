package providers

import (
	"leapp_daemon/interface/gcp"
	"sync"
)

var gcpApiSingleton *gcp.GcpApi
var gcpApiMutex sync.Mutex

func (prov *Providers) GetGcpApi() *gcp.GcpApi {
	gcpApiMutex.Lock()
	defer gcpApiMutex.Unlock()

	if gcpApiSingleton == nil {
		gcpApiSingleton = &gcp.GcpApi{}
	}
	return gcpApiSingleton
}
