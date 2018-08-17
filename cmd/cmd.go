package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	auth_helper "github.com/cyrusn/goJWTAuthHelper"
	"github.com/go-sql-driver/mysql"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DEFAULT_CONFIG_PATH   = "./config.yaml"
	DEFAULT_DSN           = "root@/lineupTestDB"
	DEFAULT_PRIVATE_KEY   = "skill-vein-planet-neigh-envoi"
	DEFAULT_OVERWRITE     = false
	DEFAULT_IMPORT_PATH   = "./data/user.json"
	DEFAULT_PUBLIC_FOLDER = "./public"
	DEFAULT_PORT          = ":5000"

	CONTEXT_KEY_NAME = "authClaim"
	JWT_KEY_NAME     = "jwt"
	ROLE_KEY_NAME    = "Role"
)

var (
	cfgFile              string
	port                 string
	dsn                  string
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
		CONTEXT_KEY_NAME, JWT_KEY_NAME, ROLE_KEY_NAME, []byte(privateKey),
	)
}

func initVar() {
	cfgFile = viper.GetString("config")
	port = viper.GetString("port")

	dsn = setDSNTimeZone(viper.GetString("dsn"))

	isOverwrite = viper.GetBool("overwrite")
	staticFolderLocation = viper.GetString("static")
	userJSONPath = viper.GetString("import")
	lifeTime = viper.GetInt64("time")
	privateKey = viper.GetString("key")
}

func setDSNTimeZone(dsn string) string {
	config, err := mysql.ParseDSN(dsn)
	if err != nil {
		log.Fatalln(err)
	}

	loc, err := time.LoadLocation("Asia/Hong_Kong")
	if err != nil {
		log.Fatalln(err)
	}

	config.Loc = loc
	config.ParseTime = true

	return config.FormatDSN()
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
