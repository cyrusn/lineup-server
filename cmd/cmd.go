package cmd

import (
	"log"

	auth_helper "github.com/cyrusn/goJWTAuthHelper"

	"github.com/spf13/cobra"
)

const (
	contextKeyName = "authClaim"
	jwtKeyName     = "jwt"
	roleKeyName    = "Role"
)

var (
	port                 string
	dbPath               string
	isOverwrite          bool
	staticFolderLocation string
	userJSONPath         string
	lifeTime             int64
	privateKey           string
	secret               auth_helper.Secret

	rootCmd = &cobra.Command{
		Use:   "lineup",
		Short: "Welcome to Line-Up System Backend Server",
		Run:   rootCmdStartupFunc,
	}
)

func init() {
	cobra.OnInitialize(func() {
		secret = auth_helper.New(
			contextKeyName, jwtKeyName, roleKeyName, []byte(privateKey),
		)
	})

	cmds := []*cobra.Command{versionCmd, serveCmd, createCmd, importCmd}

	for _, cmd := range cmds {
		rootCmd.AddCommand(cmd)
	}

	rootCmd.PersistentFlags().StringVarP(
		&privateKey,
		"key",
		"k",
		"skill-vein-planet-neigh-envoi",
		"change the private key for authentication on jwt",
	)
	rootCmd.PersistentFlags().StringVarP(
		&dbPath,
		"location",
		"l",
		"./test/test.db",
		"location of sqlite3 database file",
	)
	serveCmd.PersistentFlags().StringVarP(
		&port,
		"port",
		"p",
		":5000",
		"port value",
	)
	serveCmd.PersistentFlags().StringVarP(
		&staticFolderLocation,
		"static",
		"s",
		"./public",
		"location of static folder for serving",
	)
	serveCmd.PersistentFlags().Int64VarP(
		&lifeTime,
		"time",
		"t",
		30,
		"update the life time (minutes) of jwt",
	)
	createCmd.PersistentFlags().BoolVarP(
		&isOverwrite,
		"overwrite",
		"o",
		false,
		"overwrite database if database location exist",
	)
	importCmd.PersistentFlags().StringVarP(
		&userJSONPath,
		"import",
		"i",
		"./test/user.json",
		`path of user.json file
The schema of the json file should be as follow.
[{"userAlias": "user1", "password": "password1", "role": "teacher"}, ... ]
`,
	)
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
