package config

import (
	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/common/config"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/spf13/cobra"
)

type SetOptions struct {
	Key   string
	Value string
}

func Set(config config.Config) *cobra.Command {
	var opts SetOptions

	cmd := &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Definir configurações",
		Long:  `Definir configurações`,
		RunE: func(cmd *cobra.Command, args []string) error {
			key := opts.Key
			value := opts.Value

			if len(args) > 0 {
				key = args[0]
			}
			if len(args) > 1 {
				value = args[1]
			}

			if key == "" && value == "" {
				return cmdutils.NewCliError("missing required flags: --key=string, --value=any")
			}
			if key == "" {
				return cmdutils.NewCliError("missing required flag: --key=string")
			}
			if value == "" {
				return cmdutils.NewCliError("missing required flag: --value=any")
			}

			err := config.Set(key, value)
			if err != nil {
				return cmdutils.NewCliError(err.Error())
			}

			beautiful.NewOutput(false).PrintData(map[string]any{key: value})

			return nil
		},
	}

	cmd.Flags().StringVar(&opts.Key, "key", "", "Name of the desired config (required)")
	cmd.Flags().StringVar(&opts.Value, "value", "", "Value of the desired config (exactly one of: string, integer or boolean) (required)")

	return cmd
}
