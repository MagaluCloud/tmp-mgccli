package apikey

import (
	"context"
	"fmt"

	"github.com/magaluCloud/mgccli/cmd/common/auth"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

// CreateOptions configura opções para o processo de criar uma API Key
type CreateOptions struct {
	Name        string
	Description string
	Scopes      []string
	Expiration  string
}

// CreateCommand cria o comando de criar uma API Key
func CreateCommand(ctx context.Context) *cobra.Command {
	var opts CreateOptions

	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:     "create",
		Short:   manager.T("cli.auth.api_key.create.short"),
		Long:    manager.T("cli.auth.api_key.create.long"),
		Example: `mgc auth api-key create --description="created from MGC CLI" --expiration="2024-11-07" --name="My MGC Key" --scopes=["dbaas.read"]`,
		RunE: func(cmd *cobra.Command, args []string) error {

			return runCreate(ctx, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", manager.T("cli.auth.api_key.create.name"))
	cmd.Flags().StringVar(&opts.Description, "description", "", manager.T("cli.auth.api_key.create.description"))
	cmd.Flags().StringSliceVar(&opts.Scopes, "scopes", nil, manager.T("cli.auth.api_key.create.scopes"))
	cmd.Flags().StringVar(&opts.Expiration, "expiration", "", manager.T("cli.auth.api_key.create.expiration"))

	cmd.MarkFlagRequired("name")

	return cmd
}

// runCreate executa o processo de criar uma API Key
func runCreate(ctx context.Context, opts CreateOptions) error {
	auth := ctx.Value(cmdutils.CTX_AUTH_KEY).(auth.Auth)

	apiKey, err := auth.CreateApiKey(ctx, opts.Name, opts.Description, opts.Expiration, opts.Scopes)
	if err != nil {
		return fmt.Errorf("erro ao criar a API Key: %w", err)
	}

	fmt.Printf("API Key criada com sucesso!\nID: %s\n", apiKey.UUID)

	return nil
}
