package taskqueue

import (
	"errors"
	"jpb/scheduler/db"
	"jpb/scheduler/task"
	"sort"
	"time"
)

// The very last existing date
var veryLast = time.Unix(1<<63-1, 0)

// TaskQueue represents a queue of tasks
type TaskQueue struct {
	queue    []*task.Task // WARNING, only access queue from listener goroutine
	maxLen   int
	last     *time.Time
	listener chan *QueueRequest
	db       db.Taskdb
}

// New creates a new taskqueue
func New(maxLen int) *TaskQueue {
	tq := &TaskQueue{
		queue:    []*task.Task{},
		maxLen:   maxLen,
		last:     &veryLast,
		listener: make(chan *QueueRequest),
		db:       &db.Fakedb{},
	}

	tq.listen()
	return tq
}

// Add adds task to queue
func (q *TaskQueue) Add(date *time.Time, todo *func(string)) (string, error) {
	task := task.New(date, q.onTaskDone(todo))
	err := q.db.StoreTask(task)
	if err != nil {
		return task.ID, err
	}

	return task.ID, q.send(add, task)
}

// Remove removes task from queue
func (q *TaskQueue) Remove(id string) error {
	err := q.db.RemoveTask(id)
	if err != nil {
		// treat error
	}

	err = q.send(remove, id)
	return err
}

// Stop stops task queue
func (q *TaskQueue) Stop() error {
	err := q.send(cancel, nil)
	for _, t := range q.queue {
		t.Cancel()
	}
	return err
}

// Len returns queue length
func (q *TaskQueue) Len() int {
	return len(q.queue)
}

func (q *TaskQueue) add(task *task.Task) error {
	if q.Len() >= q.maxLen && task.Date.After(*q.last) {
		return errors.New("Max queue length reached")
	}

	q.queue = append(q.queue, task)
	q.sortQueueByDate()

	// removes tasks in excess
	for q.Len() > q.maxLen {
		index := q.Len() - 1
		task := q.queue[index]
		task.Cancel()
		q.removeByIndex(index)
	}

	q.updateLast()
	return nil
}

func (q *TaskQueue) remove(id string) error {
	task, index := q.getTaskByID(id)
	if task != nil {
		task.Cancel()
		q.removeByIndex(index)
	}

	q.updateLast()
	return nil
}

func (q *TaskQueue) send(method method, payload interface{}) error {
	r, errChan := NewQueueRequest(method, payload)
	q.listener <- r
	err := <-errChan
	return err
}

func (q *TaskQueue) listen() {
	go func() {
		for {
			r := <-q.listener
			switch r.method {
			case add:
				task := (r.payload).(*task.Task)
				err := q.add(task)
				r.err <- err
			case remove:
				id := (r.payload).(string)
				err := q.remove(id)
				r.err <- err
			case cancel:
				r.err <- nil
				return
			default:
				panic("Not handled method in listener")
			}
		}
	}()
}

func (q *TaskQueue) removeByIndex(index int) {
	q.queue = append(q.queue[:index], q.queue[index+1:]...)
}

func (q *TaskQueue) getTaskByID(id string) (*task.Task, int) {
	for k, v := range q.queue {
		if id == v.ID {
			return v, k
		}
	}
	return nil, -1
}

func (q *TaskQueue) sortQueueByDate() {
	sort.Slice(q.queue, func(i, j int) bool {
		return q.queue[j].Date.After(*q.queue[i].Date)
	})
}

func (q *TaskQueue) updateLast() {
	if q.Len() < 1 {
		q.last = &veryLast
		return
	}
	q.last = q.queue[q.Len()-1].Date
}

func (q *TaskQueue) onTaskDone(todo *func(string)) func(string) {
	return func(id string) {
		q.Remove(id)
		if todo != nil {
			(*todo)(id)
		}
	}
}
