package objectlock

import (
	"context"
	"fmt"
	"os"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type unsetOptions struct {
	Dst string
}

// UnsetCommand cria o comando de desbloquear o objeto
func UnsetCommand(ctx context.Context, bucketService objSdk.BucketService) *cobra.Command {
	manager := i18n.GetInstance()
	var opts unsetOptions

	cmd := &cobra.Command{
		Use:   "unset [dst]",
		Short: manager.T("cli.auth.object_storage.buckets.object_lock.unset.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runUnset(ctx, bucketService, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Dst, "dst", "", manager.T("cli.auth.object_storage.buckets.dst"))

	return cmd
}

// runUnset executa o processo de desbloquear o objeto
func runUnset(ctx context.Context, bucketService objSdk.BucketService, args []string, opts unsetOptions, rawMode bool) error {
	if bucketService == nil {
		return nil
	}

	bucketName := opts.Dst

	if len(args) > 0 {
		bucketName = args[0]
	}

	if bucketName == "" {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o nome do bucket como argumento ou usar a flag --dst")
		return nil
	}

	err := bucketService.UnlockBucket(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("erro ao desbloquear: %w", err)
	}

	fmt.Fprintln(os.Stderr, "✓ Desbloqueado com sucesso!")

	return nil
}
