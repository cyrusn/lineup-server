package database_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/cyrusn/goTestHelper"
	"github.com/cyrusn/lineup-system/database"

	_ "github.com/mattn/go-sqlite3"
)

// dbPath is default DB path
const dbPath = "../test/test.db"

func init() {
	if err := os.Remove(dbPath); err != nil {
		fmt.Println(err)
	}
}

// TestMain is the main function for test program
func TestMain(t *testing.T) {
	t.Run("create database file", testCreateDB)
}

var testCreateDB = func(t *testing.T) {
	err := database.CreateDBFile(dbPath, true)
	assert.OK(t, err)
}
