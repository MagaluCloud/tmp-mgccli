package cors

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type setOptions struct {
	Dst  string
	CORS string
}

// SetCommand cria o comando de configurar o CORS do bucket
func SetCommand(ctx context.Context, bucketService objSdk.BucketService) *cobra.Command {
	manager := i18n.GetInstance()
	var opts setOptions

	cmd := &cobra.Command{
		Use:   "set [dst]",
		Short: manager.T("cli.auth.object_storage.buckets.cors.set.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.CORS == "help" {
				printCORSHelp()
				return nil
			}

			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runSet(ctx, bucketService, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Dst, "dst", "", manager.T("cli.auth.object_storage.buckets.dst"))
	cmd.Flags().StringVar(&opts.CORS, "cors", "", manager.T("cli.auth.object_storage.buckets.cors.set.cors"))

	cmd.MarkFlagRequired("cors")

	return cmd
}

// runSet executa o processo de configurar o CORS do bucket
func runSet(ctx context.Context, bucketService objSdk.BucketService, args []string, opts setOptions, rawMode bool) error {
	if bucketService == nil {
		return nil
	}

	bucketName := opts.Dst

	if len(args) > 0 {
		bucketName = args[0]
	}

	if bucketName == "" {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o nome do bucket como argumento ou usar a flag --dst")

		return nil
	}

	corsData, err := resolveCORSInput(opts.CORS)
	if err != nil {
		return err
	}

	var cors *objSdk.CORSConfiguration
	if err := json.Unmarshal([]byte(corsData), &cors); err != nil {
		return fmt.Errorf("--cors JSON inválido: %w", err)
	}

	err = bucketService.SetCORS(ctx, bucketName, cors)
	if err != nil {
		return fmt.Errorf("erro ao setar a configuração do CORS: %w", err)
	}

	beautiful.NewOutput(rawMode).PrintData(cors)

	return nil
}

func printCORSHelp() {
	fmt.Fprintln(os.Stderr, "CORS format (JSON):")

	beautiful.NewOutput(false).PrintData(objSdk.CORSConfiguration{
		CORSRules: []objSdk.CORSRule{
			{
				AllowedOrigins: []string{"string"},
				AllowedMethods: []string{"string"},
				AllowedHeaders: []string{"string"},
				ExposeHeaders:  []string{"string"},
				MaxAgeSeconds:  0,
			},
		},
	})
}

func resolveCORSInput(value string) ([]byte, error) {
	// arquivo com @
	if after, ok := strings.CutPrefix(value, "@"); ok {
		path := after
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("não foi possível ler o arquivo de cors %q: %w", path, err)
		}
		return data, nil
	}

	// JSON inline
	return []byte(value), nil
}
