package db

import "jpb/scheduler/task"

// Taskdb is taskdb
type Taskdb interface {
	GetFirstTasks(int) []*task.Task
	StoreTask(*task.Task) error
	AckTask(string) error
	RemoveTask(string) error
}
