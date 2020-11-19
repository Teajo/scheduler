package taskqueue

import (
	"fmt"
	"jpb/scheduler/db"
	"jpb/scheduler/logger"
	"jpb/scheduler/task"
	"jpb/scheduler/taskqueue/queue"
	"jpb/scheduler/utils"
	"time"
)

// TaskQueue represents a queue of tasks
type TaskQueue struct {
	db        db.Taskdb
	queue     *queue.Queue
	timeChunk time.Duration
	taskDone  chan *utils.Scheduling
}

// New creates a new taskqueue
func New(db db.Taskdb, taskDone chan *utils.Scheduling, timeChunk time.Duration) *TaskQueue {
	return &TaskQueue{
		db:        db,
		queue:     queue.New(utils.LastDate),
		timeChunk: timeChunk,
		taskDone:  taskDone,
	}
}

// Add adds task to queue
func (q *TaskQueue) Add(scheduling *utils.Scheduling) (string, error) {
	err := q.db.StoreTask(scheduling)
	if err != nil {
		return scheduling.ID, err
	}

	id, err := q.createTask(scheduling)
	if err != nil {
		logger.Error(err.Error())
	}

	return id, nil
}

// LoadTasks loads tasks from db
func (q *TaskQueue) LoadTasks() {
	now := time.Now()
	nextEnd := q.getNextEnd(now).Add(time.Second)
	q.queue.UpdateEnd(nextEnd)
	setTimer(now.Add(q.timeChunk).Sub(now), q.LoadTasks)

	tasks := q.db.GetTasksToDo(utils.FirstDate, nextEnd)
	for _, t := range tasks {
		_, err := q.createTask(t)
		if err != nil {
			logger.Error(err.Error())
		}
	}
}

func (q *TaskQueue) onTaskDone() func(*utils.Scheduling) {
	return func(scheduling *utils.Scheduling) {
		logger.Info("task done", scheduling.ID)
		q.notifyTaskDone(scheduling)
		q.ackTask(scheduling.ID)
	}
}

func (q *TaskQueue) notifyTaskDone(scheduling *utils.Scheduling) {
	q.taskDone <- scheduling
}

func (q *TaskQueue) createTask(scheduling *utils.Scheduling) (string, error) {
	task := task.New(scheduling)

	err := q.queue.Add(task)
	if err != nil {
		return task.ID, err
	}

	go task.Do(q.onTaskDone())

	logger.Info(fmt.Sprintf("task %s created", scheduling.ID))

	return task.ID, nil
}

func (q *TaskQueue) ackTask(id string) error {
	err := q.db.AckTask(id)
	if err != nil {
		return err
	}

	err = q.queue.Remove(id)
	return err
}

func (q *TaskQueue) getNextEnd(now time.Time) time.Time {
	return now.Add(q.timeChunk)
}

func setTimer(delay time.Duration, cb func()) {
	go func() {
		<-time.After(delay)
		cb()
	}()
}
