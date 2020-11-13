package db

import (
	"fmt"
	"jpb/scheduler/utils"
)

// Taskdb is taskdb
type Taskdb interface {
	GetFirstTasks(int) []*utils.Scheduling
	StoreTask(*utils.Scheduling) error
	AckTask(string) error
	RemoveTask(string) error
}

// Getdb returns a db according to driver
func Getdb(driver string) Taskdb {
	switch driver {
	case "fake":
		return newFakedb()
	case "sqlite3":
		return newSqlite3()
	default:
		panic(fmt.Sprintf("driver %s not handled", driver))
	}
}
