package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

var timeoutFlag = "timeout"

func addTimeoutFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Int(
		timeoutFlag,
		60,
		"Set the timeout for the CLI in seconds",
	)
}

func getTimeoutFlag(cmd *cobra.Command) time.Duration {
	timeout, err := cmd.Root().PersistentFlags().GetInt(timeoutFlag)
	if err != nil {
		return 60 * time.Second
	}
	return time.Duration(timeout) * time.Second
}
