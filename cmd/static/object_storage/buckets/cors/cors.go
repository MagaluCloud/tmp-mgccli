package cors

import (
	"context"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

// CORSCommand cria e configura o comando de CORS
func CORSCommand(ctx context.Context, bucketService objSdk.BucketService) *cobra.Command {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "cors",
		Short: manager.T("cli.auth.object_storage.buckets.cors.short"),
	}

	cmd.AddCommand(SetCommand(ctx, bucketService))
	cmd.AddCommand(GetCommand(ctx, bucketService))
	cmd.AddCommand(DeleteCommand(ctx, bucketService))

	return cmd
}
