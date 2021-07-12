package providers

import (
	"leapp_daemon/adapter/gcp"
	"sync"
)

var gcpCredentialsTableSingleton *gcp.GcpCredentialsTable
var gcpCredentialsTableMutex sync.Mutex

func (prov *Providers) GetGcpCredentialsTable() *gcp.GcpCredentialsTable {
	gcpCredentialsTableMutex.Lock()
	defer gcpCredentialsTableMutex.Unlock()

	if gcpCredentialsTableSingleton == nil {
		gcpCredentialsTableSingleton = &gcp.GcpCredentialsTable{}
	}
	return gcpCredentialsTableSingleton
}
