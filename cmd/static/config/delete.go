package config

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/magaluCloud/mgccli/cmd/common/config"
	"github.com/spf13/cobra"
)

func Delete(config config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [config]",
		Short: "Deletar configurações",
		Long:  `Deletar configurações`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Erro: configuração não especificada")
				return
			}
			err := config.Delete(args[0])
			if err != nil {
				fmt.Println("Erro ao deletar configuração:", err)
				return
			}
			fmt.Println(color.GreenString("Configuração deletada com sucesso"))
		},
	}
	return cmd
}
