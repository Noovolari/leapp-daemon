package mock

import (
  "fmt"
  "time"
)

type TickerMock struct {
  calls []string
  channel <-chan time.Time
}

func NewTickerMock() *TickerMock {
  return &TickerMock{calls: []string{}}
}

func (tickerMock *TickerMock) GetCalls() []string {
  return tickerMock.calls
}

func (tickerMock *TickerMock) SetChannel(channel <-chan time.Time) {
  tickerMock.channel = channel
}

func (tickerMock *TickerMock) GetChannel() <-chan time.Time {
  tickerMock.calls = append(tickerMock.calls, fmt.Sprintf("GetChannel()"))
  return tickerMock.channel
}

func (tickerMock *TickerMock) Stop() {
  tickerMock.calls = append(tickerMock.calls, fmt.Sprintf("Stop()"))
}
