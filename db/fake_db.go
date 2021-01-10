package db

import (
	"jpb/scheduler/utils"
	"time"
)

// FakeDb struct
type FakeDb struct {
}

// NewFakedb returns fakedb driver
func NewFakedb() *FakeDb {
	return &FakeDb{}
}

func (f *FakeDb) GetTasks(start time.Time, end time.Time) []*utils.Scheduling {
	return []*utils.Scheduling{}
}

func (f *FakeDb) GetTasksToDo(end time.Time) []*utils.Scheduling {
	return []*utils.Scheduling{}
}

func (f *FakeDb) StoreTask(*utils.Scheduling) error {
	return nil
}

func (f *FakeDb) AckTask(string) error {
	return nil
}

func (f *FakeDb) RemoveTask(string) error {
	return nil
}
