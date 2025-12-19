package apikey

import (
	"context"
	"fmt"

	"github.com/magaluCloud/mgccli/beautiful"
	apiKeyPkg "github.com/magaluCloud/mgccli/cmd/common/api_key"
	"github.com/magaluCloud/mgccli/cmd/common/auth"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

// ListOptions configura opções para o processo de listar todas API Keys
type ListOptions struct {
	InvalidKeys bool // Incluir as API Keys inválidas
}

// ApiKeysResult representa o resultado simplificado da API Key
type apiKeysResult struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description,omitempty"`
	StartValidity string  `json:"start_validity"`
	EndValidity   *string `json:"end_validity,omitempty"`
	RevokedAt     *string `json:"revoked_at,omitempty"`
	TenantName    *string `json:"tenant_name,omitempty"`
}

// ListCommand cria o comando de listar todas API Keys
func ListCommand(ctx context.Context) *cobra.Command {
	var opts ListOptions

	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "list",
		Short: manager.T("cli.auth.api_key.list.short"),
		Long:  manager.T("cli.auth.api_key.list.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runList(ctx, opts, raw)
		},
	}

	cmd.Flags().BoolVar(&opts.InvalidKeys, "invalid-keys", false, manager.T("cli.auth.api_key.list.invalid_keys"))

	return cmd
}

// runList executa o processo de exibir todas API Keys
func runList(ctx context.Context, opts ListOptions, rawMode bool) error {
	authCtx := ctx.Value(cmdutils.CTX_AUTH_KEY).(auth.Auth)

	apiKey := apiKeyPkg.NewApiKey(authCtx)

	apiKeys, err := apiKey.List(ctx, opts.InvalidKeys)
	if err != nil {
		return fmt.Errorf("erro ao listar as API Keys: %w", err)
	}

	var result []*apiKeysResult
	for _, key := range apiKeys {
		result = append(result, &apiKeysResult{
			ID:            key.UUID,
			Name:          key.Name,
			Description:   key.Description,
			StartValidity: key.StartValidity,
			EndValidity:   key.EndValidity,
			RevokedAt:     key.RevokedAt,
			TenantName:    key.TenantName,
		})
	}

	beautiful.NewOutput(rawMode).PrintData(result)

	return nil
}
