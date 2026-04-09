package buckets

import (
	"context"
	"fmt"

	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/common/config"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/buckets/common"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type publicURLOptions struct {
	Dst string
}

type publicURLReturn struct {
	URL string `json:"url"`
}

// PublicURLCommand cria o comando de exibir a URL pública
func PublicURLCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()
	var opts publicURLOptions

	cmd := &cobra.Command{
		Use:   "public-url [dst]",
		Short: manager.T("cli.auth.object_storage.buckets.public_url.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runPublicURL(ctx, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Dst, "dst", "", manager.T("cli.auth.object_storage.buckets.public_url.dst"))

	return cmd
}

// runPublicURL executa o processo de exibir a URL pública
func runPublicURL(ctx context.Context, args []string, opts publicURLOptions, rawMode bool) error {
	bucketName := opts.Dst

	if len(args) > 0 {
		bucketName = args[0]
	}

	if bucketName == "" {
		return cmdutils.NewCliError("missing required flag: --dst=string")
	}

	configCtx := ctx.Value(cmdutils.CXT_CONFIG_KEY).(config.Config)

	region, err := configCtx.Get("region")
	if err != nil {
		return cmdutils.NewCliError(fmt.Sprintf("erro ao pegar a configuração da região: %s", err.Error()))
	}

	bucketURL, err := common.BuildHost(bucketName, region.Value.(string))
	if err != nil {
		return err
	}

	beautiful.NewOutput(rawMode).PrintData(publicURLReturn{
		URL: bucketURL,
	})

	return nil
}
