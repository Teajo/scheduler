package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"jpb/scheduler/utils"
	"time"

	_ "github.com/mattn/go-sqlite3" // sqlite3
)

type sqlite3db struct {
	conn *sql.DB
}

func newSqlite3() *sqlite3db {
	conn, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		panic(err.Error())
	}

	query, err := conn.Prepare("CREATE TABLE IF NOT EXISTS tasks (id INTEGER PRIMARY KEY, uid TEXT NOT NULL UNIQUE, date DATETIME, publisher TEXT, settings TEXT)")
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

// GetFirstTasks retrieve nb first tasks
func (f *sqlite3db) GetFirstTasks(nb int) []*utils.Scheduling {
	tasks := []*utils.Scheduling{}
	rows, _ := f.conn.Query("SELECT * FROM tasks")

	var id int
	var uid string
	var date string
	var publisher string
	var settings string

	for rows.Next() {
		rows.Scan(&id, &uid, &date, &publisher, &settings)
		fmt.Println(uid + " " + date)
		d, _ := time.Parse(time.RFC3339Nano, date)
		tasks = append(tasks, utils.NewSchedulingWithID(uid, d, publisher, jsonStringToMap(settings)))
	}

	return tasks
}

// StoreTask blabla
func (f *sqlite3db) StoreTask(t *utils.Scheduling) error {
	query, err := f.conn.Prepare("INSERT INTO tasks (uid, date, publisher, settings) VALUES (?, ?, ?, ?)")
	_, err = query.Exec(t.ID, t.Date.Format(time.RFC3339Nano), t.Publisher, mapToJSONString(t.Settings))
	return err
}

// AckTask blabla
func (f *sqlite3db) AckTask(string) error {
	return nil
}

// RemoveTask blabla
func (f *sqlite3db) RemoveTask(id string) error {
	query, err := f.conn.Prepare("DELETE FROM tasks WHERE uid = ?")
	_, err = query.Exec(id)
	return err
}

func jsonStringToMap(jsonstr string) map[string]string {
	x := make(map[string]string)
	err := json.Unmarshal([]byte(jsonstr), &x)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
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
