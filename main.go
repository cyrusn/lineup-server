package main

import (
	"github.com/cyrusn/lineup-server/cmd"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cmd.Execute()
}
