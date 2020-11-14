package sortedqueue

import (
	"errors"
	"jpb/scheduler/logger"
	"jpb/scheduler/task"
	"sort"
	"time"
)

var firstDate time.Time = time.Unix(0, 0)

type method string

const (
	add        method = "add"
	remove     method = "remove"
	removeByID method = "removeByID"
	length     method = "length"
	get        method = "get"
	lastID     method = "lastID"
	lastDate   method = "lastDate"
	cancel     method = "cancel"
)

// SortedQueue is a thread-safe sorted queue
type SortedQueue struct {
	queue    []*task.Task
	maxLen   int
	listener chan *QueueRequest
}

// New creates new sorted queue.
// WARNING all private methods must NOT call public method on pain of thread deadlock
func New(maxLen int) *SortedQueue {
	q := &SortedQueue{
		queue:    []*task.Task{},
		maxLen:   maxLen,
		listener: make(chan *QueueRequest),
	}

	go q.listen()
	return q
}

// Add adds ordered task. It returns a list of exceeding tasks
func (sq *SortedQueue) Add(t *task.Task) []*task.Task {
	res := sq.send(add, t)
	tasks, ok := (res.payload).([]*task.Task)
	if !ok {
		// treat casting error
	}
	return tasks
}

func (sq *SortedQueue) add(task *task.Task) []*task.Task {
	tasks := emptyTaskList()
	sq.queue = append(sq.queue, task)
	sq.sort()

	// removes tasks in excess
	for sq.len() > sq.maxLen {
		index := sq.len() - 1
		tasks = append(tasks, sq.queue[index])
		sq.removeByIndex(index)
	}

	return tasks
}

// Remove removes a task
func (sq *SortedQueue) Remove(t *task.Task) error {
	res := sq.send(remove, t)
	return res.err
}

func (sq *SortedQueue) remove(t *task.Task) error {
	i := sq.findIndex(t)
	if i < 0 {
		return errors.New("task does not exist")
	}

	sq.removeByIndex(i)
	return nil
}

// RemoveByID removed a task according to provided ID
func (sq *SortedQueue) RemoveByID(id string) error {
	res := sq.send(removeByID, id)
	return res.err
}

func (sq *SortedQueue) removeByID(id string) error {
	i := sq.findIndexByID(id)
	if i < 0 {
		return errors.New("task does not exist")
	}

	sq.removeByIndex(i)
	return nil
}

// Len returns queue length
func (sq *SortedQueue) Len() int {
	res := sq.send(length, nil)
	length, ok := res.payload.(int)
	if !ok {
		// treat casting error
	}
	return length
}

func (sq *SortedQueue) len() int {
	return len(sq.queue)
}

// Get returns task slice
func (sq *SortedQueue) Get() []*task.Task {
	res := sq.send(get, nil)
	tasks, ok := res.payload.([]*task.Task)
	if !ok {
		// treat casting error
	}
	return tasks
}

func (sq *SortedQueue) get() []*task.Task {
	return sq.queue
}

// LastID returns last task id, "" otherwise
func (sq *SortedQueue) LastID() string {
	res := sq.send(lastID, nil)
	tasks, ok := res.payload.(string)
	if !ok {
		// treat casting error
	}
	return tasks
}

func (sq *SortedQueue) lastID() string {
	if sq.len() > 0 {
		return sq.queue[sq.len()-1].ID
	}
	return ""
}

// EmptyLen returns number of empty items
func (sq *SortedQueue) EmptyLen() int {
	return sq.maxLen - sq.Len()
}

// LastDate returns last queue date
func (sq *SortedQueue) LastDate() time.Time {
	res := sq.send(lastDate, nil)
	date, ok := res.payload.(time.Time)
	if !ok {
		// treat casting error
	}
	return date
}

func (sq *SortedQueue) lastDate() time.Time {
	if sq.len() > 0 {
		return sq.queue[sq.len()-1].Date
	}
	return firstDate
}

// Listen make sortedqueue listening for task events
func (sq *SortedQueue) listen() {
	logger.Info("taskqueue listening for task events")
	for {
		r := <-sq.listener
		switch r.method {
		case add:
			t := (r.payload).(*task.Task)
			var res interface{} = sq.add(t)
			r.res <- NewResponse(nil, res)
		case remove:
			t := (r.payload).(*task.Task)
			err := sq.remove(t)
			r.res <- NewResponse(err, nil)
		case removeByID:
			id := (r.payload).(string)
			err := sq.removeByID(id)
			r.res <- NewResponse(err, nil)
		case length:
			l := sq.len()
			r.res <- NewResponse(nil, l)
		case get:
			t := sq.get()
			r.res <- NewResponse(nil, t)
		case lastID:
			t := sq.lastID()
			r.res <- NewResponse(nil, t)
		case lastDate:
			d := sq.lastDate()
			r.res <- NewResponse(nil, d)
		case cancel:
			r.res <- nil
			return
		default:
			panic("Not handled method in listener")
		}
	}
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

func (sq *SortedQueue) send(method method, payload interface{}) *Response {
	req, resCh := NewRequest(method, payload)
	sq.listener <- req
	res := <-resCh
	return res
}

func emptyTaskList() []*task.Task {
	return []*task.Task{}
}
