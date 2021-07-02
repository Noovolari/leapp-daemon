package providers

import (
	"leapp_daemon/use_case"
	"sync"
)

var awsSessionsWriterSingleton *use_case.AwsSessionsWriter
var awsSessionsWriterMutex sync.Mutex

func (prov *Providers) GetAwsSessionWriter() *use_case.AwsSessionsWriter {
	awsSessionsWriterMutex.Lock()
	defer awsSessionsWriterMutex.Unlock()

	if awsSessionsWriterSingleton == nil {
		awsSessionsWriterSingleton = &use_case.AwsSessionsWriter{
			ConfigurationRepository: prov.GetFileConfigurationRepository(),
		}
	}
	return awsSessionsWriterSingleton
}
