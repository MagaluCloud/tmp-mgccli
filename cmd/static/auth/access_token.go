package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/magaluCloud/mgccli/cmd/common/auth"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

// NewAccessTokenCommand cria o comando para obter o token de acesso
func NewAccessTokenCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:     "access-token",
		Aliases: []string{"access_token"},
		Short:   manager.T("cli.auth.access_token.short"),
		Long:    manager.T("cli.auth.access_token.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAccessToken(ctx)
		},
	}

	return cmd
}

// runAccessToken executa o processo de obtenção do token de acesso
func runAccessToken(ctx context.Context) error {
	auth := ctx.Value(cmdutils.CTX_AUTH_KEY).(auth.Auth)

	token := auth.GetAccessToken(ctx)
	if token == "" {
		fmt.Fprintln(os.Stderr, "✕ Seu token de acesso está vazio. Por favor, faça login novamente.")

		return nil
	}

	// Exibir token de acesso
	fmt.Fprintf(os.Stderr, "access_token:\n%s\n", token)

	return nil
}
