package model

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// DBPath is default DB path
const DBPath = "../database/test.db"

// TestMain is the main function for test program
func TestMain(t *testing.T) {
	t.Run("Create DB File", testCreate)
}

func testCreate(t *testing.T) {
	if err := CreateDBFile(DBPath, true); err != nil {
		t.Fatal(err)
	}
}
