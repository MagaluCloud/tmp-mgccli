package tenant

import (
	"context"

	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

// TenantCommand cria e configura o comando de tenant
func TenantCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "tenant",
		Short: manager.T("cli.auth.tenant.short"),
		Long:  manager.T("cli.auth.tenant.long"),
	}

	cmd.AddCommand(CurrentCommand(ctx))
	cmd.AddCommand(ListCommand(ctx))
	cmd.AddCommand(SetCommand(ctx))
	cmd.AddCommand(SelectCommand(ctx))

	return cmd
}
