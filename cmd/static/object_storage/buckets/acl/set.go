package acl

import (
	"context"
	"fmt"
	"os"

	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/buckets/common"
	cobrautils "github.com/magaluCloud/mgccli/cobra_utils/flags"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type setFlags struct {
	Dst        string
	GrantWrite string
	Private    bool
	PublicRead bool
}

type setParams struct {
	Dst        string
	GrantWrite *string
	Private    *bool
	PublicRead *bool
}

// SetCommand cria o comando de definir o ACL do bucket
func SetCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()
	var flags setFlags
	var params setParams

	cmd := &cobra.Command{
		Use:   "set [dst]",
		Short: manager.T("cli.auth.object_storage.buckets.acl.set.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if flags.GrantWrite == "help" {
				common.PrintGrantWriteHelp()
				return nil
			}

			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			params.Dst = flags.Dst

			cobrautils.NilIfNotChanged(cmd, "grant-write", &params.GrantWrite, flags.GrantWrite)
			cobrautils.NilIfNotChanged(cmd, "private", &params.Private, flags.Private)
			cobrautils.NilIfNotChanged(cmd, "public-read", &params.PublicRead, flags.PublicRead)

			return runSet(ctx, args, params, raw)
		},
	}

	cmd.Flags().StringVar(&flags.Dst, "dst", "", manager.T("cli.auth.object_storage.buckets.dst"))
	cmd.Flags().StringVar(&flags.GrantWrite, "grant-write", "", manager.T("cli.auth.object_storage.buckets.acl.grant_write"))
	cmd.Flags().BoolVar(&flags.Private, "private", false, manager.T("cli.auth.object_storage.buckets.acl.private"))
	cmd.Flags().BoolVar(&flags.PublicRead, "public-read", false, manager.T("cli.auth.object_storage.buckets.acl.public_read"))

	return cmd
}

// runSet executa o processo de retornar o ACL do bucket
func runSet(ctx context.Context, args []string, opts setParams, rawMode bool) error {
	bucketName := opts.Dst

	if len(args) > 0 {
		bucketName = args[0]
	}

	if bucketName == "" {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o nome do bucket como argumento ou usar a flag --dst")

		return nil
	}

	aclOptions := common.Options{
		GrantWrite: opts.GrantWrite,
		Private:    opts.Private,
		PublicRead: opts.PublicRead,
	}

	if common.PermissionsIsEmpty(aclOptions) {
		return fmt.Errorf("needs to pass either grant permissions or canned info")
	}

	err := common.SetACL(ctx, bucketName, aclOptions)
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "✓ ACL definido com sucesso!")

	return nil
}
