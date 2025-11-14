package auth

import (
	"github.com/magaluCloud/mgccli/i18n"

	"github.com/spf13/cobra"
)

// AuthCmd cria e configura o comando de autenticação
func AuthCmd(parent *cobra.Command) {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:     "auth",
		Short:   manager.T("cli.auth.short"),
		Long:    manager.T("cli.auth.long"),
		Aliases: []string{"auth"},
		GroupID: "settings",
	}

	// Adicionar subcomandos
	cmd.AddCommand(NewLoginCommand(parent.Context()))

	parent.AddCommand(cmd)
}
