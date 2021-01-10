package queue

import (
	"jpb/scheduler/task"
	"jpb/scheduler/utils"
	"testing"
	"time"
)

var layout = "2006-01-02T15:04:05.000Z"
var lastDate time.Time = time.Unix(1<<63-62135596801, 999999999)

func TestSortedQueueInstantiation(t *testing.T) {
	queue := New(3, lastDate)
	if queue.Len() != 0 {
		t.Error("queue length should be 0")
	}
}

func TestSortedQueueInsert(t *testing.T) {
	queue := New(3, lastDate)
	ta := task.New(utils.NewScheduling(time.Now(), make(map[string]string)))
	queue.Add(ta)
	if queue.Len() != 1 {
		t.Error("queue length should be 1")
	}
}

func TestSortedQueueRemove(t *testing.T) {
	queue := New(3, lastDate)
	ta := task.New(utils.NewScheduling(time.Now(), "test", make(map[string]string)))
	queue.Add(ta)
	queue.Remove(ta)
	if queue.Len() != 0 {
		t.Error("queue length should be 0")
	}
}

func TestSortedQueueRemoveByID(t *testing.T) {
	queue := New(3, lastDate)
	ta := task.New(utils.NewScheduling(time.Now(), "test", make(map[string]string)))
	queue.Add(ta)
	queue.RemoveByID(ta.ID)
	if queue.Len() != 0 {
		t.Error("queue length should be 0")
	}
}

func TestSortedQueueGet(t *testing.T) {
	queue := New(3, lastDate)
	ta := task.New(utils.NewScheduling(time.Now(), "test", make(map[string]string)))
	queue.Add(ta)
	q := queue.Get()
	if q[0] != ta {
		t.Error("queue elem must be the task")
	}
}
