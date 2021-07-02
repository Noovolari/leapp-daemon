package mock

import (
	"errors"
	"fmt"
)

type FileSystemMock struct {
	calls                 []string
	ExpErrorOnGetHomeDir  bool
	ExpErrorOnWriteToFile bool
	ExpErrorOnRemoveFile  bool
	ExpDoesFileExist      bool
}

func (fileSystem *FileSystemMock) GetCalls() []string {
	return fileSystem.calls
}

func NewFileSystemMock() FileSystemMock {
	return FileSystemMock{calls: []string{}}
}

func (fileSystem *FileSystemMock) DoesFileExist(path string) bool {
	fileSystem.calls = append(fileSystem.calls, fmt.Sprintf("DoesFileExist(%v)", path))
	return fileSystem.ExpDoesFileExist
}

func (fileSystem *FileSystemMock) GetHomeDir() (string, error) {
	if fileSystem.ExpErrorOnGetHomeDir {
		return "", errors.New("error getting home dir")
	}
	fileSystem.calls = append(fileSystem.calls, "GetHomeDir()")
	return "/user/home", nil
}

func (fileSystem *FileSystemMock) WriteToFile(path string, data []byte) error {
	if fileSystem.ExpErrorOnWriteToFile {
		return errors.New("error writing file")
	}
	fileSystem.calls = append(fileSystem.calls, fmt.Sprintf("WriteToFile(%v, %v)", path, data))
	return nil
}

func (fileSystem *FileSystemMock) RemoveFile(path string) error {
	if fileSystem.ExpErrorOnRemoveFile {
		return errors.New("error removing file")
	}
	fileSystem.calls = append(fileSystem.calls, fmt.Sprintf("RemoveFile(%v)", path))
	return nil
}
