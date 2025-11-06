package config

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Set() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Definir configurações",
		Long:  `Definir configurações`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Definindo configurações")
		},
	}
	return cmd
}
