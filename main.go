package main

import (
	"github.com/cyrusn/lineup-system/cmd"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cmd.Execute()
}
