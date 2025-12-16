package apikey

import (
	"context"

	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

// ApiKeyCommand cria e configura o comando de API Key
func ApiKeyCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "api-key",
		Short: manager.T("cli.auth.api_key.short"),
		Long:  manager.T("cli.auth.api_key.long"),
	}

	cmd.AddCommand(ListCommand(ctx))
	cmd.AddCommand(CreateCommand(ctx))
	cmd.AddCommand(GetCommand(ctx))
	cmd.AddCommand(RevokeCommand(ctx))

	return cmd
}
