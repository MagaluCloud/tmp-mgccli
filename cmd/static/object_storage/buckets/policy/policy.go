package policy

import (
	"context"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

// PolicyCommand cria e configura o comando de política
func PolicyCommand(ctx context.Context, bucketService objSdk.BucketService) *cobra.Command {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "policy",
		Short: manager.T("cli.auth.object_storage.buckets.policy.short"),
		Long:  manager.T("cli.auth.object_storage.buckets.policy.long"),
	}

	cmd.AddCommand(SetCommand(ctx, bucketService))
	cmd.AddCommand(GetCommand(ctx, bucketService))
	cmd.AddCommand(DeleteCommand(ctx, bucketService))

	return cmd
}
