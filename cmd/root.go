package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var versionTag string

var rootCmd = &cobra.Command{
	Use:   "congo [param]",
	Short: "Connect to AWS and ECS easily.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("congo", versionTag)
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
