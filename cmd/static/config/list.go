package config

import (
	"reflect"

	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/common/config"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/spf13/cobra"
)

func List(config config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Listar configurações",
		Long:  `Listar configurações`,
		RunE: func(cmd *cobra.Command, args []string) error {
			configMap, err := config.List()
			if err != nil {
				return cmdutils.NewCliError(err.Error())
			}

			for key, value := range configMap {
				configMap[key].Value = valueOrDefault(value.Value, value.Default)
			}

			beautiful.NewOutput(false).PrintData(configMap)

			return nil
		},
	}
	return cmd
}

func valueOrDefault(value any, defaultValue any) any {
	if value == nil || reflect.ValueOf(value).IsZero() {
		return defaultValue
	}
	return value
}
