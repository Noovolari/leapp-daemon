package websocket

import (
  "github.com/gorilla/websocket"
  "leapp_daemon/infrastructure/logging"
  "time"
)

const (
  // Maximum Message size allowed from peer.
  maxMessageSize = 512

  // Time allowed to read the next pong Message from the peer.
  pongWait = 60 * time.Second

  // Time allowed to write a Message to the peer.
  writeWait = 10 * time.Second
)

type WebsocketConnectionWrapper struct {
  WebsocketConnection *websocket.Conn
}

func NewWebsocketConnectionWrapper(websocketConnection *websocket.Conn) *WebsocketConnectionWrapper {
  return &WebsocketConnectionWrapper{
    WebsocketConnection: websocketConnection,
  }
}

func (connection *WebsocketConnectionWrapper) Setup() error {
  connection.WebsocketConnection.SetReadLimit(maxMessageSize)
  err := connection.WebsocketConnection.SetReadDeadline(time.Now().Add(pongWait))
  if err != nil {
    return err
  }
  connection.WebsocketConnection.SetPongHandler(func(string) error {
    _ = connection.WebsocketConnection.SetReadDeadline(time.Now().Add(pongWait))
    return nil
  })
  connection.WebsocketConnection.SetCloseHandler(func(int, string) error {
    _ = connection.WebsocketConnection.Close()
    return nil
  })
  return nil
}

func (connection *WebsocketConnectionWrapper) ReadMessage() error {
  _, _, err := connection.WebsocketConnection.ReadMessage()
  if err != nil {
    if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
      logging.Entry().Errorf("error reading WebSocket message: %+v", err)
    }
    return err
  }
  return nil
}

func (connection *WebsocketConnectionWrapper) WriteMessage(payload []byte) error {
  err := connection.writeMessage(websocket.TextMessage, payload)
  if err != nil {
    return err
  }
  return nil
}

func (connection *WebsocketConnectionWrapper) Ping() error {
  var payload []byte
  err := connection.writeMessage(websocket.PingMessage, payload)
  if err != nil {
    return err
  }
  return nil
}

func (connection *WebsocketConnectionWrapper) Close() {
  _ = connection.WebsocketConnection.Close()
}

func (connection *WebsocketConnectionWrapper) writeMessage(messageType int, payload []byte) error {
  err := connection.WebsocketConnection.SetWriteDeadline(time.Now().Add(writeWait))
  if err != nil {
    return err
  }

  err = connection.WebsocketConnection.WriteMessage(messageType, payload)
  if err != nil {
    return err
  }

  return nil
}
