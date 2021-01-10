package db

import (
	"fmt"
	"jpb/scheduler/utils"
	"time"
)

// Taskdb is taskdb
type Taskdb interface {
	// Get tasks scheduled before provided end date
	GetTasks(start time.Time, end time.Time) []*utils.Scheduling

	// Returns tasks to do
	GetTasksToDo(end time.Time) []*utils.Scheduling

	// Store task scheduling
	StoreTask(*utils.Scheduling) error

	// Ack task by ID
	AckTask(string) error

	// Remove a task by ID
	RemoveTask(string) error
}

// Getdb returns a db according to driver
func Getdb(driver string) Taskdb {
	switch driver {
	case "fake":
		return NewFakedb()
	case "sqlite3":
		return NewSqlite3()
	default:
		panic(fmt.Sprintf("driver %s not handled", driver))
	}
}
