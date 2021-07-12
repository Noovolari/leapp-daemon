package mock

import (
	"errors"
	"fmt"
)

type FileSystemMock struct {
	calls                  []string
	ExpErrorOnGetHomeDir   bool
	ExpErrorOnWriteToFile  bool
	ExpErrorOnRemoveFile   bool
	ExpErrorOnRenamingFile bool
	ExpHomeDir             string
	ExpDoesFileExist       bool
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
	fileSystem.calls = append(fileSystem.calls, "GetHomeDir()")
	if fileSystem.ExpErrorOnGetHomeDir {
		return "", errors.New("error getting home dir")
	}

	return fileSystem.ExpHomeDir, nil
}

func (fileSystem *FileSystemMock) WriteToFile(path string, data []byte) error {
	fileSystem.calls = append(fileSystem.calls, fmt.Sprintf("WriteToFile(%v, %v)", path, data))
	if fileSystem.ExpErrorOnWriteToFile {
		return errors.New("error writing file")
	}

	return nil
}

func (fileSystem *FileSystemMock) RemoveFile(path string) error {
	fileSystem.calls = append(fileSystem.calls, fmt.Sprintf("RemoveFile(%v)", path))
	if fileSystem.ExpErrorOnRemoveFile {
		return errors.New("error removing file")
	}

	return nil
}

func (fileSystem *FileSystemMock) RenameFile(from string, to string) error {
	fileSystem.calls = append(fileSystem.calls, fmt.Sprintf("RenameFile(%v, %v)", from, to))
	if fileSystem.ExpErrorOnRenamingFile {
		return errors.New("error renaming file")
	}

	return nil
}
