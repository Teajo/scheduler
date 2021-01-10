package taskqueue

import (
	"jpb/scheduler/config"
	"jpb/scheduler/db"
	"jpb/scheduler/events"
	"jpb/scheduler/logger"
	"jpb/scheduler/task"
	"jpb/scheduler/taskqueue/queue"
	"jpb/scheduler/utils"
	"time"
)

// TaskQueue represents a queue of tasks
type TaskQueue struct {
	queue  *queue.Queue
	bus    *events.Bus
	db     db.Taskdb
	config *config.Config
	end    time.Time
}

// New creates a new taskqueue
func New() *TaskQueue {
	return &TaskQueue{
		queue: queue.New(),
		end:   time.Now(),
	}
}

// Start starts taskqueue ticking
func (q *TaskQueue) Start() {
	tick := q.bus.Subscribe(events.TICK)

	go q.onTick(tick)
}

// Add adds task to queue and store it
func (q *TaskQueue) Add(s *utils.Scheduling) error {
	err := q.db.StoreTask(s)
	if err != nil {
		return err
	}

	if s.Date.Before(q.end) {
		q.appendTask(s)
	}

	return nil
}

// Remove removes task from queue
func (q *TaskQueue) Remove(s *utils.Scheduling) error {
	q.queue.Remove(s.ID)
	return q.db.RemoveTask(s.ID)
}

func (q *TaskQueue) appendTask(s *utils.Scheduling) {
	task := task.New(s)
	err := q.queue.Add(task)
	if err != nil {
		logger.Error(err.Error())
	} else {
		go task.Do(q.onTaskDone())
	}
}

func (q *TaskQueue) onTaskDone() func(*utils.Scheduling) {
	return func(scheduling *utils.Scheduling) {
		logger.Info("task done", scheduling.ID)
		q.queue.Remove(scheduling.ID)
	}
}

func (q *TaskQueue) onTick(tick chan interface{}) {
	for {
		payload := <-tick
		now, ok := payload.(time.Time)
		if !ok {
			continue
		}

		q.end = now.Add(q.config.TimeChunk)
		tasks := q.db.GetTasksToDo(now, q.end)
		for _, t := range tasks {
			q.appendTask(t)
		}
	}
}
