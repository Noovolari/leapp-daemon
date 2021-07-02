package providers

import (
	"leapp_daemon/infrastructure/file_system"
	"sync"
)

var fileSystemSingleton *file_system.FileSystem
var fileSystemMutex sync.Mutex

func (prov *Providers) GetFileSystem() *file_system.FileSystem {
	fileSystemMutex.Lock()
	defer fileSystemMutex.Unlock()

	if fileSystemSingleton == nil {
		fileSystemSingleton = &file_system.FileSystem{}
	}
	return fileSystemSingleton
}
