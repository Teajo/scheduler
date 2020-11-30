package queue

import (
	"errors"
	"fmt"
	"jpb/scheduler/task"
	"sync"
	"time"
)

// Queue is a thread-safe queue
type Queue struct {
	mu    sync.Mutex
	end   time.Time
	queue map[string]*task.Task
}

// New creates new thread safe queue
func New(end time.Time) *Queue {
	return &Queue{
		end:   end,
		queue: make(map[string]*task.Task),
	}
}

// Add adds a task to queue
func (sq *Queue) Add(t *task.Task) error {
	sq.mu.Lock()
	defer sq.mu.Unlock()

	if t.Date.After(sq.end) {
		return errors.New("task ends after queue end")
	}

	// ensure task is instantiated only once
	if _, ok := sq.queue[t.ID]; !ok {
		sq.queue[t.ID] = t
	}

	return nil
}

// Remove removes a task
func (sq *Queue) Remove(id string) error {
	sq.mu.Lock()
	defer sq.mu.Unlock()

	if _, ok := sq.queue[id]; !ok {
		return fmt.Errorf("task %s does not exist in queue", id)
	}

	sq.queue[id].Cancel()
	delete(sq.queue, id)

	return nil
}

// Len returns queue length
func (sq *Queue) Len() int {
	sq.mu.Lock()
	defer sq.mu.Unlock()

	return len(sq.queue)
}

// UpdateEnd updates end date of queue
func (sq *Queue) UpdateEnd(date time.Time) {
	sq.end = date
}
