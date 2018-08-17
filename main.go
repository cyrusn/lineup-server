package main

import (
	"github.com/cyrusn/lineup-server/cmd"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cmd.Execute()
}
