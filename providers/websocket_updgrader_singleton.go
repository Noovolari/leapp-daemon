package providers

import (
  "github.com/gorilla/websocket"
  "net/http"
  "sync"
)

var websocketUpgraderSingleton *websocket.Upgrader
var websocketUpgraderMutex sync.Mutex

func (prov *Providers) GetWebsocketUpgrader() *websocket.Upgrader {
  websocketUpgraderMutex.Lock()
  defer websocketUpgraderMutex.Unlock()

  if websocketUpgraderSingleton == nil {
    websocketUpgraderSingleton = &websocket.Upgrader{
      ReadBufferSize:  1024,
      WriteBufferSize: 1024,
      CheckOrigin: func(r *http.Request) bool {
        return true
      },
    }
  }

  return websocketUpgraderSingleton
}
