package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pamcli",
	Short: "pamcli is a Command Line Interface of Privilege Account Manager",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		return
	},
}

// Execute is provided to main.main() acting as a entrance
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
