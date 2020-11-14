package db

import (
	"jpb/scheduler/utils"
	"time"
)

// fakedb represents a fake db
type fakedb struct {
}

// NewFakedb creates a new fake db
func newFakedb() *fakedb {
	return &fakedb{}
}

func (f *fakedb) GetTasks(lastuid string, nb int, first time.Time) []*utils.Scheduling {
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
