package providers

import (
  "leapp_daemon/use_case/notification"
  "sync"
)

var notificationHubSingleton *notification.Hub
var notificationHubMutex sync.Mutex

func (prov *Providers) GetNotificationHub() *notification.Hub {
  notificationHubMutex.Lock()
  defer notificationHubMutex.Unlock()

  if notificationHubSingleton == nil {
    notificationHubSingleton = notification.NewHub()
  }

  return notificationHubSingleton
}
