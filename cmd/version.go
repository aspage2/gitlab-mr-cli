package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var AppVersion = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version and exit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("glmr %s\n", green(AppVersion))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
