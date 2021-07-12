package websocket

import "time"

const (
  // Time allowed to read the next pong Message from the peer.
  PongWait = 60 * time.Second

  // Send pings to peer with this period. Must be less than pongWait.
  PingPeriod = (PongWait * 9) / 10
)
