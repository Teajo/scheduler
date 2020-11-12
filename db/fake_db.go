package db

import "jpb/scheduler/task"

// Fakedb represents a fake db
type Fakedb struct {
}

// GetFirstTasks retrieve nb first tasks
func (f *Fakedb) GetFirstTasks(nb int) []*task.Task {
	return []*task.Task{}
}

// StoreTask blabla
func (f *Fakedb) StoreTask(*task.Task) error {
	return nil
}

// AckTask blabla
func (f *Fakedb) AckTask(string) error {
	return nil
}

// RemoveTask blabla
func (f *Fakedb) RemoveTask(string) error {
	return nil
}
