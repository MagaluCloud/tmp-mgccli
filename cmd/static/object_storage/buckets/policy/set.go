package policy

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/beautiful"
	objectstorage "github.com/magaluCloud/mgccli/cmd/common/object_storage"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type setOptions struct {
	Dst    string
	Policy string
}

// SetCommand cria o comando de configurar a política do bucket
func SetCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()
	var opts setOptions

	cmd := &cobra.Command{
		Use:   "set [dst]",
		Short: manager.T("cli.auth.object_storage.buckets.policy.set.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Policy == "help" {
				printPolicyHelp()
				return nil
			}

			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runSet(ctx, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Dst, "dst", "", manager.T("cli.auth.object_storage.buckets.dst"))
	cmd.Flags().StringVar(&opts.Policy, "policy", "", manager.T("cli.auth.object_storage.buckets.policy.set.policy"))

	cmd.MarkFlagRequired("policy")

	return cmd
}

// runSet executa o processo de configurar a política do bucket
func runSet(ctx context.Context, args []string, opts setOptions, rawMode bool) error {
	objectStorageService, err := objectstorage.NewObjectStorage(ctx)
	if err != nil {
		return cmdutils.NewCliError(err.Error())
	}
	bucketService := objectStorageService.GetBucketService()

	bucketName := opts.Dst

	if len(args) > 0 {
		bucketName = args[0]
	}

	if bucketName == "" {
		return cmdutils.NewCliError("missing required flag: --dst=string")
	}

	policyData, err := resolvePolicyInput(opts.Policy)
	if err != nil {
		return err
	}

	var policy *objSdk.Policy
	if err := json.Unmarshal([]byte(policyData), &policy); err != nil {
		return cmdutils.NewCliError(fmt.Errorf("--policy JSON inválido: %w", err).Error())
	}

	err = bucketService.SetPolicy(ctx, bucketName, policy)
	if err != nil {
		return cmdutils.NewCliError(err.Error())
	}

	beautiful.NewOutput(rawMode).PrintData(policy)

	return nil
}

func printPolicyHelp() {
	fmt.Fprintln(os.Stderr, "Policy format (JSON):")

	beautiful.NewOutput(false).PrintData(objSdk.Policy{
		Id:      "string",
		Version: "string",
		Statement: []objSdk.Statement{
			{
				Sid:       "string",
				Effect:    "string",
				Principal: "any",
				Action:    "any",
				Resource:  "any",
			},
		},
	})
}

func resolvePolicyInput(value string) ([]byte, error) {
	// arquivo com @
	if after, ok := strings.CutPrefix(value, "@"); ok {
		path := after
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("não foi possível ler o arquivo de policy %q: %w", path, err)
		}
		return data, nil
	}

	// JSON inline
	return []byte(value), nil
}
