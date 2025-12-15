package tenant

import (
	"context"
	"fmt"

	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/common/auth"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

// ListCommand cria o comando de listar todos os tenants
func ListCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "list",
		Short: manager.T("cli.auth.tenant.list.short"),
		Long:  manager.T("cli.auth.tenant.list.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runList(ctx, raw)
		},
	}

	return cmd
}

// runList executa o processo de exibir todos os tenants
func runList(ctx context.Context, rawMode bool) error {
	auth := ctx.Value(cmdutils.CTX_AUTH_KEY).(auth.Auth)

	tenants, err := auth.ListTenants(ctx)

	if err != nil {
		return fmt.Errorf("erro ao listar os tenants: %w", err)
	}

	beautiful.NewOutput(rawMode).PrintData(tenants)

	return nil
}
