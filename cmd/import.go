package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/cyrusn/lineup-system/model/auth"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import user in database",
	Run: func(cmd *cobra.Command, args []string) {
		paths := []string{dbPath, userJSONPath}
		checkPathExist(paths)

		file, err := ioutil.ReadFile(userJSONPath)
		if err != nil {
			fmt.Printf("File error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Parsing:", userJSONPath)

		var credentials []auth.Credential

		if err := json.Unmarshal(file, &credentials); err != nil {
			fmt.Printf("File error: %v\n", err)
			os.Exit(1)
		}

		sqldb, err := sql.Open("sqlite3", dbPath)

		if err != nil {
			log.Fatal(err)
		}

		db := &auth.DB{sqldb, &secret}

		for _, c := range credentials {
			err := db.Insert(c.UserAlias, c.Password, c.Role)
			if err != nil {
				fmt.Printf("Import error: %v\n", err)
				os.Exit(1)
			}
		}
		fmt.Println("users are imported")
	},
}
