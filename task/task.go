package task

import (
	"jpb/scheduler/utils"
	"time"
)

// Task structure
type Task struct {
	ID         string
	Date       time.Time
	onDone     func(*utils.Scheduling)
	done       bool
	cancel     *chan interface{}
	scheduling *utils.Scheduling
}

// New creates Task
func New(scheduling *utils.Scheduling, onDone func(*utils.Scheduling)) *Task {
	cancel := make(chan interface{})

	t := &Task{
		ID:         scheduling.ID,
		onDone:     onDone,
		cancel:     &cancel,
		Date:       scheduling.Date,
		done:       false,
		scheduling: scheduling,
	}

	go t.doTask()
	return t
}

// Cancel cancels task
func (t *Task) Cancel() {
	if !t.done {
		*t.cancel <- struct{}{}
	}
}

func (t *Task) doTask() {
	now := time.Now()
	duration := t.Date.Sub(now)

	for {
		select {
		case <-time.After(duration):
			t.done = true
			t.onDone(t.scheduling)
			return
		case <-*t.cancel:
			return
		}
	}
}
