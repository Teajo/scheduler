package ticker

import (
	"jpb/scheduler/events"
	"time"
)

// Ticker is a ticker
type Ticker struct {
	Bus *events.Bus `inject:""`
}

// Start starts a ticker
func (t *Ticker) Start(period time.Duration) {
	ticker := time.NewTicker(period)

	go func() {
		t.Bus.Emit(events.TICK, time.Now())

		for {
			d := <-ticker.C
			t.Bus.Emit(events.TICK, d)
		}
	}()
}
