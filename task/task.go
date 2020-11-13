package task

import (
	"jpb/scheduler/utils"
	"time"
)

// Task structure
type Task struct {
	ID         string
	Date       time.Time
	Scheduling *utils.Scheduling
	done       bool
	cancel     chan interface{}
}

// New creates Task
func New(scheduling *utils.Scheduling) *Task {
	return &Task{
		ID:         scheduling.ID,
		cancel:     make(chan interface{}),
		Date:       scheduling.Date,
		done:       false,
		Scheduling: scheduling,
	}
}

// Cancel cancels task
func (t *Task) Cancel() {
	if !t.done {
		t.cancel <- struct{}{}
	}
}

// Do starts scheduled task
func (t *Task) Do(onDone func(*utils.Scheduling)) {
	now := time.Now()
	duration := t.Date.Sub(now)

	for {
		select {
		case <-time.After(duration):
			t.done = true
			onDone(t.Scheduling)
			return
		case <-t.cancel:
			return
		}
	}
}
