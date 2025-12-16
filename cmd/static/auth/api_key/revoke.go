package apikey

import (
	"context"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/common/auth"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type RevokeOptions struct {
	ID string
}

// RevokeCommand cria o comando de revogar uma API Key
func RevokeCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()

	var opts RevokeOptions

	cmd := &cobra.Command{
		Use:   "revoke [id]",
		Short: manager.T("cli.auth.api_key.revoke.short"),
		Long:  manager.T("cli.auth.api_key.revoke.long"),
		RunE: func(cmd *cobra.Command, args []string) error {

			return runRevoke(ctx, args, opts)
		},
	}

	cmd.Flags().StringVar(&opts.ID, "id", "", manager.T("cli.auth.api_key.revoke.id"))

	return cmd
}

// runRevoke executa o processo de revogar uma API Key
func runRevoke(ctx context.Context, args []string, opts RevokeOptions) error {
	ID := opts.ID

	if len(args) > 0 {
		ID = args[0]
	}

	if ID == "" {
		beautiful.NewOutput(false).PrintError("é necessário fornecer o ID como argumento ou usar a flag --id")

		return nil
	}

	var confirm bool
	huh.NewConfirm().Title(fmt.Sprintf("This operation will permanently revoke the API Key %s. Do you wish to continue?", ID)).
		Affirmative("Yes").
		Negative("No").Value(&confirm).Run()
	if !confirm {
		return nil
	}

	auth := ctx.Value(cmdutils.CTX_AUTH_KEY).(auth.Auth)

	err := auth.RevokeApiKey(ctx, ID)
	if err != nil {
		return fmt.Errorf("erro ao revogar a API Key: %w", err)
	}

	fmt.Printf("API Key %s revogada com sucesso!\n", ID)

	return nil
}
