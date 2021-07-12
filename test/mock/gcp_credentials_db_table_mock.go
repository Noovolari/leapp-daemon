package mock

import (
	"errors"
	"fmt"
)

type GcpCredentialsDbTableMock struct {
	calls                     []string
	ExpErrorOnExecInsertQuery bool
	ExpErrorOnExecDeleteQuery bool
}

func NewGcpCredentialsDbTableMock() GcpCredentialsDbTableMock {
	return GcpCredentialsDbTableMock{calls: []string{}}
}

func (table *GcpCredentialsDbTableMock) GetCalls() []string {
	return table.calls
}

func (table *GcpCredentialsDbTableMock) WriteCredentials(sqlFilePath string, accountId string, value string) error {
	table.calls = append(table.calls, fmt.Sprintf("WriteCredentials(%v, %v, %v)", sqlFilePath, accountId, value))

	if table.ExpErrorOnExecInsertQuery {
		return errors.New("error executing insert query")
	}
	return nil
}

func (table *GcpCredentialsDbTableMock) RemoveCredentials(sqlFilePath string, accountId string) error {
	table.calls = append(table.calls, fmt.Sprintf("RemoveCredentials(%v, %v)", sqlFilePath, accountId))

	if table.ExpErrorOnExecDeleteQuery {
		return errors.New("error executing delete query")
	}
	return nil
}
