package providers

import (
	"leapp_daemon/use_case"
	"sync"
)

var namedProfilesWriterSingleton *use_case.NamedProfilesWriter
var namedProfilesWriterMutex sync.Mutex

func (prov *Providers) GetNamedProfilesWriter() *use_case.NamedProfilesWriter {
	namedProfilesWriterMutex.Lock()
	defer namedProfilesWriterMutex.Unlock()

	if namedProfilesWriterSingleton == nil {
		namedProfilesWriterSingleton = &use_case.NamedProfilesWriter{
			ConfigurationRepository: prov.GetFileConfigurationRepository(),
		}
	}
	return namedProfilesWriterSingleton
}
