package policy

import (
	"context"

	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

// PolicyCommand cria e configura o comando de política
func PolicyCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "policy",
		Short: manager.T("cli.auth.object_storage.buckets.policy.short"),
		Long:  manager.T("cli.auth.object_storage.buckets.policy.long"),
	}

	cmd.AddCommand(SetCommand(ctx))
	cmd.AddCommand(GetCommand(ctx))
	cmd.AddCommand(DeleteCommand(ctx))

	return cmd
}
