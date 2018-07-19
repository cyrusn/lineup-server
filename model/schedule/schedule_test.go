package schedule_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/cyrusn/goTestHelper"
	"github.com/cyrusn/lineup-system/database"
	"github.com/cyrusn/lineup-system/schedule"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbPath = "../test/test.db"
)

type mySchedule struct {
	classCode string
	classNo   int
}

var mySchedules = []mySchedule{
	mySchedule{"3A", 12},
	mySchedule{"3C", 8},
	mySchedule{"3B", 2},
	mySchedule{"3C", 5},
	mySchedule{"3D", 1},
	mySchedule{"3B", 4},
	mySchedule{"3A", 9},
	mySchedule{"3D", 10},
	mySchedule{"3C", 14},
	mySchedule{"3C", 6},
	mySchedule{"3D", 2},
}

var scheduleDB schedule.DB

func init() {
	if err := os.Remove(dbPath); err != nil {
		fmt.Println(err)
	}

	if err := database.CreateDBFile(dbPath, true); err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	scheduleDB = schedule.DB{db}
}

// TestMain is main test function
func TestMain(t *testing.T) {
	t.Run("insert schedule", testInsertUser)
	t.Run("insert duplicated schedule", testDuplicatedSchedule)
	t.Run("Toggle all 3A isNotified", testToggleIsNotified)
	t.Run("update all 3B priority to 2", testUpdatePriority)
	t.Run("Toggle all 3B isMeeting", testToggleIsMeeting)
	t.Run("delete all 3C schedule", testDelete)
	t.Run("Toggle all 3D isComplete", testToggleIsComplete)
	t.Run("query by class", testSelectByClassCode)
}

var testInsertUser = func(t *testing.T) {

	for _, s := range mySchedules {
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

var testDelete = func(t *testing.T) {
	for _, s := range mySchedules {
		if s.classCode == "3C" {
			err := scheduleDB.Delete(s.classCode, s.classNo)
			assert.OK(t, err)
		}
	}
}

var testUpdatePriority = func(t *testing.T) {
	for _, s := range mySchedules {
		if s.classCode == "3B" {
			err := scheduleDB.UpdatePriority(s.classCode, s.classNo, 2)
			assert.OK(t, err)
		}
	}
}

var testToggleIsNotified = func(t *testing.T) {
	for _, s := range mySchedules {
		if s.classCode == "3A" {
			err := scheduleDB.ToggleIsNotified(s.classCode, s.classNo)
			assert.OK(t, err)
		}
	}
}

var testToggleIsMeeting = func(t *testing.T) {
	for _, s := range mySchedules {
		if s.classCode == "3B" {
			err := scheduleDB.ToggleIsMeeting(s.classCode, s.classNo)
			assert.OK(t, err)
		}
	}
}

var testToggleIsComplete = func(t *testing.T) {
	for _, s := range mySchedules {
		if s.classCode == "3D" {
			err := scheduleDB.ToggleIsComplete(s.classCode, s.classNo)
			assert.OK(t, err)
		}
	}
}

func selectByClass(t *testing.T, classCode string) {
	mySchedules, err := scheduleDB.SelectByClassCode(classCode)
	assert.OK(t, err)

	for _, s := range mySchedules {

		if s.ClassCode == "3A" {
			assert.Equal(s.IsNotified, true, t)
		}
		if s.ClassCode == "3B" {
			assert.Equal(s.Priority, 2, t)
			assert.Equal(s.IsMeeting, true, t)
		}
		if s.ClassCode == "3C" {
			assert.Panic("No 3C", t, func() {
				panic("3C should all be removed")
			})
		}
		if s.ClassCode == "3D" {
			assert.Equal(s.IsComplete, true, t)
		}
	}
}
