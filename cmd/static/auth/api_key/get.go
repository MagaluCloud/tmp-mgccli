package apikey

import (
	"context"
	"fmt"

	"github.com/magaluCloud/mgccli/beautiful"
	apiKeyPkg "github.com/magaluCloud/mgccli/cmd/common/api_key"
	authPkg "github.com/magaluCloud/mgccli/cmd/common/auth"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type GetOptions struct {
	ID string
}

type GetApiKeyResult struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	ApiKey        string          `json:"api_key"`
	Description   string          `json:"description,omitempty"`
	KeyPairID     string          `json:"key_pair_id"`
	KeyPairSecret string          `json:"key_pair_secret"`
	StartValidity string          `json:"start_validity"`
	EndValidity   *string         `json:"end_validity,omitempty"`
	RevokedAt     *string         `json:"revoked_at,omitempty"`
	TenantName    *string         `json:"tenant_name,omitempty"`
	Scopes        []authPkg.Scope `json:"scopes"`
}

// GetCommand cria o comando de exibir uma API Key
func GetCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()

	var opts GetOptions

	cmd := &cobra.Command{
		Use:   "get [id]",
		Short: manager.T("cli.auth.api_key.get.short"),
		Long:  manager.T("cli.auth.api_key.get.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runGet(ctx, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.ID, "id", "", manager.T("cli.auth.api_key.get.id"))

	return cmd
}

// runGet executa o processo de exibir uma API Key
func runGet(ctx context.Context, args []string, opts GetOptions, rawMode bool) error {
	ID := opts.ID

	if len(args) > 0 {
		ID = args[0]
	}

	if ID == "" {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o ID como argumento ou usar a flag --id")

		return nil
	}

	auth := ctx.Value(cmdutils.CTX_AUTH_KEY).(authPkg.Auth)
	apiKey := apiKeyPkg.NewApiKey(auth)

	apiKeys, err := apiKey.List(ctx, true)
	if err != nil {
		return fmt.Errorf("erro ao listar as API Keys: %w", err)
	}

	for _, key := range apiKeys {
		if key.UUID == ID {
			scopes := []authPkg.Scope{}
			for _, scope := range key.Scopes {
				scopes = append(scopes, authPkg.Scope{
					Name:  scope.Name,
					Title: scope.Title,
					UUID:  scope.UUID,
				})
			}

			beautiful.NewOutput(rawMode).PrintData(&GetApiKeyResult{
				ID:            key.UUID,
				Name:          key.Name,
				Description:   key.Description,
				StartValidity: key.StartValidity,
				EndValidity:   key.EndValidity,
				RevokedAt:     key.RevokedAt,
				TenantName:    key.TenantName,
				ApiKey:        key.ApiKey,
				KeyPairID:     key.KeyPairID,
				KeyPairSecret: key.KeyPairSecret,
				Scopes:        scopes,
			})

			return nil
		}
	}

	return fmt.Errorf("não foi possível encontrar a API Key com o ID %q", ID)
}
