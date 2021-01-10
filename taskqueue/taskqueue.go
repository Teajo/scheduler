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
	Bus   *events.Bus    `inject:""`
	DB    db.Taskdb      `inject:""`
	Cfg   *config.Config `inject:""`
	queue *queue.Queue
	end   time.Time
}

// Start starts taskqueue ticking
func (q *TaskQueue) Start() {
	q.queue = queue.New()
	q.end = time.Now()
	tick := q.Bus.Subscribe(events.TICK)
	go q.onTick(tick)
}

// Add adds task to queue and store it
func (q *TaskQueue) Add(s *utils.Scheduling) error {
	err := q.DB.StoreTask(s)
	if err != nil {
		return err
	}

	if s.Date.Before(q.end) {
		q.appendTask(s)
	}

	return nil
}

// Remove removes task from queue
func (q *TaskQueue) Remove(id string) error {
	q.queue.Remove(id)
	return q.DB.RemoveTask(id)
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
		q.DB.AckTask(scheduling.ID)
	}
}

func (q *TaskQueue) onTick(tick chan interface{}) {
	for {
		payload := <-tick
		logger.Info("on tick")
		now, ok := payload.(time.Time)
		if !ok {
			continue
		}

		q.end = now.Add(q.Cfg.TimeChunk)
		tasks := q.DB.GetTasksToDo(q.end)
		for _, t := range tasks {
			q.appendTask(t)
		}
	}
}
