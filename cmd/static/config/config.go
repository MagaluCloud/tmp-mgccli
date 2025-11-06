package config

import (
	"github.com/magaluCloud/mgccli/i18n"

	sdk "github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/spf13/cobra"
)

func ConfigCmd(parent *cobra.Command, sdkCoreConfig sdk.CoreClient) {
	manager := i18n.GetInstance()
	cmd := &cobra.Command{
		Use:     "config",
		Short:   manager.T("cli.config.short"),
		Long:    manager.T("cli.config.long"),
		Aliases: []string{"cfg"},
		GroupID: "settings",
	}

	cmd.AddCommand(List())
	cmd.AddCommand(Delete())
	cmd.AddCommand(Get())
	cmd.AddCommand(Set())

	parent.AddCommand(cmd)
}
