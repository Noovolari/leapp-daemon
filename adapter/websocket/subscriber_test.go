package websocket_test

import (
  "encoding/json"
  "fmt"
  "leapp_daemon/adapter/websocket"
  "leapp_daemon/test/mock"
  "leapp_daemon/use_case/notification"
  "reflect"
  "testing"
  "time"
)

var websocketConnectionMock *mock.ConnectionMock
var tickerMock *mock.TickerMock
var websocketSubscriber *websocket.Subscriber
var messageMock notification.Message

func websocketSubscriberSetup() {
  websocketConnectionMock = mock.NewConnectionMock()
  tickerMock = mock.NewTickerMock()
  websocketSubscriber = websocket.NewWebsocketSubscriber(websocketConnectionMock, tickerMock)
  messageMock = notification.Message{
    MessageType: 0,
    Data:        "fake-message",
  }
}

func websocketSubscriberVerifyExpectedCalls(t *testing.T, websocketConnectionMockCalls []string, tickerMockCalls []string) {
  if !reflect.DeepEqual(websocketConnectionMock.GetCalls(), websocketConnectionMockCalls) {
    t.Fatalf("websocketConnectionMock expectation violation.\nMock calls: %v", websocketConnectionMock.GetCalls())
  }
  if !reflect.DeepEqual(tickerMock.GetCalls(), tickerMockCalls) {
    t.Fatalf("tickerMock expectation violation.\nMock calls: %v", tickerMock.GetCalls())
  }
}

func TestNotify_InvokesConnectionWriteMessageOnceWithMarshalledMessage(t *testing.T) {
  websocketSubscriberSetup()
  websocketSubscriber.Notify(messageMock)
  marshalledMessage, _ := json.Marshal(messageMock)
  websocketSubscriberVerifyExpectedCalls(
    t,
    []string{fmt.Sprintf("WriteMessage(%+v)", marshalledMessage)},
    []string{},
  )
}

func TestNotify_InvokesConnectionClose_IfConnectionWriteMessageReturnsAnError(t *testing.T) {
  websocketSubscriberSetup()
  websocketConnectionMock.MakeWriteMessageReturnError = true
  websocketSubscriber.Notify(messageMock)
  marshalledMessage, _ := json.Marshal(messageMock)
  websocketSubscriberVerifyExpectedCalls(
    t,
    []string{
      fmt.Sprintf("WriteMessage(%+v)", marshalledMessage),
      fmt.Sprintf("Close()"),
    },
    []string{},
  )
}

func TestNotify_ReturnsConnectionWriteMessageReturnedError(t *testing.T) {
  websocketSubscriberSetup()
  websocketConnectionMock.MakeWriteMessageReturnError = true
  err := websocketSubscriber.Notify(messageMock)
  reflect.DeepEqual(err, fmt.Errorf("fake-error"))
}

func TestReadPump_InvokesConnectionReadMessageAndClose_IfConnectionReadMessageReturnsAnError(t *testing.T) {
  websocketSubscriberSetup()
  websocketConnectionMock.MakeReadMessageReturnError = true
  websocketSubscriber.ReadPump()
  websocketSubscriberVerifyExpectedCalls(
    t,
    []string{
      fmt.Sprintf("ReadMessage()"),
      fmt.Sprintf("Close()"),
    },
    []string{},
  )
}

func TestWritePump_InvokesExpectedMethods_IfConnectionPingReturnsAnError(t *testing.T) {
  websocketSubscriberSetup()
  websocketConnectionMock.MakePingReturnError = true
  tickerMock.SetChannel(time.NewTicker(1 * time.Second).C)
  websocketSubscriber.WritePump()
  websocketSubscriberVerifyExpectedCalls(
    t,
    []string{
      fmt.Sprintf("Ping()"),
      fmt.Sprintf("Close()"),
    },
    []string{
      fmt.Sprintf("GetChannel()"),
      fmt.Sprintf("Stop()"),
    },
  )
}
