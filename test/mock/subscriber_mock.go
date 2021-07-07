package mock

import (
  "fmt"
  "leapp_daemon/use_case/notification"
)

type SubscriberMock struct {
  calls                 []string
  MakeNotifyReturnError bool
}

func NewSubscriberMock() SubscriberMock {
  return SubscriberMock{calls: []string{}}
}

func (subscriberMock *SubscriberMock) GetCalls() []string {
  return subscriberMock.calls
}

func (subscriberMock *SubscriberMock) Notify(message notification.Message) error {
  subscriberMock.calls = append(subscriberMock.calls, fmt.Sprintf("Notify(%+v)", message))
  if subscriberMock.MakeNotifyReturnError {
    return fmt.Errorf("fake-error")
  }
  return nil
}
