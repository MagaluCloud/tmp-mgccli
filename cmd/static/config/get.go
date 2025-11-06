package config

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Get() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Obter configurações",
		Long:  `Obter configurações`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Obter configurações")
		},
	}
	return cmd
}
