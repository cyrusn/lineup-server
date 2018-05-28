package database

import (
	"database/sql"
	"os"
)

const (
	scheduleSchema = `
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

	authSchema = `
	  CREATE TABLE IF NOT EXISTS AUTHENTICATION (
	    useralias TEXT PRIMARY KEY,
	    password BLOB NOT NULL
	  );`
)

var schemas = []string{scheduleSchema, authSchema}

// CreateDBFile create new database and write `schedule` schema on it.
// To create long-lived db for server, use `sql.Open()` to open an database,
// and developer should check if the database is already existed.
func CreateDBFile(dbPath string, isOverWrite bool) error {
	exist, err := IsFileExist(dbPath)
	if err != nil {
		return err
	}

	if exist && !isOverWrite {
		return os.ErrExist
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
	if err != nil {
		return err
	}
	defer db.Close()

	if err := createTable(db, dbPath); err != nil {
		return err
	}

	return nil
}

// createTable create table for the schema of ScheduleSchema
func createTable(db *sql.DB, dbPath string) error {
	for _, schema := range schemas {
		if _, err := db.Exec(schema); err != nil {
			return err
		}
	}
	return nil
}

// IsFileExist check if file exist
func IsFileExist(path string) (bool, error) {
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
