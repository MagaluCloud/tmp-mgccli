package buckets

import (
	"context"

	"github.com/magaluCloud/mgccli/beautiful"
	objectstorage "github.com/magaluCloud/mgccli/cmd/common/object_storage"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type bucketResponse struct {
	CreationDate string `xml:"CreationDate"`
	Name         string `xml:"Name"`
}

// ListCommand cria o comando de listar todos os buckets
func ListCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "list",
		Short: manager.T("cli.auth.object_storage.buckets.list.short"),
		Long:  manager.T("cli.auth.object_storage.buckets.list.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runList(ctx, raw)
		},
	}

	return cmd
}

// runList executa o processo de exibir todos os buckets
func runList(ctx context.Context, rawMode bool) error {
	objectStorageService, err := objectstorage.NewObjectStorage(ctx)
	if err != nil {
		return cmdutils.NewCliError(err.Error())
	}
	bucketService := objectStorageService.GetBucketService()

	buckets, err := bucketService.List(ctx)
	if err != nil {
		return cmdutils.NewCliError(err.Error())
	}

	result := []bucketResponse{}

	for _, bucket := range buckets {
		result = append(result, bucketResponse{
			Name:         bucket.Name,
			CreationDate: bucket.CreationDate.UTC().Format("2006-01-02T15:04:05.000Z"),
		})
	}

	beautiful.NewOutput(rawMode).PrintData(result)

	return nil
}
