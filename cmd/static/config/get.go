package config

import (
	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/common/config"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/spf13/cobra"
)

type GetOptions struct {
	Key string
}

func Get(config config.Config) *cobra.Command {
	var opts GetOptions

	cmd := &cobra.Command{
		Use:   "get [key]",
		Short: "Obter configurações",
		Long:  `Obter configurações`,
		RunE: func(cmd *cobra.Command, args []string) error {
			key := opts.Key
			if len(args) > 0 {
				key = args[0]
			}

			if key == "" {
				return cmdutils.NewCliError("missing required flag: --key=string")
			}

			value, err := config.Get(key)
			if err != nil {
				return cmdutils.NewCliError(err.Error())
			}

			beautiful.NewOutput(false).PrintData(map[string]any{key: value.Value})

			return nil
		},
	}

	cmd.Flags().StringVar(&opts.Key, "key", "", "Name of the desired config (required)")

	return cmd
}
