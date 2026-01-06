package acl

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"

	"github.com/magaluCloud/mgccli/beautiful"
	configPkg "github.com/magaluCloud/mgccli/cmd/common/config"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/buckets/common"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type getOptions struct {
	Dst string
}

type AccessControlPolicy struct {
	Owner             Owner             `xml:"Owner"`
	AccessControlList AccessControlList `xml:"AccessControlList"`
}

type Owner struct {
	ID          string `xml:"ID"`
	DisplayName string `xml:"DisplayName"`
}

type AccessControlList struct {
	Grants []Grant `xml:"Grant"`
}

type Grant struct {
	Grantee    Grantee `xml:"Grantee"`
	Permission string  `xml:"Permission"`
}

type Grantee struct {
	ID          string `xml:"ID"`
	DisplayName string `xml:"DisplayName"`
	URI         string `xml:"URI"`
}

// GetCommand cria o comando de retornar o ACL do bucket
func GetCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()
	var opts getOptions

	cmd := &cobra.Command{
		Use:   "get [dst]",
		Short: manager.T("cli.auth.object_storage.buckets.acl.get.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runGet(ctx, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Dst, "dst", "", manager.T("cli.auth.object_storage.buckets.dst"))

	return cmd
}

// runGet executa o processo de o ACL do bucket
func runGet(ctx context.Context, args []string, opts getOptions, rawMode bool) error {
	bucketName := opts.Dst

	if len(args) > 0 {
		bucketName = args[0]
	}

	if bucketName == "" {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o nome do bucket como argumento ou usar a flag --dst")

		return nil
	}

	config := ctx.Value(cmdutils.CXT_CONFIG_KEY).(configPkg.Config)

	region, err := config.Get("region")
	if err != nil {
		return fmt.Errorf("erro ao pegar a região: %w", err)
	}

	host, err := common.BuildHost(bucketName, region.Value.(string))
	if err != nil {
		return err
	}

	bucketURL, err := url.Parse(host)
	if err != nil {
		return err
	}

	query := bucketURL.Query()
	query.Add("acl", "")
	bucketURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, bucketURL.String(), nil)
	if err != nil {
		return err
	}

	resp, err := common.SendRequest(ctx, req, region.Value.(string))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return cmdutils.NewHttpErrorFromResponse(resp, req)
	}

	defer resp.Body.Close()
	var result AccessControlPolicy
	if err = xml.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	beautiful.NewOutput(rawMode).PrintData(result)

	return nil
}
