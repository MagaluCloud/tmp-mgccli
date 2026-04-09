package buckets

import (
	"context"

	"github.com/magaluCloud/mgccli/cmd/static/object_storage/buckets/acl"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/buckets/cors"
	objectlock "github.com/magaluCloud/mgccli/cmd/static/object_storage/buckets/object_lock"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/buckets/policy"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/buckets/versioning"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

// BucketsCommand cria e configura o comando de buckets
func BucketsCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "buckets",
		Short: manager.T("cli.auth.object_storage.buckets.short"),
		Long:  manager.T("cli.auth.object_storage.buckets.long"),
	}

	cmd.AddCommand(ListCommand(ctx))
	cmd.AddCommand(CreateCommand(ctx))
	cmd.AddCommand(DeleteCommand(ctx))
	cmd.AddCommand(PublicURLCommand(ctx))

	cmd.AddCommand(cors.CORSCommand(ctx))
	cmd.AddCommand(policy.PolicyCommand(ctx))
	cmd.AddCommand(versioning.VersioningCommand(ctx))
	cmd.AddCommand(objectlock.ObjectLockCommand(ctx))
	cmd.AddCommand(acl.ACLCommand(ctx))

	return cmd
}
