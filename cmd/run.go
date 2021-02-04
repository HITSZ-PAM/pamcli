package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a program within PAM context",
	Long:  `Resolve environment variables and pass them to the target program`,
	RunE:  runCmdExecute,
}

func runCmdExecute(cmd *cobra.Command, args []string) error {
	return nil
}
