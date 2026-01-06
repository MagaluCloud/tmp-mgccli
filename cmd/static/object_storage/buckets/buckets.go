package buckets

import (
	"context"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/buckets/acl"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/buckets/cors"
	objectlock "github.com/magaluCloud/mgccli/cmd/static/object_storage/buckets/object_lock"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/buckets/policy"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/buckets/versioning"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

// BucketsCommand cria e configura o comando de buckets
func BucketsCommand(ctx context.Context, bucketService objSdk.BucketService) *cobra.Command {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "buckets",
		Short: manager.T("cli.auth.object_storage.buckets.short"),
		Long:  manager.T("cli.auth.object_storage.buckets.long"),
	}

	cmd.AddCommand(ListCommand(ctx, bucketService))
	cmd.AddCommand(CreateCommand(ctx, bucketService))
	cmd.AddCommand(DeleteCommand(ctx, bucketService))
	cmd.AddCommand(PublicURLCommand(ctx))

	cmd.AddCommand(cors.CORSCommand(ctx, bucketService))
	cmd.AddCommand(policy.PolicyCommand(ctx, bucketService))
	cmd.AddCommand(versioning.VersioningCommand(ctx, bucketService))
	cmd.AddCommand(objectlock.ObjectLockCommand(ctx, bucketService))
	cmd.AddCommand(acl.ACLCommand(ctx))

	return cmd
}
