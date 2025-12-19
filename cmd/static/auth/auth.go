package auth

import (
	apikey "github.com/magaluCloud/mgccli/cmd/static/auth/api_key"
	"github.com/magaluCloud/mgccli/cmd/static/auth/clients"
	"github.com/magaluCloud/mgccli/cmd/static/auth/tenant"
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
		GroupID: "settings",
	}

	// Adicionar subcomandos
	cmd.AddCommand(NewLoginCommand(parent.Context()))
	cmd.AddCommand(NewLogoutCommand(parent.Context()))
	cmd.AddCommand(NewAccessTokenCommand(parent.Context()))

	cmd.AddCommand(tenant.TenantCommand(parent.Context()))
	cmd.AddCommand(apikey.ApiKeyCommand(parent.Context()))
	cmd.AddCommand(clients.ClientsCommand(parent.Context()))

	parent.AddCommand(cmd)
}
