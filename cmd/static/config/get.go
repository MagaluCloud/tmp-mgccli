package config

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/magaluCloud/mgccli/cmd/common/config"
	"github.com/spf13/cobra"
)

func Get(config config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [config]",
		Short: "Obter configurações",
		Long:  `Obter configurações`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Erro: configuração não especificada")
				return
			}
			value, err := config.Get(args[0])
			if err != nil {
				fmt.Println("Erro ao obter configuração:", err)
				return
			}
			fmt.Printf("%s: %v\n", color.BlueString(args[0]), color.YellowString(fmt.Sprintf("%v", value.Value)))
		},
	}
	return cmd
}
