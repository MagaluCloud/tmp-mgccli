package versioning

import (
	"context"
	"fmt"
	"os"

	objectstorage "github.com/magaluCloud/mgccli/cmd/common/object_storage"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type suspendOptions struct {
	Bucket string
}

// SuspendCommand cria o comando de suspender o versionamento do bucket
func SuspendCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()
	var opts suspendOptions

	cmd := &cobra.Command{
		Use:   "suspend [bucket]",
		Short: manager.T("cli.auth.object_storage.buckets.versioning.suspend.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runSuspend(ctx, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Bucket, "bucket", "", manager.T("cli.auth.object_storage.buckets.dst"))

	return cmd
}

// runSuspend executa o processo de suspender o versionamento do bucket
func runSuspend(ctx context.Context, args []string, opts suspendOptions, rawMode bool) error {
	objectStorageService, err := objectstorage.NewObjectStorage(ctx)
	if err != nil {
		return cmdutils.NewCliError(err.Error())
	}
	bucketService := objectStorageService.GetBucketService()

	bucketName := opts.Bucket

	if len(args) > 0 {
		bucketName = args[0]
	}

	if bucketName == "" {
		return cmdutils.NewCliError("missing required flag: --bucket=string")
	}

	err = bucketService.SuspendVersioning(ctx, bucketName)
	if err != nil {
		return cmdutils.NewCliError(err.Error())
	}

	fmt.Fprintln(os.Stderr, "✓ Suspendido o versionamento do bucket!")

	return nil
}
