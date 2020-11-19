package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"jpb/scheduler/logger"
	"jpb/scheduler/utils"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3" // sqlite3
)

type sqlite3db struct {
	conn *sql.DB
	mu   sync.Mutex
}

func newSqlite3() *sqlite3db {
	conn, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		panic(err.Error())
	}

	query, err := conn.Prepare(`CREATE TABLE IF NOT EXISTS tasks (id INTEGER PRIMARY KEY, uid TEXT NOT NULL UNIQUE, date DATETIME, publishers TEXT, settings TEXT, done INTEGER DEFAULT 0)`)
	if err != nil {
		panic(err.Error())
	}

	_, err = query.Exec()
	if err != nil {
		panic(err.Error())
	}

	return &sqlite3db{
		conn: conn,
	}
}

func (f *sqlite3db) GetTasks(start time.Time, end time.Time) []*utils.Scheduling {
	f.mu.Lock()
	defer f.mu.Unlock()

	logger.Info(fmt.Sprintf("Get all tasks which end before %s", end.String()))
	tasks := []*utils.Scheduling{}
	rows, err := f.conn.Query("SELECT uid, date, publishers, settings, done FROM tasks WHERE datetime(date) >= datetime(?) AND datetime(date) <= datetime(?) ORDER BY date ASC", start, end)
	if err != nil {
		panic(err)
	}

	var uid string
	var date string
	var publishers []string
	var settings string
	var done bool

	for rows.Next() {
		rows.Scan(&uid, &date, &publishers, &settings, &done)
		d, _ := time.Parse(time.RFC3339Nano, date)
		tasks = append(tasks, utils.NewSchedulingWithID(uid, d, publishers, jsonStringToMap(settings), done))
	}

	return tasks
}

func (f *sqlite3db) GetTasksToDo(start time.Time, end time.Time) []*utils.Scheduling {
	f.mu.Lock()
	defer f.mu.Unlock()

	logger.Info(fmt.Sprintf("Get all tasks which end before %s", end.String()))
	tasks := []*utils.Scheduling{}
	rows, err := f.conn.Query("SELECT uid, date, publishers, settings FROM tasks WHERE datetime(date) >= datetime(?) AND datetime(date) <= datetime(?) AND done = 0 ORDER BY date ASC", start, end)
	if err != nil {
		panic(err)
	}

	var uid string
	var date string
	var publishers []string
	var settings string

	for rows.Next() {
		rows.Scan(&uid, &date, &publishers, &settings)
		d, _ := time.Parse(time.RFC3339Nano, date)
		tasks = append(tasks, utils.NewSchedulingWithID(uid, d, publishers, jsonStringToMap(settings), false))
	}

	return tasks
}

func (f *sqlite3db) StoreTask(t *utils.Scheduling) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	query, err := f.conn.Prepare("INSERT INTO tasks (uid, date, publishers, settings) VALUES (?, ?, ?, ?)")
	_, err = query.Exec(t.ID, t.Date.Format(time.RFC3339Nano), t.Publishers, mapToJSONString(t.Settings))
	return err
}

func (f *sqlite3db) AckTask(id string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	query, err := f.conn.Prepare("UPDATE tasks SET done = 1 WHERE uid = ?")
	_, err = query.Exec(id)
	return err
}

func (f *sqlite3db) RemoveTask(id string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	query, err := f.conn.Prepare("UPDATE tasks SET done = 1 WHERE uid = ?")
	_, err = query.Exec(id)
	return err
}

func jsonStringToMap(jsonstr string) map[string]string {
	x := make(map[string]string)
	err := json.Unmarshal([]byte(jsonstr), &x)
	if err != nil {
		logger.Error("Error: %s\n", err)
		return make(map[string]string)
	}
	return x
}

func mapToJSONString(data map[string]string) string {
	str, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(str)
}
