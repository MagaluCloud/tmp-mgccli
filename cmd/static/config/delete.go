package config

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Delete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Deletar configurações",
		Long:  `Deletar configurações`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Deletando configurações")
		},
	}
	return cmd
}
