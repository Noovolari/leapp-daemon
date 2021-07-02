package providers

import (
	"leapp_daemon/use_case"
	"sync"
)

var alibabaSessionsWriterSingleton *use_case.AlibabaSessionsWriter
var alibabaSessionsWriterMutex sync.Mutex

func (prov *Providers) GetAlibabaSessionWriter() *use_case.AlibabaSessionsWriter {
	alibabaSessionsWriterMutex.Lock()
	defer alibabaSessionsWriterMutex.Unlock()

	if alibabaSessionsWriterSingleton == nil {
		alibabaSessionsWriterSingleton = &use_case.AlibabaSessionsWriter{
			ConfigurationRepository: prov.GetFileConfigurationRepository(),
		}
	}
	return alibabaSessionsWriterSingleton
}