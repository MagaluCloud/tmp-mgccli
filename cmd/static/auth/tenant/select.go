package tenant

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

// SelectCommand cria o comando de selecionar o tenant
func SelectCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "select",
		Short: manager.T("cli.auth.tenant.select.short"),
		Long: fmt.Sprintf("%s\n\n⚠️  %s",
			manager.T("cli.auth.tenant.select.long"),
			manager.T("cli.auth.tenant.set.observation")),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runSelect(ctx, raw)
		},
	}
	return cmd
}

// runSelect executa o processo de selecionar o tenant atual
func runSelect(ctx context.Context, rawMode bool) error {
	auth := ctx.Value(cmdutils.CTX_AUTH_KEY).(auth.Auth)

	tenants, err := auth.ListTenants(ctx)
	if err != nil {
		return fmt.Errorf("erro ao listar os tenants: %w", err)
	}

	options := []huh.Option[string]{}
	for _, t := range tenants {
		options = append(options, huh.NewOption(t.Name, t.UUID))
	}

	selectedTenant := huh.NewSelect[string]()
	selectedTenant.Options(options...)
	err = selectedTenant.Run()
	if err != nil {
		return cmdutils.NewCliError(err.Error())
	}

	var uuid string
	selectedTenant.Value(&uuid)
	tokenInfo, err := auth.SetTenant(ctx, uuid)
	if err != nil {
		return cmdutils.NewCliError(err.Error())
	}

	accessKeyId := auth.GetAccessKeyID()
	secretAccessKey := auth.GetSecretAccessKey()

	if accessKeyId != "" && secretAccessKey != "" {
		fmt.Print("🔐 This operation unset the current api key. \n\n")

		err := auth.SetAccessKeyID("")
		if err != nil {
			return fmt.Errorf("erro ao remover o access key id: %w", err)
		}

		err = auth.SetSecretAccessKey("")
		if err != nil {
			return fmt.Errorf("erro ao remover o secret access key: %w", err)
		}
	}

	beautiful.NewOutput(rawMode).PrintData(tokenInfo)

	return nil
}
