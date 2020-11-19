package db

import (
	"jpb/scheduler/utils"
	"time"
)

type fakedb struct {
}

func newFakedb() *fakedb {
	return &fakedb{}
}

func (f *fakedb) GetTasks(start time.Time, end time.Time) []*utils.Scheduling {
	return []*utils.Scheduling{}
}

func (f *fakedb) GetTasksToDo(start time.Time, end time.Time) []*utils.Scheduling {
	return []*utils.Scheduling{}
}

func (f *fakedb) StoreTask(*utils.Scheduling) error {
	return nil
}

func (f *fakedb) AckTask(string) error {
	return nil
}

func (f *fakedb) RemoveTask(string) error {
	return nil
}
