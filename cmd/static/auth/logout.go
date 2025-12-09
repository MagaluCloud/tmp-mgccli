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

// NewLogoutCommand cria o comando de logout para o CLI
func NewLogoutCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "logout",
		Short: manager.T("cli.auth.logout.short"),
		Long:  manager.T("cli.auth.logout.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogout(ctx)
		},
	}

	return cmd
}

// runLogout executa o processo de logout
func runLogout(ctx context.Context) error {
	auth := ctx.Value(cmdutils.CTX_AUTH_KEY).(auth.Auth)

	// Executar logout
	err := auth.Logout()
	if err != nil {
		return fmt.Errorf("falha ao encerrar a sessão: %w", err)
	}

	// Exibir mensagem de sucesso
	fmt.Fprintln(os.Stderr, "✓ Sessão encerrada com sucesso!")

	return nil
}
