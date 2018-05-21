package model_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/cyrusn/lineup-system/model"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbPath = "../database/test.db"
)

var scheudleDB model.ScheduleDB

func init() {
	if err := os.Remove(dbPath); err != nil {
		log.Print(err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	scheudleDB = model.ScheduleDB{db}
}

// TestMain is main test function
func TestMain(t *testing.T) {

	t.Run("create database file", testCreateDB)

	t.Run("insert user", testInsertUser)

	t.Run("query by class", testSelectByClassCode)
}

var testCreateDB = func(t *testing.T) {
	if err := model.CreateDBFile(dbPath, true); err != nil {
		t.Fatal(err)
	}
}

var testInsertUser = func(t *testing.T) {
	type schedule struct {
		classCode string
		classNo   int
	}

	var schedules = []schedule{
		schedule{"3A", 12},
		schedule{"3C", 8},
		schedule{"3B", 2},
		schedule{"3B", 4},
		schedule{"3C", 5},
		schedule{"3A", 9},
		schedule{"3C", 10},
	}

	for _, s := range schedules {
		if err := scheudleDB.Insert(s.classCode, s.classNo); err != nil {
			t.Fatal(err)
		}
	}
}

var testSelectByClassCode = func(t *testing.T) {
	var classCodes = []string{"3A", "3B", "3C"}

	for _, c := range classCodes {
		selectByClass(t, c)
		fmt.Println("======")
	}
}

func selectByClass(t *testing.T, classCode string) {
	schedules, err := scheudleDB.SelectByClassCode(classCode)
	if err != nil {
		t.Fatal(err)
	}

	for _, s := range schedules {
		fmt.Println(s)
	}
}
