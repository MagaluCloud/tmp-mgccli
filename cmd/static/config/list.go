package config

import (
	"fmt"
	"reflect"

	"github.com/fatih/color"
	"github.com/magaluCloud/mgccli/cmd/common/config"
	"github.com/spf13/cobra"
)

func List(config config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Listar configurações",
		Long:  `Listar configurações`,
		Run: func(cmd *cobra.Command, args []string) {

			configMap, err := config.List()
			if err != nil {
				fmt.Println("Erro ao listar configurações:", err)
				return
			}
			for key, value := range configMap {
				fmt.Printf("Name: %s\n   Value: %v\n   Type: %s\n   Description: %s\n   Validator: %s\n   Default: %v\n   Scope: %s\n\n",
					color.BlueString(key),
					color.YellowString(fmt.Sprintf("%v", valueOrDefault(value.Value, value.Default))),
					value.Type,
					value.Description,
					validatorOrEmpty(value.Validator),
					value.Default,
					value.Scope,
				)
			}
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

func validatorOrEmpty(validator *string) string {
	if validator == nil {
		return ""
	}
	return *validator
}
