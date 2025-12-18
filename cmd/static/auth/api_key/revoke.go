package apikey

import (
	"context"
	"fmt"
	"net/http"

	"github.com/charmbracelet/huh"
	"github.com/magaluCloud/mgccli/beautiful"
	authPkg "github.com/magaluCloud/mgccli/cmd/common/auth"
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

	err := revokeApiKey(ctx, ID)
	if err != nil {
		return fmt.Errorf("erro ao revogar a API Key: %w", err)
	}

	fmt.Printf("API Key %s revogada com sucesso!\n", ID)

	return nil
}

func revokeApiKey(ctx context.Context, ID string) error {
	authCtx := ctx.Value(cmdutils.CTX_AUTH_KEY).(authPkg.Auth)

	config := authCtx.GetConfig()

	client, err := authPkg.NewOAuthClient(config)
	if err != nil {
		return fmt.Errorf("failed to create OAuth client: %w", err)
	}

	httpClient := client.AuthenticatedHttpClientFromContext(ctx)
	if httpClient == nil {
		return fmt.Errorf("programming error: unable to get HTTP Client from context")
	}

	url := fmt.Sprintf("%s/%s/revoke", config.ApiKeysURLV1, ID)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return cmdutils.NewHttpErrorFromResponse(resp, req)
	}

	return nil
}
