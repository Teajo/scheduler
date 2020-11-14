package taskqueue

import (
	"fmt"
	"jpb/scheduler/db"
	"jpb/scheduler/logger"
	"jpb/scheduler/task"
	"jpb/scheduler/taskqueue/sortedqueue"
	"jpb/scheduler/utils"
)

// TaskQueue represents a queue of tasks
type TaskQueue struct {
	db       db.Taskdb
	queue    *sortedqueue.SortedQueue
	taskDone chan *utils.Scheduling
}

// New creates a new taskqueue
func New(db db.Taskdb, maxLen int, taskDone chan *utils.Scheduling) *TaskQueue {
	return &TaskQueue{
		db:       db,
		queue:    sortedqueue.New(maxLen),
		taskDone: taskDone,
	}
}

// Add adds task to queue
func (q *TaskQueue) Add(scheduling *utils.Scheduling) (string, error) {
	err := q.db.StoreTask(scheduling)
	if err != nil {
		return scheduling.ID, err
	}
	return q.createTask(scheduling)
}

// LoadTasks loads tasks from db
func (q *TaskQueue) LoadTasks() error {
	tasks := q.db.GetTasks(q.queue.LastID(), q.queue.EmptyLen(), q.queue.LastDate())
	for _, t := range tasks {
		q.createTask(t)
	}
	return nil
}

// Stop stops task queue
func (q *TaskQueue) Stop() error {
	for _, t := range q.queue.Get() {
		t.Cancel()
	}
	return nil
}

func (q *TaskQueue) onTaskDone() func(*utils.Scheduling) {
	return func(scheduling *utils.Scheduling) {
		logger.Info("task done", scheduling.ID)
		q.notifyTaskDone(scheduling)
		q.ackTask(scheduling.ID)
		q.LoadTasks()
	}
}

func (q *TaskQueue) notifyTaskDone(scheduling *utils.Scheduling) {
	q.taskDone <- scheduling
}

func (q *TaskQueue) createTask(scheduling *utils.Scheduling) (string, error) {
	logger.Info(fmt.Sprintf("create task %s", scheduling.ID))
	task := task.New(scheduling)
	go task.Do(q.onTaskDone())

	tasks := q.queue.Add(task)
	for _, t := range tasks {
		t.Cancel()
	}

	return task.ID, nil
}

func (q *TaskQueue) ackTask(id string) error {
	err := q.db.AckTask(id)
	if err != nil {
		return err
	}

	err = q.queue.RemoveByID(id)
	return err
}
