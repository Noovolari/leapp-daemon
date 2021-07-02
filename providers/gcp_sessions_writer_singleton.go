package providers

import (
	"leapp_daemon/use_case"
	"sync"
)

var gcpSessionsWriterSingleton *use_case.GcpSessionsWriter
var gcpSessionsWriterMutex sync.Mutex

func (prov *Providers) GetGcpSessionWriter() *use_case.GcpSessionsWriter {
	gcpSessionsWriterMutex.Lock()
	defer gcpSessionsWriterMutex.Unlock()

	if gcpSessionsWriterSingleton == nil {
		gcpSessionsWriterSingleton = &use_case.GcpSessionsWriter{
			ConfigurationRepository: prov.GetFileConfigurationRepository(),
		}
	}
	return gcpSessionsWriterSingleton
}
