package clients

import (
	"context"

	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

// ClientsCommand cria e configura o comando de clients
func ClientsCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "clients",
		Short: manager.T("cli.auth.clients.short"),
		Long:  manager.T("cli.auth.clients.long"),
	}

	cmd.AddCommand(ListCommand(ctx))
	cmd.AddCommand(CreateCommand(ctx))
	cmd.AddCommand(UpdateCommand(ctx))

	return cmd
}
