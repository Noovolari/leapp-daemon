package providers

import (
	"leapp_daemon/interface/gcp"
	"sync"
)

var gcpAccessTokensTableSingleton *gcp.GcpAccessTokensTable
var gcpAccessTokensTableMutex sync.Mutex

func (prov *Providers) GetGcpAccessTokensTable() *gcp.GcpAccessTokensTable {
	gcpAccessTokensTableMutex.Lock()
	defer gcpAccessTokensTableMutex.Unlock()

	if gcpAccessTokensTableSingleton == nil {
		gcpAccessTokensTableSingleton = &gcp.GcpAccessTokensTable{}
	}
	return gcpAccessTokensTableSingleton
}
