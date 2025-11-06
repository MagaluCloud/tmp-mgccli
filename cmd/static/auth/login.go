package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/magaluCloud/mgccli/cmd/common/auth"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/spf13/cobra"
)

// NewLoginCommand cria o comando de login para o CLI
func NewLoginCommand(ctx context.Context) *cobra.Command {
	var opts auth.LoginOptions

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Autenticar na Magalu Cloud",
		Long:  "Executa o fluxo de autenticação OAuth para fazer login na Magalu Cloud",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogin(ctx, opts)
		},
	}

	// Configurar flags
	cmd.Flags().BoolVar(&opts.Headless, "headless", false, "Login sem abrir navegador (device flow)")
	cmd.Flags().BoolVar(&opts.QRCode, "qrcode", false, "Exibir QR code para login")
	cmd.Flags().BoolVar(&opts.Show, "show", false, "Exibir token de acesso após o login")

	// Marcar flags como mutuamente exclusivas
	cmd.MarkFlagsMutuallyExclusive("headless", "qrcode")

	return cmd
}

// runLogin executa o processo de login
func runLogin(ctx context.Context, opts auth.LoginOptions) error {
	auth := ctx.Value(cmdutils.CTX_AUTH_KEY).(auth.Auth)

	// Executar login
	fmt.Println("Iniciando processo de autenticação...")
	token, err := auth.GetService().Login(ctx, opts)
	if err != nil {
		return fmt.Errorf("falha na autenticação: %w", err)
	}

	// Exibir mensagem de sucesso
	fmt.Fprintln(os.Stderr, "\n✓ Autenticação realizada com sucesso!")

	auth.SetAccessToken(token.AccessToken)
	auth.SetRefreshToken(token.RefreshToken)

	return nil
}
