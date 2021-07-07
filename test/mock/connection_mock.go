package mock

import (
  "fmt"
)

type ConnectionMock struct {
  calls []string
  MakeWriteMessageReturnError bool
  MakeReadMessageReturnError bool
  MakePingReturnError bool
}

func NewConnectionMock() *ConnectionMock {
  return &ConnectionMock{calls: []string{}}
}

func (connectionMock *ConnectionMock) GetCalls() []string {
  return connectionMock.calls
}

func (connectionMock *ConnectionMock) Setup() error {
  connectionMock.calls = append(connectionMock.calls, fmt.Sprintf("Setup()"))
  return nil
}

func (connectionMock *ConnectionMock) ReadMessage() error {
  connectionMock.calls = append(connectionMock.calls, fmt.Sprintf("ReadMessage()"))
  if connectionMock.MakeReadMessageReturnError {
    return fmt.Errorf("fake-error")
  }
  return nil
}

func (connectionMock *ConnectionMock) WriteMessage(payload []byte) error {
  connectionMock.calls = append(connectionMock.calls, fmt.Sprintf("WriteMessage(%+v)", payload))
  if connectionMock.MakeWriteMessageReturnError {
    return fmt.Errorf("fake-error")
  }
  return nil
}

func (connectionMock *ConnectionMock) Ping() error {
  connectionMock.calls = append(connectionMock.calls, fmt.Sprintf("Ping()"))
  if connectionMock.MakePingReturnError {
    return fmt.Errorf("fake-error")
  }
  return nil
}

func (connectionMock *ConnectionMock) Close() {
  connectionMock.calls = append(connectionMock.calls, fmt.Sprintf("Close()"))
}
