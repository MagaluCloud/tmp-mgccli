package objectlock

import (
	"context"
	"fmt"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type getOptions struct {
	Dst string
}

// GetCommand cria o comando de exibir a configuração de bloqueio de objetos
func GetCommand(ctx context.Context, bucketService objSdk.BucketService) *cobra.Command {
	manager := i18n.GetInstance()
	var opts getOptions

	cmd := &cobra.Command{
		Use:   "get [dst]",
		Short: manager.T("cli.auth.object_storage.buckets.object_lock.get.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runGet(ctx, bucketService, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Dst, "dst", "", manager.T("cli.auth.object_storage.buckets.dst"))

	return cmd
}

// runGet executa o processo de exibir a configuração de bloqueio de objetos
func runGet(ctx context.Context, bucketService objSdk.BucketService, args []string, opts getOptions, rawMode bool) error {
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

	config, err := bucketService.GetBucketLockConfig(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("erro ao pegar a configuração de bloqueio: %w", err)
	}

	beautiful.NewOutput(rawMode).PrintData(config)

	return nil
}
