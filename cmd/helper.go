package cmd

import (
	"log"
	"os"
)

func checkPathExist(paths []string) {
	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Fatalf(`"%s" doesn't exist`, path)
		}
	}
}
