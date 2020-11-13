package taskqueue

import (
	"fmt"
	"jpb/scheduler/db"
	"jpb/scheduler/task"
	"jpb/scheduler/taskqueue/sortedqueue"
	"jpb/scheduler/utils"
)

type method string

const (
	add    method = "add"
	remove method = "remove"
	cancel method = "cancel"
)

// TaskQueue represents a queue of tasks
type TaskQueue struct {
	db       db.Taskdb
	queue    *sortedqueue.SortedQueue // WARNING, only access queue from listener goroutine
	listener chan *QueueRequest
	taskDone chan *utils.Scheduling
}

// New creates a new taskqueue
func New(db db.Taskdb, maxLen int, taskDone chan *utils.Scheduling) *TaskQueue {
	return &TaskQueue{
		db:       db,
		queue:    sortedqueue.New(maxLen),
		listener: make(chan *QueueRequest),
		taskDone: taskDone,
	}
}

// LoadTasks loads tasks from db
func (q *TaskQueue) LoadTasks() error {
	tasks := q.db.GetFirstTasks(0)
	for _, t := range tasks {
		fmt.Println("load task", t.ID)
		q.Add(t)
	}
	return nil
}

// Add adds task to queue
func (q *TaskQueue) Add(scheduling *utils.Scheduling) (string, error) {
	task := task.New(scheduling)
	err := q.db.StoreTask(task.Scheduling)
	if err != nil {
		return task.ID, err
	}

	go task.Do(q.onTaskDone())

	return task.ID, q.send(add, task)
}

// Remove removes task from queue
func (q *TaskQueue) Remove(id string) error {
	err := q.db.RemoveTask(id)
	if err != nil {
		return err
	}

	err = q.send(remove, id)
	return err
}

// Stop stops task queue
func (q *TaskQueue) Stop() error {
	err := q.send(cancel, nil)
	for _, t := range q.queue.Get() {
		t.Cancel()
	}
	return err
}

// Len returns queue length
func (q *TaskQueue) Len() int {
	return q.queue.Len()
}

// Listen make taskqueue listening for task events
func (q *TaskQueue) Listen() {
	fmt.Println("taskqueue listening for task events")
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
}

func (q *TaskQueue) add(task *task.Task) error {
	tasks := q.queue.Insert(task)
	for _, t := range tasks {
		t.Cancel()
	}
	return nil
}

func (q *TaskQueue) remove(id string) error {
	return q.queue.RemoveByID(id)
}

func (q *TaskQueue) onTaskDone() func(*utils.Scheduling) {
	return func(scheduling *utils.Scheduling) {
		q.taskDone <- scheduling
		q.Remove(scheduling.ID)
	}
}

func (q *TaskQueue) send(method method, payload interface{}) error {
	r, errChan := NewQueueRequest(method, payload)
	q.listener <- r
	err := <-errChan
	return err
}
