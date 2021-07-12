package ticker

import (
  "time"
)

type TickerWrapper struct {
  ticker *time.Ticker
}

func NewTickerWrapper(ticker *time.Ticker) *TickerWrapper {
  tickerWrapper := &TickerWrapper{}
  tickerWrapper.ticker = ticker
  return tickerWrapper
}

func (t *TickerWrapper) GetChannel() <-chan time.Time {
  if t.ticker == nil {
    return nil
  }
  return t.ticker.C
}

func (t *TickerWrapper) Stop() {
  if t.ticker != nil {
    t.ticker.Stop()
  }
}
