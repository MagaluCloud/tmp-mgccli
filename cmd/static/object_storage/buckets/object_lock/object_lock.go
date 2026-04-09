package objectlock

import (
	"context"

	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

// ObjectLockCommand cria e configura o comando de bloqueio de objeto
func ObjectLockCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "object-lock",
		Short: manager.T("cli.auth.object_storage.buckets.object_lock.short"),
	}

	cmd.AddCommand(SetCommand(ctx))
	cmd.AddCommand(UnsetCommand(ctx))
	cmd.AddCommand(GetCommand(ctx))

	return cmd
}
