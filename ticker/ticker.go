package ticker

import (
	"jpb/scheduler/events"
	"time"
)

// Ticker is a ticker
type Ticker struct {
	stop chan interface{}
}

// New returns a new ticker
func New() *Ticker {
	return &Ticker{make(chan interface{})}
}

// Start starts a ticker
func (t *Ticker) Start(period time.Duration, bus events.Bus) {
	ticker := time.NewTicker(period)

	go func() {
		for {
			select {
			case t := <-ticker.C:
				bus.Emit(events.TICK, t)
			case <-t.stop:
				return
			}
		}
	}()
}

// Stop stops ticker
func (t *Ticker) Stop() {
	t.stop <- nil
}
