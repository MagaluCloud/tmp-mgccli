package config

import (
	"fmt"

	"github.com/spf13/cobra"
)

func List() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Listar configurações",
		Long:  `Listar configurações`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Listando configurações")
		},
	}
	return cmd
}
