package websocket

import (
  "encoding/json"
  "leapp_daemon/use_case/notification"
  "time"
)

type Connection interface {
  Setup() error
  ReadMessage() error
  WriteMessage(payload []byte) error
  Ping() error
  Close()
}

type Ticker interface {
  GetChannel() <-chan time.Time
  Stop()
}

type Subscriber struct {
  Connection Connection
  Ticker     Ticker
}

func NewWebsocketSubscriber(connection Connection, ticker Ticker) *Subscriber {
  subscriber := Subscriber{
    Connection: connection,
    Ticker: ticker,
  }
  return &subscriber
}

func (subscriber *Subscriber) Notify(message notification.Message) error {
  payload, _ := json.Marshal(message)

  err := subscriber.Connection.WriteMessage(payload)
  if err != nil {
    subscriber.Connection.Close()
    return err
  }

  return nil
}

func (subscriber *Subscriber) ReadPump() {
  defer func() {
    subscriber.Connection.Close()
  }()

  for {
    err := subscriber.Connection.ReadMessage()
    if err != nil {
      break
    }
  }
}

func (subscriber *Subscriber) WritePump() {
  defer func() {
    subscriber.Ticker.Stop()
    subscriber.Connection.Close()
  }()

  for {
    select {
    case <-subscriber.Ticker.GetChannel():
      err := subscriber.Connection.Ping()
      if err != nil {
        return
      }
    }
  }
}
