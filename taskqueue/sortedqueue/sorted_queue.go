package sortedqueue

import (
	"errors"
	"jpb/scheduler/task"
	"sort"
)

// SortedQueue represents a sorted queue
type SortedQueue struct {
	queue  []*task.Task
	maxLen int
}

// New creates new sorted queue
func New(maxLen int) *SortedQueue {
	return &SortedQueue{
		queue:  []*task.Task{},
		maxLen: maxLen,
	}
}

// Insert inserts ordered task. It returns a list of exceeding task
func (sq *SortedQueue) Insert(task *task.Task) []*task.Task {
	tasks := emptyTaskList()
	sq.queue = append(sq.queue, task)
	sq.sort()

	// removes tasks in excess
	for sq.Len() > sq.maxLen {
		index := sq.Len() - 1
		tasks = append(tasks, sq.queue[index])
		sq.removeByIndex(index)
	}

	return tasks
}

// Remove removes a task
func (sq *SortedQueue) Remove(t *task.Task) error {
	i := sq.findIndex(t)
	if i < 0 {
		return errors.New("task does not exist")
	}

	sq.removeByIndex(i)
	return nil
}

// RemoveByID removed a task according to provided ID
func (sq *SortedQueue) RemoveByID(id string) error {
	i := sq.findIndexByID(id)
	if i < 0 {
		return errors.New("task does not exist")
	}

	sq.removeByIndex(i)
	return nil
}

// Len returns queue length
func (sq *SortedQueue) Len() int {
	return len(sq.queue)
}

// Get returns task slice
func (sq *SortedQueue) Get() []*task.Task {
	return sq.queue
}

func (sq *SortedQueue) sort() {
	sort.Slice(sq.queue, func(i, j int) bool {
		return sq.queue[j].Date.After(sq.queue[i].Date)
	})
}

func (sq *SortedQueue) removeByIndex(index int) {
	sq.queue = append(sq.queue[:index], sq.queue[index+1:]...)
}

func (sq *SortedQueue) findIndex(t *task.Task) int {
	for i, el := range sq.queue {
		if el == t {
			return i
		}
	}
	return -1
}

func (sq *SortedQueue) findIndexByID(id string) int {
	for i, el := range sq.queue {
		if el.ID == id {
			return i
		}
	}
	return -1
}

func emptyTaskList() []*task.Task {
	return []*task.Task{}
}
