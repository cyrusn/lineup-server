package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "lineup",
	Short: "Welcome to Line-Up System Backend Server",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&cfgFile,
		"config",
		"c",
		DEFAULT_CONFIG_PATH,
		"config file",
	)
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	rootCmd.PersistentFlags().StringVarP(
		&privateKey,
		"key",
		"k",
		DEFAULT_PRIVATE_KEY,
		"change the private key for authentication on jwt",
	)
	viper.BindPFlag("key", rootCmd.PersistentFlags().Lookup("key"))

	rootCmd.PersistentFlags().StringVarP(
		&dsn,
		"dsn",
		"d",
		DEFAULT_DSN,
		"Data source name of mysql. [ref https://github.com/go-sql-driver/mysql]",
	)
	viper.BindPFlag("location", rootCmd.PersistentFlags().Lookup("location"))

}
