package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/cyrusn/lineup-server/model/auth"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import user in database",
	Run: func(cmd *cobra.Command, args []string) {
		checkPathExist(userJSONPath)

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

		sqldb, err := sql.Open("mysql", dsn)

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

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.PersistentFlags().StringVarP(
		&userJSONPath,
		"import",
		"i",
		DEFAULT_IMPORT_PATH,
		`path of user.json file
The schema of the json file should be as follow.
[{"userAlias": "user1", "password": "password1", "role": "teacher"}, ... ]
`,
	)
	viper.BindPFlag("import", importCmd.PersistentFlags().Lookup("import"))

}
