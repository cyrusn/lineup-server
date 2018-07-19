package cmd

import (
	"fmt"
	"log"
	"os"

	auth_helper "github.com/cyrusn/goJWTAuthHelper"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	contextKeyName = "authClaim"
	jwtKeyName     = "jwt"
	roleKeyName    = "Role"
)

var (
	cfgFile              string
	port                 string
	dbPath               string
	isOverwrite          bool
	staticFolderLocation string
	userJSONPath         string
	lifeTime             int64
	privateKey           string
	secret               auth_helper.Secret
)

func initConfig() {
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}

func initSecret() {
	secret = auth_helper.New(
		contextKeyName, jwtKeyName, roleKeyName, []byte(privateKey),
	)
}

func initVar() {
	cfgFile = viper.GetString("config")
	port = viper.GetString("port")
	dbPath = viper.GetString("location")
	isOverwrite = viper.GetBool("overwrite")
	staticFolderLocation = viper.GetString("static")
	userJSONPath = viper.GetString("import")
	lifeTime = viper.GetInt64("time")
	privateKey = viper.GetString("key")
}

func init() {
	cobra.OnInitialize(initConfig, initVar, initSecret)
}

// Execute run all cmds
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
