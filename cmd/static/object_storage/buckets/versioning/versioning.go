package versioning

import (
	"context"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

// VersioningCommand cria e configura o comando de versionamento
func VersioningCommand(ctx context.Context, bucketService objSdk.BucketService) *cobra.Command {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "versioning",
		Short: manager.T("cli.auth.object_storage.buckets.versioning.short"),
	}

	cmd.AddCommand(EnableCommand(ctx, bucketService))
	cmd.AddCommand(SuspendCommand(ctx, bucketService))
	cmd.AddCommand(GetCommand(ctx, bucketService))

	return cmd
}
