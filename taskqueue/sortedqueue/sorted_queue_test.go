package sortedqueue

import (
	"jpb/scheduler/task"
	"jpb/scheduler/utils"
	"testing"
	"time"
)

var layout = "2006-01-02T15:04:05.000Z"

func TestSortedQueueInstantiation(t *testing.T) {
	queue := New(3)
	if queue.Len() != 0 {
		t.Error("queue length should be 0")
	}
}

func TestSortedQueueInsert(t *testing.T) {
	queue := New(3)
	ta := task.New(utils.NewScheduling(time.Now(), "test", make(map[string]string)), func(*utils.Scheduling) {})
	queue.Insert(ta)
	if queue.Len() != 1 {
		t.Error("queue length should be 1")
	}
}

func TestSortedQueueRemove(t *testing.T) {
	queue := New(3)
	ta := task.New(utils.NewScheduling(time.Now(), "test", make(map[string]string)), func(*utils.Scheduling) {})
	queue.Insert(ta)
	queue.Remove(ta)
	if queue.Len() != 0 {
		t.Error("queue length should be 0")
	}
}

func TestSortedQueueRemoveByID(t *testing.T) {
	queue := New(3)
	ta := task.New(utils.NewScheduling(time.Now(), "test", make(map[string]string)), func(*utils.Scheduling) {})
	queue.Insert(ta)
	queue.RemoveByID(ta.ID)
	if queue.Len() != 0 {
		t.Error("queue length should be 0")
	}
}

func TestSortedQueueGet(t *testing.T) {
	queue := New(3)
	ta := task.New(utils.NewScheduling(time.Now(), "test", make(map[string]string)), func(*utils.Scheduling) {})
	queue.Insert(ta)
	q := queue.Get()
	if q[0] != ta {
		t.Error("queue elem must be task")
	}
}

func TestSortedQueueInsertInExcess(t *testing.T) {
	date1, _ := time.Parse(layout, "2019-11-12T12:45:26.371Z")
	date2, _ := time.Parse(layout, "2020-11-12T12:45:26.371Z")

	queue := New(1)
	ta1 := task.New(utils.NewScheduling(date1, "test", make(map[string]string)), func(*utils.Scheduling) {})
	queue.Insert(ta1)
	ta2 := task.New(utils.NewScheduling(date2, "test2", make(map[string]string)), func(*utils.Scheduling) {})
	excess := queue.Insert(ta2)

	if excess[0] != ta2 {
		t.Error("queue must reject last task")
	}
}
