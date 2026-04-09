package buckets

import (
	"context"
	"encoding/json"
	"fmt"

	bws "github.com/geffersonFerraz/brazilian-words-sorter"
	"github.com/magaluCloud/mgccli/beautiful"
	objectstorage "github.com/magaluCloud/mgccli/cmd/common/object_storage"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/buckets/common"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	cobrautils "github.com/magaluCloud/mgccli/cobra_utils/flags"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type createFlags struct {
	Name             string
	NameIsPrefix     bool
	EnableVersioning bool
	CliListLinks     bool
	GrantWrite       string
	Private          bool
	PublicRead       bool
}

type createParams struct {
	Name             *string
	NameIsPrefix     *bool
	EnableVersioning *bool
	CliListLinks     bool
	GrantWrite       *string
	Private          *bool
	PublicRead       *bool
}

type result struct {
	Bucket           string                 `json:"bucket,omitempty"`
	BucketIsPrefix   bool                   `json:"bucket_is_prefix,omitempty"`
	EnableVersioning bool                   `json:"enable_versioning,omitempty"`
	GrantWrite       []common.ACLPermission `json:"grant_write,omitempty"`
	Private          *bool                  `json:"private,omitempty"`
	PublicRead       *bool                  `json:"public_read,omitempty"`
}

// CreateCommand cria o comando de criar um novo bucket
func CreateCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()

	var flags createFlags
	var params createParams

	cmd := &cobra.Command{
		Use:   "create",
		Short: manager.T("cli.auth.object_storage.buckets.create.short"),
		Long:  manager.T("cli.auth.object_storage.buckets.create.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if flags.GrantWrite == "help" {
				common.PrintGrantWriteHelp()
				return nil
			}

			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			params.CliListLinks = flags.CliListLinks

			cobrautils.NilIfNotChanged(cmd, "bucket", &params.Name, flags.Name)
			cobrautils.NilIfNotChanged(cmd, "bucket-is-prefix", &params.NameIsPrefix, flags.NameIsPrefix)
			cobrautils.NilIfNotChanged(cmd, "enable-versioning", &params.EnableVersioning, flags.EnableVersioning)
			cobrautils.NilIfNotChanged(cmd, "grant-write", &params.GrantWrite, flags.GrantWrite)
			cobrautils.NilIfNotChanged(cmd, "private", &params.Private, flags.Private)
			cobrautils.NilIfNotChanged(cmd, "public-read", &params.PublicRead, flags.PublicRead)

			return runCreate(ctx, params, raw)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "bucket", "", manager.T("cli.auth.object_storage.buckets.create.name"))
	cmd.Flags().BoolVar(&flags.NameIsPrefix, "bucket-is-prefix", false, manager.T("cli.auth.object_storage.buckets.create.name_is_prefix"))
	cmd.Flags().BoolVar(&flags.EnableVersioning, "enable-versioning", true, manager.T("cli.auth.object_storage.buckets.create.enable_versioning"))
	cmd.Flags().StringVar(&flags.GrantWrite, "grant-write", "", manager.T("cli.auth.object_storage.buckets.acl.grant_write"))
	cmd.Flags().BoolVar(&flags.Private, "private", false, manager.T("cli.auth.object_storage.buckets.acl.private"))
	cmd.Flags().BoolVar(&flags.PublicRead, "public-read", false, manager.T("cli.auth.object_storage.buckets.acl.public_read"))
	cmd.Flags().BoolVar(&flags.CliListLinks, "cli.list-links", false, manager.T("cli.auth.object_storage.buckets.create.cli_list_links"))

	return cmd
}

// runCreate executa o processo de criar um novo bucket
func runCreate(ctx context.Context, opts createParams, rawMode bool) error {
	objectStorageService, err := objectstorage.NewObjectStorage(ctx)
	if err != nil {
		return cmdutils.NewCliError(err.Error())
	}
	bucketService := objectStorageService.GetBucketService()

	if opts.CliListLinks {
		beautiful.NewOutput(rawMode).PrintTable(
			[]string{"Description", "Command"},
			[][]string{
				{"Delete an existing bucket", "delete"},
				{"List all existing buckets", "list"},
			},
		)

		return nil
	}

	if opts.Name == nil {
		return cmdutils.NewCliError("missing required flag: --bucket=string")
	}
	if opts.NameIsPrefix == nil {
		return cmdutils.NewCliError("missing required flag: --bucket-is-prefix=bool")
	}
	if opts.EnableVersioning == nil {
		return cmdutils.NewCliError("missing required flag: --enable-versioning=bool")
	}

	if *opts.NameIsPrefix {
		bwords := bws.BrazilianWords(3, "-")
		name := fmt.Sprintf("%s-%s", *opts.Name, bwords.Sort())
		opts.Name = &name
	}

	err = bucketService.Create(ctx, *opts.Name)
	if err != nil {
		return cmdutils.NewCliError(fmt.Sprintf("erro ao criar o bucket: %s", err.Error()))
	}

	if *opts.EnableVersioning {
		err := bucketService.EnableVersioning(ctx, *opts.Name)
		if err != nil {
			_ = bucketService.Delete(ctx, *opts.Name, true)

			return cmdutils.NewCliError(fmt.Sprintf("erro ao habilitar o versionamento do bucket: %s", err.Error()))
		}
	} else {
		err := bucketService.SuspendVersioning(ctx, *opts.Name)
		if err != nil {
			_ = bucketService.Delete(ctx, *opts.Name, true)

			return cmdutils.NewCliError(fmt.Sprintf("erro ao suspender o versionamento do bucket: %s", err.Error()))
		}
	}

	err = common.SetACL(ctx, *opts.Name, common.Options{
		GrantWrite: opts.GrantWrite,
		Private:    opts.Private,
		PublicRead: opts.PublicRead,
	})
	if err != nil {
		return err
	}

	var grantsFormatted []common.ACLPermission

	if opts.GrantWrite != nil && *opts.GrantWrite != "" {
		err := json.Unmarshal([]byte(*opts.GrantWrite), &grantsFormatted)
		if err != nil {
			return cmdutils.NewCliError(fmt.Sprintf("invalid --grant-write JSON: %s", err.Error()))
		}
	}

	response := result{
		Bucket:           *opts.Name,
		BucketIsPrefix:   *opts.NameIsPrefix,
		EnableVersioning: *opts.EnableVersioning,
		Private:          opts.Private,
		PublicRead:       opts.PublicRead,
	}

	if len(grantsFormatted) > 0 {
		response.GrantWrite = grantsFormatted
	}

	beautiful.NewOutput(rawMode).PrintData(response)

	return nil
}
