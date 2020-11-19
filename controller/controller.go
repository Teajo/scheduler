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
	db       db.Taskdb
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
	queue := taskqueue.New(db, taskDone, cfg.TimeChunk)
	queue.LoadTasks()

	return &Ctrl{
		db:       db,
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

// GetTasks returns tasks from db
func (c *Ctrl) GetTasks(start time.Time, end time.Time) []*utils.Scheduling {
	return c.db.GetTasks(start, end)
}
