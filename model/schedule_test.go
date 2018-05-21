package model_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/cyrusn/goTestHelper"
	"github.com/cyrusn/lineup-system/model"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbPath = "../database/test.db"
)

type schedule struct {
	classCode string
	classNo   int
}

var schedules = []schedule{
	schedule{"3A", 12},
	schedule{"3C", 8},
	schedule{"3B", 2},
	schedule{"3C", 5},
	schedule{"3D", 1},
	schedule{"3B", 4},
	schedule{"3A", 9},
	schedule{"3D", 10},
	schedule{"3C", 14},
	schedule{"3C", 6},
	schedule{"3D", 2},
}

var scheduleDB model.ScheduleDB

func init() {
	if err := os.Remove(dbPath); err != nil {
		log.Print(err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	scheduleDB = model.ScheduleDB{db}
}

// TestMain is main test function
func TestMain(t *testing.T) {

	t.Run("create database file", testCreateDB)
	t.Run("insert schedule", testInsertUser)
	t.Run("query by class", testSelectByClassCode)
	t.Run("insert duplicated schedule", testDuplicatedSchedule)
	t.Run("delete all 3C schedule", testDeleteSchedule)
	t.Run("update all 3B priority to it's index", testUpdatePriority)
	t.Run("Toggle all 3A isNotified", testToggleIsNotified)
	t.Run("Toggle all 3B isMeeting", testToggleIsMeeting)
	t.Run("Toggle all 3D isComplete", testToggleIsComplete)
	t.Run("query by class", testSelectByClassCode)
}

var testCreateDB = func(t *testing.T) {
	err := model.CreateDBFile(dbPath, true)
	assert.OK(t, err)
}

var testInsertUser = func(t *testing.T) {

	for _, s := range schedules {
		err := scheduleDB.Insert(s.classCode, s.classNo)
		assert.OK(t, err)
	}
}

var testDuplicatedSchedule = func(t *testing.T) {
	assert.Panic("Duplicated Schedule", t, func() {
		if err := scheduleDB.Insert("3A", 12); err != nil {
			panic(err)
		}
	})
}

var testSelectByClassCode = func(t *testing.T) {
	var classCodes = []string{"3A", "3B", "3C", "3D"}

	for _, c := range classCodes {
		selectByClass(t, c)
	}
}

var testDeleteSchedule = func(t *testing.T) {
	for _, s := range schedules {
		if s.classCode == "3C" {
			err := scheduleDB.DeleteSchedule(s.classCode, s.classNo)
			assert.OK(t, err)
		}
	}
}

var testUpdatePriority = func(t *testing.T) {
	for i, s := range schedules {
		if s.classCode == "3B" {
			err := scheduleDB.UpdatePriority(s.classCode, s.classNo, i)
			assert.OK(t, err)
		}
	}
}

var testToggleIsNotified = func(t *testing.T) {
	for _, s := range schedules {
		if s.classCode == "3A" {
			err := scheduleDB.ToggleIsNotified(s.classCode, s.classNo)
			assert.OK(t, err)
		}
	}
}

var testToggleIsMeeting = func(t *testing.T) {
	for _, s := range schedules {
		if s.classCode == "3B" {
			err := scheduleDB.ToggleIsMeeting(s.classCode, s.classNo)
			assert.OK(t, err)
		}
	}
}

var testToggleIsComplete = func(t *testing.T) {
	for _, s := range schedules {
		if s.classCode == "3D" {
			err := scheduleDB.ToggleIsComplete(s.classCode, s.classNo)
			assert.OK(t, err)
		}
	}
}

func selectByClass(t *testing.T, classCode string) {
	schedules, err := scheduleDB.SelectByClassCode(classCode)
	assert.OK(t, err)

	for _, s := range schedules {
		fmt.Println(s)
	}
}
