package versioning

import (
	"context"
	"fmt"
	"os"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type suspendOptions struct {
	Bucket string
}

// SuspendCommand cria o comando de suspender o versionamento do bucket
func SuspendCommand(ctx context.Context, bucketService objSdk.BucketService) *cobra.Command {
	manager := i18n.GetInstance()
	var opts suspendOptions

	cmd := &cobra.Command{
		Use:   "suspend [bucket]",
		Short: manager.T("cli.auth.object_storage.buckets.versioning.suspend.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runSuspend(ctx, bucketService, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Bucket, "bucket", "", manager.T("cli.auth.object_storage.buckets.dst"))

	return cmd
}

// runSuspend executa o processo de suspender o versionamento do bucket
func runSuspend(ctx context.Context, bucketService objSdk.BucketService, args []string, opts suspendOptions, rawMode bool) error {
	if bucketService == nil {
		return nil
	}

	bucketName := opts.Bucket

	if len(args) > 0 {
		bucketName = args[0]
	}

	if bucketName == "" {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o nome do bucket como argumento ou usar a flag --bucket")

		return nil
	}

	err := bucketService.SuspendVersioning(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("erro ao suspender o versionamento: %w", err)
	}

	fmt.Fprintln(os.Stderr, "✓ Suspendido o versionamento do bucket!")

	return nil
}
