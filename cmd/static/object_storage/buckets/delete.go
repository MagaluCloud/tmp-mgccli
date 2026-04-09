package buckets

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	objectstorage "github.com/magaluCloud/mgccli/cmd/common/object_storage"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type deleteOptions struct {
	Bucket    string
	Recursive bool
}

// DeleteCommand cria o comando de remover o bucket
func DeleteCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()
	var opts deleteOptions

	cmd := &cobra.Command{
		Use:   "delete [bucket]",
		Short: manager.T("cli.auth.object_storage.buckets.delete.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runDelete(ctx, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Bucket, "bucket", "", manager.T("cli.auth.object_storage.buckets.dst"))
	cmd.Flags().BoolVar(&opts.Recursive, "recursive", false, manager.T("cli.auth.object_storage.buckets.delete.recursive"))

	cmd.MarkFlagRequired("recursive")

	return cmd
}

// runDelete executa o processo de remover o bucket
func runDelete(ctx context.Context, args []string, opts deleteOptions, rawMode bool) error {
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

	var input string
	huh.NewInput().
		Title(fmt.Sprintf("This command will delete bucket %s, and its result is NOT reversible. Please confirm by retyping: %s", bucketName, bucketName)).
		Value(&input).
		Run()
	if input != bucketName {
		fmt.Println("Não foi possível deletar. O texto digitado não corresponde ao nome do bucket informado!")

		return nil
	}

	err = bucketService.Delete(ctx, bucketName, opts.Recursive)
	if err != nil {
		return cmdutils.NewCliError(err.Error())
	}

	fmt.Fprintln(os.Stderr, "✓ Deletado com sucesso!")

	return nil
}
