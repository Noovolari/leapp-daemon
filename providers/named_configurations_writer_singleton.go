package providers

import (
	"leapp_daemon/use_case"
	"sync"
)

var namedConfigurationsWriterSingleton *use_case.NamedConfigurationsWriter
var namedConfigurationsWriterMutex sync.Mutex

func (prov *Providers) GetNamedConfigurationsWriter() *use_case.NamedConfigurationsWriter {
	namedConfigurationsWriterMutex.Lock()
	defer namedConfigurationsWriterMutex.Unlock()

	if namedConfigurationsWriterSingleton == nil {
		namedConfigurationsWriterSingleton = &use_case.NamedConfigurationsWriter{
			ConfigurationRepository: prov.GetFileConfigurationRepository(),
		}
	}
	return namedConfigurationsWriterSingleton
}
