package auth

import (
	"context"

	"github.com/magaluCloud/mgccli/beautiful"
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
		return cmdutils.NewCliError("your access token is empty. Please log in again.")
	}

	beautiful.NewOutput(false).PrintData(map[string]any{"access_token": token})

	return nil
}
