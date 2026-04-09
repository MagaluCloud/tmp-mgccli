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

type SetOptions struct {
	UUID string
}

// SetCommand cria o comando de definir o tenant atual
func SetCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()

	var opts SetOptions

	cmd := &cobra.Command{
		Use:   "set [uuid]",
		Short: manager.T("cli.auth.tenant.set.short"),
		Long: fmt.Sprintf("%s\n\n⚠️  %s",
			manager.T("cli.auth.tenant.set.long"),
			manager.T("cli.auth.tenant.set.observation")),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runSet(ctx, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.UUID, "uuid", "", "The UUID of the desired Tenant. To list all possible IDs, run auth tenant list (required)")

	return cmd
}

// runSet executa o processo de definir o tenant atual
func runSet(ctx context.Context, args []string, opts SetOptions, rawMode bool) error {
	uuid := opts.UUID

	if len(args) > 0 {
		uuid = args[0]
	}

	if uuid == "" {
		return cmdutils.NewCliError("missing required flag: --uuid")
	}

	auth := ctx.Value(cmdutils.CTX_AUTH_KEY).(auth.Auth)

	tokenInfo, err := auth.SetTenant(ctx, uuid)
	if err != nil {
		return cmdutils.NewCliError(fmt.Errorf("erro ao alterar o tenant: %w", err).Error())
	}

	accessKeyId := auth.GetAccessKeyID()
	secretAccessKey := auth.GetSecretAccessKey()

	if accessKeyId != "" && secretAccessKey != "" {
		fmt.Print("🔐 This operation unset the current api key. \n\n")

		err := auth.SetAccessKeyID("")
		if err != nil {
			return cmdutils.NewCliError(fmt.Errorf("erro ao remover o access key id: %w", err).Error())
		}

		err = auth.SetSecretAccessKey("")
		if err != nil {
			return cmdutils.NewCliError(fmt.Errorf("erro ao remover o secret access key: %w", err).Error())
		}
	}

	beautiful.NewOutput(rawMode).PrintData(tokenInfo)

	return nil
}
