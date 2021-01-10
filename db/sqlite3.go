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

// Sqlite3db struct
type Sqlite3db struct {
	conn *sql.DB
	mu   sync.Mutex
}

// NewSqlite3 returns a sqlite3 driver
func NewSqlite3() *Sqlite3db {
	conn, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		panic(err.Error())
	}

	query, err := conn.Prepare(`CREATE TABLE IF NOT EXISTS tasks (id INTEGER PRIMARY KEY, uid TEXT NOT NULL UNIQUE, date DATETIME, publishers TEXT, done INTEGER DEFAULT 0)`)
	if err != nil {
		panic(err.Error())
	}

	_, err = query.Exec()
	if err != nil {
		panic(err.Error())
	}

	return &Sqlite3db{
		conn: conn,
	}
}

func (f *Sqlite3db) GetTasks(start time.Time, end time.Time) []*utils.Scheduling {
	f.mu.Lock()
	defer f.mu.Unlock()

	logger.Info(fmt.Sprintf("Get all tasks which end before %s and start after %s", end.String(), start.String()))
	tasks := []*utils.Scheduling{}
	rows, err := f.conn.Query("SELECT uid, date, publishers, done FROM tasks WHERE datetime(date) >= datetime(?) AND datetime(date) <= datetime(?) ORDER BY date ASC", start, end)
	if err != nil {
		panic(err)
	}

	var uid string
	var date string
	var publishers string
	var done bool

	for rows.Next() {
		var pubs []*utils.Publisher
		rows.Scan(&uid, &date, &publishers, &done)
		json.Unmarshal([]byte(publishers), &pubs)
		d, _ := time.Parse(time.RFC3339Nano, date)
		tasks = append(tasks, utils.NewSchedulingWithID(uid, d, pubs, done))
	}

	return tasks
}

func (f *Sqlite3db) GetTasksToDo(end time.Time) []*utils.Scheduling {
	f.mu.Lock()
	defer f.mu.Unlock()

	logger.Info(fmt.Sprintf("Get all tasks to do which end before %s", end.String()))
	tasks := []*utils.Scheduling{}
	rows, err := f.conn.Query("SELECT uid, date, publishers FROM tasks WHERE datetime(date) <= datetime(?) AND done = 0 ORDER BY date ASC", end)
	if err != nil {
		panic(err)
	}

	var uid string
	var date string
	var publishers string

	for rows.Next() {
		var pubs []*utils.Publisher
		rows.Scan(&uid, &date, &publishers)
		json.Unmarshal([]byte(publishers), &pubs)
		d, _ := time.Parse(time.RFC3339Nano, date)
		tasks = append(tasks, utils.NewSchedulingWithID(uid, d, pubs, false))
	}

	return tasks
}

func (f *Sqlite3db) StoreTask(t *utils.Scheduling) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	b, err := json.Marshal(t.Publishers)
	if err != nil {
		return err
	}

	query, err := f.conn.Prepare("INSERT INTO tasks (uid, date, publishers) VALUES (?, ?, ?)")
	_, err = query.Exec(t.ID, t.Date.Format(time.RFC3339Nano), string(b))
	return err
}

func (f *Sqlite3db) AckTask(id string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	query, err := f.conn.Prepare("UPDATE tasks SET done = 1 WHERE uid = ?")
	_, err = query.Exec(id)
	return err
}

func (f *Sqlite3db) RemoveTask(id string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	query, err := f.conn.Prepare("DELETE FROM tasks WHERE uid = ?")
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
