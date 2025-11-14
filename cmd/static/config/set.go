package config

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/magaluCloud/mgccli/cmd/common/config"
	"github.com/spf13/cobra"
)

func Set(config config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [config] [value]",
		Short: "Definir configurações",
		Long:  `Definir configurações`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				fmt.Println("Erro: configuração e valor não especificados")
				return
			}

			err := config.Set(args[0], args[1])
			if err != nil {
				fmt.Println("Erro ao definir configuração:", err)
				return
			}
			fmt.Printf("%s: %v\n", color.BlueString(args[0]), color.YellowString(args[1]))
		},
	}
	return cmd
}
