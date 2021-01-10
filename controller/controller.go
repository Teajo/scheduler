package controller

import (
	"fmt"
	"jpb/scheduler/config"
	"jpb/scheduler/db"
	"jpb/scheduler/events"
	"jpb/scheduler/logger"
	"jpb/scheduler/publisher"
	"jpb/scheduler/taskqueue"
	"jpb/scheduler/utils"
	"time"
)

// Ctrl represents a scheduler controller
type Ctrl struct {
	DB    db.Taskdb             `inject:""`
	Queue *taskqueue.TaskQueue  `inject:""`
	Bus   *events.Bus           `inject:""`
	Cfg   *config.Config        `inject:""`
	Pubs  *publisher.PubManager `inject:""`
}

// Schedule schedules a task
func (c *Ctrl) Schedule(s *utils.Scheduling) (string, error) {
	logger.Info("Schedule a task at", s.Date.Format(time.RFC3339Nano))

	for _, publ := range s.Publishers {
		pub, ok := c.Pubs.Get(publ.Publisher)
		if !ok {
			return "", fmt.Errorf("Publisher %s does not exist", publ.Publisher)
		}

		err := publisher.CheckConfig(pub, publ.Settings)
		if err != nil {
			return "", err
		}
	}

	err := c.Queue.Add(s)
	if err != nil {
		return "", err
	}

	return s.ID, nil
}

// GetTasks returns tasks from db
func (c *Ctrl) GetTasks(start time.Time, end time.Time) []*utils.Scheduling {
	return c.DB.GetTasks(start, end)
}

// GetPublishers returns publishers
func (c *Ctrl) GetPublishers() interface{} {
	return c.Pubs.GetAvailable()
}

// RemoveTask removes a task according to id
func (c *Ctrl) RemoveTask(id string) error {
	return c.Queue.Remove(id)
}
