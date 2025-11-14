package config

import (
	"github.com/magaluCloud/mgccli/cmd/common/config"
	"github.com/magaluCloud/mgccli/i18n"

	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/spf13/cobra"
)

func ConfigCmd(parent *cobra.Command) {
	manager := i18n.GetInstance()
	cmd := &cobra.Command{
		Use:     "config",
		Short:   manager.T("cli.config.short"),
		Long:    manager.T("cli.config.long"),
		Aliases: []string{"cfg"},
		GroupID: "settings",
	}

	config := parent.Context().Value(cmdutils.CXT_CONFIG_KEY).(config.Config)

	cmd.AddCommand(List(config))
	cmd.AddCommand(Delete(config))
	cmd.AddCommand(Get(config))
	cmd.AddCommand(Set(config))

	parent.AddCommand(cmd)
}
