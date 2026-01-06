package common

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/google/uuid"
	"github.com/magaluCloud/mgccli/beautiful"
	configPkg "github.com/magaluCloud/mgccli/cmd/common/config"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
)

type ACLPermission struct {
	ID string
}

type Options struct {
	GrantWrite *string
	Private    *bool
	PublicRead *bool
}

func SetACL(ctx context.Context, bucketName string, opts Options) error {
	if PermissionsIsEmpty(opts) {
		return nil
	}

	var grants []ACLPermission

	if opts.GrantWrite != nil && *opts.GrantWrite != "" {
		err := json.Unmarshal([]byte(*opts.GrantWrite), &grants)
		if err != nil {
			return fmt.Errorf("invalid --grant-write JSON: %w", err)
		}
	}

	err := validatePermissions(opts, grants)
	if err != nil {
		return err
	}

	config := ctx.Value(cmdutils.CXT_CONFIG_KEY).(configPkg.Config)

	region, err := config.Get("region")
	if err != nil {
		return fmt.Errorf("erro ao pegar a região: %w", err)
	}

	host, err := BuildHost(bucketName, region.Value.(string))
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, bucketURL.String(), nil)
	if err != nil {
		return err
	}

	if opts.PublicRead != nil && *opts.PublicRead {
		req.Header.Set("x-amz-acl", "public-read")
	}

	if opts.Private != nil && *opts.Private {
		req.Header.Set("x-amz-acl", "private")
	}

	if len(grants) > 0 {
		err = setGrantWriteHeader(req.Header, grants, region.Value.(string))
		if err != nil {
			return err
		}
	}

	resp, err := SendRequest(ctx, req, region.Value.(string))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return cmdutils.NewHttpErrorFromResponse(resp, req)
	}

	return nil
}

func validatePermissions(opts Options, grants []ACLPermission) error {
	for _, grant := range grants {
		if grant.ID == "" {
			return fmt.Errorf("ID for ACLPermission may not be empty")
		}
	}

	if opts.Private != nil && opts.PublicRead != nil {
		if *opts.Private && *opts.PublicRead {
			return fmt.Errorf("canned ACL cannot have more than one 'true' field: [private, public-read]")
		}
	}

	hasTrueOption := (opts.Private != nil && *opts.Private) || (opts.PublicRead != nil && *opts.PublicRead)

	if len(grants) > 0 && hasTrueOption {
		return fmt.Errorf("Specifying both Canned ACLs and Header Grants is not allowed")
	}

	return nil
}

func setGrantWriteHeader(header http.Header, grants []ACLPermission, region string) error {
	value := ""
	for i, permission := range grants {
		if i > 0 {
			value += ","
		}

		var tenantId, userProject string

		_, err := uuid.Parse(permission.ID)
		isTenantId := err == nil

		if isTenantId {
			tenantId = permission.ID
			userProject = userProjectFromTenantId(permission.ID, region)
		} else {
			userProject = permission.ID
			tenantId, err = tenantIdFromUserProject(permission.ID)
			if err != nil {
				return err
			}
		}

		value += fmt.Sprintf("id=%s,id=%s", tenantId, userProject)
	}

	if value == "" {
		return nil
	}

	header.Set("x-amz-grant-write", value)
	return nil
}

func tenantIdFromUserProject(userProject string) (string, error) {
	var userProjectRegex = regexp.MustCompile("cloud_(?P<region>[^_]+)_(?P<env>[^_]+).(?P<tenant_id>[^:]+):cloud_(?P<region1>[^_]+)_(?P<env1>[^_]+).(?P<tenant_id1>[^:]+)")

	match := userProjectRegex.FindStringSubmatch(userProject)
	for i, substr := range match {
		if userProjectRegex.SubexpNames()[i] == "tenant_id" {
			return substr, nil
		}
	}

	return "", fmt.Errorf("unable to find 'tenant_id' inside 'user_project' ACL ID permission: %q", userProject)
}

func userProjectFromTenantId(tenantId string, region string) string {
	pattern := fmt.Sprintf("cloud_%s_prod_%s", region, tenantId)
	return fmt.Sprintf("%s:%s", pattern, pattern)
}

func PrintGrantWriteHelp() {
	fmt.Fprintln(os.Stderr, "Grant Write format:")

	beautiful.NewOutput(false).PrintData([]ACLPermission{
		{ID: "string"},
	})
}

func PermissionsIsEmpty(opts Options) bool {
	hasTrueOption := (opts.Private != nil && *opts.Private) || (opts.PublicRead != nil && *opts.PublicRead)

	if opts.GrantWrite == nil && !hasTrueOption {
		return true
	}

	return false
}
