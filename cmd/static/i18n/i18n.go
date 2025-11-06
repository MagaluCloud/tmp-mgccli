package i18n

import (
	"github.com/magaluCloud/mgccli/i18n"

	"github.com/spf13/cobra"
)

func I18nCmd(parent *cobra.Command) {
	manager := i18n.GetInstance()
	cmd := &cobra.Command{
		Use:     "i18n",
		Short:   manager.T("cli.i18n.short"),
		GroupID: "other",
		Long:    manager.T("cli.i18n.long"),
	}

	cmd.AddCommand(listCmd())
	cmd.AddCommand(setCmd())

	parent.AddCommand(cmd)
}
