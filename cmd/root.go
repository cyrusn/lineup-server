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
		"./config.yaml",
		"config file",
	)
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	rootCmd.PersistentFlags().StringVarP(
		&privateKey,
		"key",
		"k",
		"skill-vein-planet-neigh-envoi",
		"change the private key for authentication on jwt",
	)
	viper.BindPFlag("key", rootCmd.PersistentFlags().Lookup("key"))

	rootCmd.PersistentFlags().StringVarP(
		&dbPath,
		"location",
		"l",
		"./test/test.db",
		"location of sqlite3 database file",
	)
	viper.BindPFlag("location", rootCmd.PersistentFlags().Lookup("location"))

}
