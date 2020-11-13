package db

import (
	"jpb/scheduler/utils"
)

// fakedb represents a fake db
type fakedb struct {
}

// NewFakedb creates a new fake db
func newFakedb() *fakedb {
	return &fakedb{}
}

// GetFirstTasks retrieve nb first tasks
func (f *fakedb) GetFirstTasks(nb int) []*utils.Scheduling {
	return []*utils.Scheduling{}
}

// StoreTask blabla
func (f *fakedb) StoreTask(*utils.Scheduling) error {
	return nil
}

// AckTask blabla
func (f *fakedb) AckTask(string) error {
	return nil
}

// RemoveTask blabla
func (f *fakedb) RemoveTask(string) error {
	return nil
}
