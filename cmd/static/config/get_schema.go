package config

import (
	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/common/config"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/spf13/cobra"
)

type GetSchemaOptions struct {
	Key string
}

func GetSchema(config config.Config) *cobra.Command {
	var opts GetSchemaOptions

	cmd := &cobra.Command{
		Use:   "get-schema [key]",
		Short: "Obter schema da configuração especificada",
		Long:  `Obter schema da configuração especificada`,
		RunE: func(cmd *cobra.Command, args []string) error {
			key := opts.Key
			if len(args) > 0 {
				key = args[0]
			}

			if key == "" {
				return cmdutils.NewCliError("missing required flag: --key=string")
			}

			configMap, err := config.GetSchema(key)
			if err != nil {
				return cmdutils.NewCliError(err.Error())
			}

			beautiful.NewOutput(false).PrintData(configMap)

			return nil
		},
	}

	cmd.Flags().StringVar(&opts.Key, "key", "", "Name of the desired config (required)")

	return cmd
}
