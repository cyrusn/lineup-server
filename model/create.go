package model

import (
	"database/sql"
	"fmt"
	"os"
)

const userSchema = `
CREATE TABLE IF NOT EXISTS USER (
  username TEXT PRIMARY KEY,
  password TEXT NOT NULL,
  name TEXT NOT NULL,
  cname TEXT
)
`

const scheduleSchema = `
CREATE TABLE IF NOT EXISTS SCHEDULE (
  classcode TEXT NOT NULL,
  classno INT NOT NULL,
  arrived_at TIMESTAMP,
  priority INT NOT NULL,
  is_notified BOOL NOT NULL,
  is_meeting BOOL NOT NULL,
  is_complete BOOL NOT NULL,
  CONSTRAINT unique_student UNIQUE (classcode, classno)
);`

// CreateDBFile create new database and write `schedule` schema on it.
// To create long-lived db for server, use `sql.Open()` to open an database,
// and developer should check if the database is already existed.
func CreateDBFile(dbPath string, isOverWrite bool) error {
	exist, err := isFileExist(dbPath)
	if err != nil {
		return err
	}

	if exist && !isOverWrite {
		return nil
	}

	if _, err := os.Create(dbPath); err != nil {
		return err
	}

	if err := openDBAndCreateTable(dbPath); err != nil {
		return err
	}

	return nil
}

// openDBAndCreateTable open existed database and create table
func openDBAndCreateTable(dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	defer db.Close()

	if err != nil {
		return err
	}
	if err := createTable(db, dbPath); err != nil {
		return err
	}

	return nil
}

// createTable create table for the schema of ScheduleSchema
func createTable(db *sql.DB, dbPath string) error {
	if _, err := db.Exec(scheduleSchema); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// isFileExist check if file exist
func isFileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		// file exist
		return true, nil

	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
