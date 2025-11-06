package cmd

import "github.com/spf13/cobra"

const LangFlag = "lang"

func addLangFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String(
		LangFlag,
		"en-US",
		"Set the language for the CLI",
	)
}
