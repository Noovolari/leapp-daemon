package notification

import "leapp_daemon/infrastructure/logging"

type Subscriber interface {
  Notify(message Message) error
}

type Hub struct {
  subscribers []Subscriber
}

func NewHub() *Hub {
  return &Hub{
    subscribers: make([]Subscriber, 0),
  }
}

func (hub *Hub) Subscribe(subscriber Subscriber) {
  hub.subscribers = append(hub.subscribers, subscriber)
}

func (hub *Hub) BroadcastMessage(message Message) {
  for i, subscriber := range hub.subscribers {
    err := subscriber.Notify(message)
    if err != nil {
      logging.Entry().Errorf("error broadcasting message through WebSocket connection: %v", err)
      hub.subscribers = append(hub.subscribers[:i], hub.subscribers[i+1:]...)
    }
  }
}
