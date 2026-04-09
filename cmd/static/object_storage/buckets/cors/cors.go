package cors

import (
	"context"

	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

// CORSCommand cria e configura o comando de CORS
func CORSCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "cors",
		Short: manager.T("cli.auth.object_storage.buckets.cors.short"),
	}

	cmd.AddCommand(SetCommand(ctx))
	cmd.AddCommand(GetCommand(ctx))
	cmd.AddCommand(DeleteCommand(ctx))

	return cmd
}
