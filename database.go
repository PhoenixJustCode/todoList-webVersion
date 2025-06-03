
package main

import (
    "database/sql"
    "log"

    _ "github.com/lib/pq"
)


type Task struct {
	ID   int64  `json:"id"`
	Task string `json:"task"`
	Days int16  `json:"days"`
}

type DB struct {
	Conn *sql.DB
}

func NewDB(dataSourceName string) (*DB, error) {
	dbConn, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err := dbConn.Ping(); err != nil {
		return nil, err
	}

	return &DB{Conn: dbConn}, nil
}

func (db *DB) Close() {
	db.Conn.Close()
}


func (db *DB) GetTasksByDays(days int) ([]Task, error) {
	rows, err := db.Conn.Query("SELECT id, task, days FROM taskinf WHERE days = $1 ORDER BY id", days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Task, &t.Days)
		if err != nil {
			log.Println("Scan error:", err)
			continue
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}


func (db *DB) GetAllTasks() ([]Task, error) {
	rows, err := db.Conn.Query("SELECT id, task, days FROM taskinf ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Task, &t.Days)
		if err != nil {
			log.Println("Scan error:", err)
			continue
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (db *DB) InsertTask(t Task) error {
	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO taskinf(task, days) VALUES ($1, $2)", t.Task, t.Days)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO logs(entity, action) VALUES ($1, $2)", "task", "created")
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (db *DB) DeleteTask(id int64) error {
	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM taskinf WHERE id = $1", id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO logs(entity, action) VALUES ($1, $2)", "task", "deleted")
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (db *DB) UpdateTask(id int64, newTask Task) error {
	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE taskinf SET task = $1, days = $2 WHERE id = $3", newTask.Task, newTask.Days, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO logs(entity, action) VALUES ($1, $2)", "task", "updated")
	if err != nil {
		return err
	}

	return tx.Commit()
}
