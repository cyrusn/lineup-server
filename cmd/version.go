package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "v2.0.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Line-Up System Backend Server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}
