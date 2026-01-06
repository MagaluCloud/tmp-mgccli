package buckets

import (
	"context"
	"fmt"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type bucketResponse struct {
	CreationDate string `xml:"CreationDate"`
	Name         string `xml:"Name"`
}

// ListCommand cria o comando de listar todos os buckets
func ListCommand(ctx context.Context, bucketService objSdk.BucketService) *cobra.Command {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "list",
		Short: manager.T("cli.auth.object_storage.buckets.list.short"),
		Long:  manager.T("cli.auth.object_storage.buckets.list.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runList(ctx, bucketService, raw)
		},
	}

	return cmd
}

// runList executa o processo de exibir todos os buckets
func runList(ctx context.Context, bucketService objSdk.BucketService, rawMode bool) error {
	if bucketService == nil {
		return nil
	}

	buckets, err := bucketService.List(ctx)
	if err != nil {
		return fmt.Errorf("erro ao listar os buckets: %w", err)
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
