package controller

import (
	"fmt"
	"jpb/scheduler/config"
	"jpb/scheduler/db"
	"jpb/scheduler/logger"
	"jpb/scheduler/publisher"
	"jpb/scheduler/taskqueue"
	"jpb/scheduler/utils"
	"time"
)

// Ctrl represents a scheduler controller
type Ctrl struct {
	queue    *taskqueue.TaskQueue
	pubs     *publisher.PubManager
	taskDone chan *utils.Scheduling
}

// New creates a scheduler controller
func New() *Ctrl {
	cfg := config.Get()
	db := db.Getdb(cfg.DbDriver)
	taskDone := make(chan *utils.Scheduling)
	pubs := publisher.New(taskDone)
	queue := taskqueue.New(db, cfg.MaxQueueLen, taskDone)
	queue.LoadTasks()

	return &Ctrl{
		queue:    queue,
		pubs:     pubs,
		taskDone: make(chan *utils.Scheduling),
	}
}

// Schedule schedules a task
func (c *Ctrl) Schedule(scheduling *utils.Scheduling) (string, error) {
	logger.Info("Schedule a task at", scheduling.Date.Format(time.RFC3339Nano))

	publisher, ok := c.pubs.Get(scheduling.Publisher)
	if ok {
		err := publisher.CheckConfig(scheduling.Settings)
		if err != nil {
			return "", err
		}

		return c.queue.Add(scheduling)
	}

	return "", fmt.Errorf("Publisher %s does not exist", scheduling.Publisher)
}

func (c *Ctrl) newQueue(length int) *taskqueue.TaskQueue {
	cfg := config.Get()
	db := db.Getdb(cfg.DbDriver)
	queue := taskqueue.New(db, length, c.taskDone)
	return queue
}
