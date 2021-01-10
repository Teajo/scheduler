package queue

import (
	"fmt"
	"jpb/scheduler/task"
	"sync"
)

// Queue is a thread-safe queue
type Queue struct {
	mu    sync.Mutex
	queue map[string]*task.Task
}

// New creates new thread safe queue
func New() *Queue {
	return &Queue{
		queue: make(map[string]*task.Task),
	}
}

// Add adds a task to queue
func (sq *Queue) Add(t *task.Task) error {
	sq.mu.Lock()
	defer sq.mu.Unlock()

	// ensure task is instantiated only once
	if _, ok := sq.queue[t.ID]; ok {
		return fmt.Errorf("task %s already exists", t.ID)
	}

	sq.queue[t.ID] = t
	return nil
}

// Remove removes a task
func (sq *Queue) Remove(id string) {
	sq.mu.Lock()
	defer sq.mu.Unlock()

	if _, ok := sq.queue[id]; !ok {
		return
	}

	sq.queue[id].Cancel()
	delete(sq.queue, id)
}

// Len returns queue length
func (sq *Queue) Len() int {
	sq.mu.Lock()
	defer sq.mu.Unlock()

	return len(sq.queue)
}
