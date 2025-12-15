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

// CurrentCommand cria o comando de exibir o tenant atual
func CurrentCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "current",
		Short: manager.T("cli.auth.tenant.current.short"),
		Long:  manager.T("cli.auth.tenant.current.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runCurrent(ctx, raw)
		},
	}

	return cmd
}

// runCurrent executa o processo de exibir o tenant atual
func runCurrent(ctx context.Context, rawMode bool) error {
	auth := ctx.Value(cmdutils.CTX_AUTH_KEY).(auth.Auth)

	tenant, err := auth.GetCurrentTenant(ctx)
	if err != nil {
		return fmt.Errorf("erro ao pegar o tenant atual: %w", err)
	}

	beautiful.NewOutput(rawMode).PrintData(tenant)

	return nil
}
