package objectlock

import (
	"context"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

// ObjectLockCommand cria e configura o comando de bloqueio de objeto
func ObjectLockCommand(ctx context.Context, bucketService objSdk.BucketService) *cobra.Command {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "object-lock",
		Short: manager.T("cli.auth.object_storage.buckets.object_lock.short"),
	}

	cmd.AddCommand(SetCommand(ctx, bucketService))
	cmd.AddCommand(UnsetCommand(ctx, bucketService))
	cmd.AddCommand(GetCommand(ctx, bucketService))

	return cmd
}
