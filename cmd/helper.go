package cmd

import (
	"log"
	"os"
)

func checkPathExist(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf(`"%s" doesn't exist`, path)
	}
}
