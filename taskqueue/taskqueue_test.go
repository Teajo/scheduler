package taskqueue

import (
	"jpb/scheduler/db"
	"jpb/scheduler/utils"
	"testing"
	"time"
)

var layout = "2006-01-02T15:04:05.000Z"

func TestInstantiateQueue(t *testing.T) {
	c := make(chan *utils.Scheduling)
	queue := New(db.Getdb("fake"), 3, c)
	go queue.Listen()
	if queue.Len() != 0 {
		t.Error("queue length should be 0")
	}
}

func TestAddTask(t *testing.T) {
	c := make(chan *utils.Scheduling)
	queue := New(db.Getdb("fake"), 1, c)
	go queue.Listen()

	now, _ := time.Parse(layout, "2020-11-12T12:45:26.371Z")
	queue.Add(utils.NewScheduling(now, "test", make(map[string]string)))

	if queue.Len() != 1 {
		t.Error("element not added")
	}
}

func TestRemoveTask(t *testing.T) {
	c := make(chan *utils.Scheduling)
	queue := New(db.Getdb("fake"), 3, c)
	go queue.Listen()

	now, _ := time.Parse(layout, "2020-11-12T12:45:26.371Z")
	id, _ := queue.Add(utils.NewScheduling(now, "test", make(map[string]string)))
	queue.Remove(id)
	if queue.Len() != 0 {
		t.Error("element not removed")
	}
}

func TestDoneTask(t *testing.T) {
	c := make(chan *utils.Scheduling)
	queue := New(db.Getdb("fake"), 3, c)
	go queue.Listen()

	now, _ := time.Parse(layout, "2019-11-12T12:45:26.371Z")
	queue.Add(utils.NewScheduling(now, "test", make(map[string]string)))

	select {
	case <-c:
		<-time.After(50 * time.Millisecond)
		if queue.Len() != 0 {
			t.Error("task done but not removed")
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("task not done and not removed")
	}
}

func TestGracefullStop(t *testing.T) {
	c := make(chan *utils.Scheduling)
	queue := New(db.Getdb("fake"), 2, c)
	go queue.Listen()

	now, _ := time.Parse(layout, "2020-11-12T12:45:26.371Z")
	queue.Add(utils.NewScheduling(now, "test", make(map[string]string)))

	err := queue.Stop()
	if err != nil {
		t.Error("should stop gracefully")
	}
}
