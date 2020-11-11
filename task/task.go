package task

import (
	"time"

	"github.com/google/uuid"
)

// Task structure
type Task struct {
	ID     string
	Date   *time.Time
	onDone func(string)
	done   bool
	cancel *chan interface{}
}

// New creates Task
func New(date *time.Time, onDone func(string)) *Task {
	cancel := make(chan interface{})
	id, _ := uuid.NewUUID()

	t := &Task{
		ID:     id.String(),
		onDone: onDone,
		cancel: &cancel,
		Date:   date,
		done:   false,
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
			t.onDone(t.ID)
			return
		case <-*t.cancel:
			return
		}
	}
}
