package taskqueue

import (
	"testing"
	"time"
)

var layout = "2006-01-02T15:04:05.000Z"

func TestInstantiateQueue(t *testing.T) {
	queue := New(3)
	if queue.Len() != 0 {
		t.Error("queue length should be 0")
	}
}

func TestAddTask(t *testing.T) {
	queue := New(3)
	now, _ := time.Parse(layout, "2020-11-12T12:45:26.371Z")
	queue.Add(&now, nil)

	if queue.Len() != 1 {
		t.Error("element not added")
	}
}

func TestRemoveTask(t *testing.T) {
	queue := New(3)
	now, _ := time.Parse(layout, "2020-11-12T12:45:26.371Z")
	id, _ := queue.Add(&now, nil)
	queue.Remove(id)
	if queue.Len() != 0 {
		t.Error("element not removed")
	}
}

func TestDoneTask(t *testing.T) {
	queue := New(3)
	now, _ := time.Parse(layout, "2019-11-12T12:45:26.371Z")
	queue.Add(&now, nil)

	<-time.After(100 * time.Millisecond) // give time to scheduler to delete task as this is asynchronous
	if queue.Len() != 0 {
		t.Error("element not done")
	}
}

func TestAddTaskWhenQueueIsFull(t *testing.T) {
	queue := New(2)
	date1, _ := time.Parse(layout, "2021-11-12T12:45:26.371Z")
	date2, _ := time.Parse(layout, "2021-11-12T11:45:26.371Z")
	date3, _ := time.Parse(layout, "2021-11-12T10:45:26.371Z")
	date4, _ := time.Parse(layout, "2021-11-12T13:45:26.371Z")

	_, err := queue.Add(&date1, nil)
	if err != nil {
		t.Error("error has been triggered", err.Error())
	}

	_, err = queue.Add(&date2, nil)
	if err != nil {
		t.Error("error has been triggered", err.Error())
	}

	_, err = queue.Add(&date3, nil)
	if err != nil {
		t.Error("error has been triggered", err.Error())
	}

	_, err = queue.Add(&date4, nil)
	if err == nil {
		t.Error("error not triggered")
	}

	for _, k := range queue.queue {
		t.Log(k)
	}

	if queue.Len() != 2 {
		t.Error("maxLen not respected")
	}

	if queue.last != &date2 {
		t.Error("last date is not correct")
	}
}

func TestGracefullStop(t *testing.T) {
	queue := New(2)
	err := queue.Stop()
	if err != nil {
		t.Error("should stop gracefully")
	}
}

func TestTodoFunc(t *testing.T) {
	queue := New(2)
	isOk := make(chan bool)
	date := time.Now().Add(100 * time.Millisecond)
	t.Log("should execute task at ", date)
	todo := func(id string) {
		t.Log("OK", id, time.Now().String())
		isOk <- true
	}
	queue.Add(&date, &todo)

	for {
		select {
		case <-isOk:
			if queue.Len() > 0 {
				t.Error("Queue should be empty")
			}
			return
		case <-time.After(200 * time.Millisecond):
			t.Error("Task should have been executed")
			return
		}
	}
}
