package controller

import (
	"fmt"
	"jpb/scheduler/config"
	"jpb/scheduler/publisher"
	"jpb/scheduler/taskqueue"
	"net/http"
	"time"
)

// Ctrl represents a scheduler controller
type Ctrl struct {
	queue     *taskqueue.TaskQueue
	publisher *publisher.HTTPPublisher
}

// singleton
var ctrl *Ctrl = nil

// New creates a scheduler controller
func New() *Ctrl {
	c := config.Get()

	if ctrl == nil {
		ctrl = &Ctrl{
			queue:     taskqueue.New(c.MaxQueueLen),
			publisher: publisher.NewHTTPPublisher(http.MethodPost, "http://127.0.0.1:3003/", "{}"),
		}
	}
	return ctrl
}

// Schedule schedules a task
func (c *Ctrl) Schedule(date *time.Time) (string, error) {
	fmt.Println("Schedule a task at", date.Format(time.RFC3339))
	pub := c.publish
	return c.queue.Add(date, &pub)
}

func (c *Ctrl) publish(id string) {
	fmt.Println(fmt.Sprintf("on publish at %s", time.Now().Format(time.RFC3339)))
	err := c.publisher.Publish()
	if err != nil {
		fmt.Println(err.Error())
	}
}
