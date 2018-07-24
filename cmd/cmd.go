package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	port                 string
	dbPath               string
	isOverwrite          bool
	staticFolderLocation string
	userJSONPath         string
	lifeTime             int64

	rootCmd = &cobra.Command{
		Use:   "lineup",
		Short: "Welcome to Line-Up System Backend Server",
		Run:   rootCmdStartupFunc,
	}
)

func init() {
	cmds := []*cobra.Command{versionCmd, serveCmd, createCmd, importCmd}

	for _, cmd := range cmds {
		rootCmd.AddCommand(cmd)
	}

	serveCmd.PersistentFlags().StringVarP(&port, "port", "p", ":5000", "Port value")
	serveCmd.PersistentFlags().StringVarP(&dbPath, "location", "l", "./test/test.db", "Location of sqlite3 database file")
	serveCmd.PersistentFlags().StringVarP(&staticFolderLocation, "static", "s", "./public", "Location of static folder for serving")
	serveCmd.PersistentFlags().Int64VarP(&lifeTime, "time", "t", 30, "update the life time of jwt token")

	createCmd.PersistentFlags().StringVarP(&dbPath, "location", "l", "./test/test.db", "Location of sqlite3 database file")
	createCmd.PersistentFlags().BoolVarP(&isOverwrite, "overwrite", "o", false, "Overwrite database if database location exist")

	importCmd.PersistentFlags().StringVarP(&userJSONPath, "import", "i", "./test/user.json", "path to user.json file")
	importCmd.PersistentFlags().StringVarP(&dbPath, "location", "l", "./test/test.db", "Location of sqlite3 database file")
}

// Execute run all cmds
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func rootCmdStartupFunc(cmd *cobra.Command, args []string) {
	cmd.Help()
}
