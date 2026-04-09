package versioning

import (
	"context"

	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

// VersioningCommand cria e configura o comando de versionamento
func VersioningCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "versioning",
		Short: manager.T("cli.auth.object_storage.buckets.versioning.short"),
	}

	cmd.AddCommand(EnableCommand(ctx))
	cmd.AddCommand(SuspendCommand(ctx))
	cmd.AddCommand(GetCommand(ctx))

	return cmd
}
