package cmd

import "github.com/spf13/cobra"

const noConfirmationFlag = "no-confirm"

func addNoConfirmationFlag(cmd *cobra.Command) {
	cmd.Root().PersistentFlags().Bool(
		noConfirmationFlag,
		false,
		"Bypasses confirmation step for commands that ask a confirmation from the user",
	)
}

func getNoConfirmationFlag(cmd *cobra.Command) bool {
	allow, err := cmd.Root().PersistentFlags().GetBool(noConfirmationFlag)
	if err != nil {
		return false
	}
	return allow
}
