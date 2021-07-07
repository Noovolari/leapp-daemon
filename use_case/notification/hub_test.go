package notification_test

import (
  "leapp_daemon/test/mock"
  "leapp_daemon/use_case/notification"
  "reflect"
  "testing"
)

var websocketSubscriberMock mock.SubscriberMock
var notificationHub *notification.Hub

func notificationHubSetup() {
  websocketSubscriberMock = mock.NewSubscriberMock()
  notificationHub = notification.NewHub()
  notificationHub.Subscribe(&websocketSubscriberMock)
}

func notificationHubVerifyExpectedCalls(t *testing.T, websocketSubscriberMockCalls []string) {
  if !reflect.DeepEqual(websocketSubscriberMock.GetCalls(), websocketSubscriberMockCalls) {
    t.Fatalf("websocketSubscriberMock expectation violation.\nMock calls: %v", websocketSubscriberMock.GetCalls())
  }
}

func TestBroadcastMessage_NotifiesAllSubscribers(t *testing.T) {
  notificationHubSetup()
  notificationHub.BroadcastMessage(notification.Message{
    MessageType: notification.MfaTokenRequest,
    Data:        "fake-data",
  })
  notificationHubVerifyExpectedCalls(t, []string{"Notify({MessageType:0 Data:fake-data})"})
}

func TestBroadcastMessage_RemovesSubscriber_IfNotifyGeneratesAnError(t *testing.T) {
  notificationHubSetup()
  websocketSubscriberMock.MakeNotifyReturnError = true
  notificationHub.BroadcastMessage(notification.Message{
    MessageType: notification.MfaTokenRequest,
    Data:        "fake-data",
  })
  notificationHub.BroadcastMessage(notification.Message{
    MessageType: notification.MfaTokenRequest,
    Data:        "fake-data",
  })
  notificationHubVerifyExpectedCalls(t, []string{"Notify({MessageType:0 Data:fake-data})"})
}
