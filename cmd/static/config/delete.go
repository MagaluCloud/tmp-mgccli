package config

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/magaluCloud/mgccli/cmd/common/config"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/spf13/cobra"
)

type DeleteOptions struct {
	Key string
}

func Delete(config config.Config) *cobra.Command {
	var opts GetOptions

	cmd := &cobra.Command{
		Use:   "delete [config]",
		Short: "Deletar configurações",
		Long:  `Deletar configurações`,
		RunE: func(cmd *cobra.Command, args []string) error {
			key := opts.Key
			if len(args) > 0 {
				key = args[0]
			}

			if key == "" {
				return cmdutils.NewCliError("missing required flag: --key=string")
			}

			err := config.Delete(key)
			if err != nil {
				return cmdutils.NewCliError(err.Error())
			}

			fmt.Println(color.GreenString("Configuration deleted successfully"))
			return nil
		},
	}

	cmd.Flags().StringVar(&opts.Key, "key", "", "Name of the desired config (required)")

	return cmd
}
