package cors

import (
	"context"
	"fmt"
	"os"

	objectstorage "github.com/magaluCloud/mgccli/cmd/common/object_storage"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type deleteOptions struct {
	Dst string
}

// DeleteCommand cria o comando de remover a configuração de CORS do bucket
func DeleteCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()
	var opts deleteOptions

	cmd := &cobra.Command{
		Use:   "delete [dst]",
		Short: manager.T("cli.auth.object_storage.buckets.cors.delete.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runDelete(ctx, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Dst, "dst", "", manager.T("cli.auth.object_storage.buckets.dst"))

	return cmd
}

// runDelete executa o processo de remover a configuração de CORS do bucket
func runDelete(ctx context.Context, args []string, opts deleteOptions, rawMode bool) error {
	objectStorageService, err := objectstorage.NewObjectStorage(ctx)
	if err != nil {
		return cmdutils.NewCliError(err.Error())
	}
	bucketService := objectStorageService.GetBucketService()

	bucketName := opts.Dst

	if len(args) > 0 {
		bucketName = args[0]
	}

	if bucketName == "" {
		return cmdutils.NewCliError("missing required flag: --dst=string")
	}

	err = bucketService.DeleteCORS(ctx, bucketName)
	if err != nil {
		return cmdutils.NewCliError(err.Error())
	}

	fmt.Fprintln(os.Stderr, "✓ Deletado com sucesso!")

	return nil
}
