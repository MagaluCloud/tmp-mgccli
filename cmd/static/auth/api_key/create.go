package apikey

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	authPkg "github.com/magaluCloud/mgccli/cmd/common/auth"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

// CreateOptions configura opções para o processo de criar uma API Key
type CreateOptions struct {
	Name        string
	Description string
	Scopes      []string
	Expiration  string
}

type ScopesCreate struct {
	ID            string `json:"id"`
	RequestReason string `json:"request_reason"`
}

type CreateApiKey struct {
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	TenantID      string         `json:"tenant_id"`
	ScopesList    []ScopesCreate `json:"scopes"`
	StartValidity string         `json:"start_validity"`
	EndValidity   string         `json:"end_validity"`
}

type ApiKeyResult struct {
	UUID string `json:"uuid,omitempty"`
	Used bool   `json:"used,omitempty"`
}

type APIProduct struct {
	Name   string          `json:"name"`
	Scopes []authPkg.Scope `json:"scopes"`
	UUID   string          `json:"uuid"`
}

type Platform struct {
	APIProducts []APIProduct `json:"api_products"`
	Name        string       `json:"name"`
	UUID        string       `json:"uuid"`
}

type PlatformsResponse []Platform

// CreateCommand cria o comando de criar uma API Key
func CreateCommand(ctx context.Context) *cobra.Command {
	var opts CreateOptions

	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:     "create",
		Short:   manager.T("cli.auth.api_key.create.short"),
		Long:    manager.T("cli.auth.api_key.create.long"),
		Example: `mgc auth api-key create --description="created from MGC CLI" --expiration="2024-11-07" --name="My MGC Key" --scopes=["dbaas.read"]`,
		RunE: func(cmd *cobra.Command, args []string) error {

			return runCreate(ctx, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", manager.T("cli.auth.api_key.create.name"))
	cmd.Flags().StringVar(&opts.Description, "description", "", manager.T("cli.auth.api_key.create.description"))
	cmd.Flags().StringSliceVar(&opts.Scopes, "scopes", nil, manager.T("cli.auth.api_key.create.scopes"))
	cmd.Flags().StringVar(&opts.Expiration, "expiration", "", manager.T("cli.auth.api_key.create.expiration"))

	cmd.MarkFlagRequired("name")

	return cmd
}

// runCreate executa o processo de criar uma API Key
func runCreate(ctx context.Context, opts CreateOptions) error {
	apiKey, err := createApiKey(ctx, opts.Name, opts.Description, opts.Expiration, opts.Scopes)
	if err != nil {
		return fmt.Errorf("erro ao criar a API Key: %w", err)
	}

	fmt.Printf("API Key criada com sucesso!\nID: %s\n", apiKey.UUID)

	return nil
}

func createApiKey(ctx context.Context, name, description, expiration string, scopes []string) (*ApiKeyResult, error) {
	authCtx := ctx.Value(cmdutils.CTX_AUTH_KEY).(authPkg.Auth)

	config := authCtx.GetConfig()

	client, err := authPkg.NewOAuthClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create OAuth client: %w", err)
	}

	httpClient := client.AuthenticatedHttpClientFromContext(ctx)
	if httpClient == nil {
		return nil, fmt.Errorf("programming error: unable to get HTTP Client from context")
	}

	scopesCreateList, err := processScopes(ctx, scopes)
	if err != nil {
		return nil, err
	}

	currentTenantID, err := authCtx.GetCurrentTenantID()
	if err != nil {
		return nil, err
	}

	newApiKey := &CreateApiKey{
		Name:          name,
		TenantID:      currentTenantID,
		ScopesList:    scopesCreateList,
		StartValidity: time.Now().Format(time.DateOnly),
		Description:   description,
	}

	if expiration != "" {
		if _, err = time.Parse(time.DateOnly, expiration); err != nil {
			return nil, fmt.Errorf("invalid date format for expiration, use YYYY-MM-DD")
		}

		newApiKey.EndValidity = expiration
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(newApiKey)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		config.ApiKeysURLV2,
		&buf,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, cmdutils.NewHttpErrorFromResponse(resp, req)
	}

	defer resp.Body.Close()

	var result ApiKeyResult
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func processScopes(ctx context.Context, scopes []string) ([]ScopesCreate, error) {
	authCtx := ctx.Value(cmdutils.CTX_AUTH_KEY).(authPkg.Auth)

	config := authCtx.GetConfig()

	client, err := authPkg.NewOAuthClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create OAuth client: %w", err)
	}

	httpClient := client.AuthenticatedHttpClientFromContext(ctx)
	if httpClient == nil {
		return nil, fmt.Errorf("programming error: unable to get HTTP Client from context")
	}

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		config.ScopesURL,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	var scopesListFile PlatformsResponse

	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&scopesListFile); err != nil {
		return nil, err
	}

	scopesTitleMap := make(map[string]string)
	scopesNameMap := make(map[string]string)

	for _, company := range scopesListFile {
		if company.Name == "Magalu Cloud" {
			for _, product := range company.APIProducts {
				for _, scope := range product.Scopes {
					scopeName := product.Name + " [" + scope.Name + "]" + " - " + scope.Title
					scopesTitleMap[scopeName] = scope.UUID
					scopesNameMap[strings.ToLower(scope.Name)] = scope.UUID
				}
			}
		}
	}

	var scopesCreateList []ScopesCreate
	var invalidScopes []string

	if len(scopes) > 0 {
		for _, scope := range scopes {
			if id, ok := scopesNameMap[strings.ToLower(scope)]; ok {
				scopesCreateList = append(scopesCreateList, ScopesCreate{
					ID: id,
				})
			} else {
				invalidScopes = append(invalidScopes, scope)
			}
		}

		if len(invalidScopes) > 0 {
			return nil, fmt.Errorf("invalid scopes: %s", strings.Join(invalidScopes, ", "))
		}
	} else {
		options := []huh.Option[string]{}
		for title, id := range scopesTitleMap {
			options = append(options, huh.NewOption(title, id))
		}

		var selectedScopes []string

		multiSelect := huh.NewMultiSelect[string]()
		multiSelect.Title("Scopes:")
		multiSelect.Description("enter: confirm | space: select | ctrl + a: select/unselect all | /: to filter")
		multiSelect.Options(options...)
		multiSelect.Height(14)
		multiSelect.Filterable(true)
		multiSelect.Value(&selectedScopes)
		err = multiSelect.Run()
		if err != nil {
			return nil, cmdutils.NewCliError(err.Error())
		}

		if len(selectedScopes) == 0 {
			return nil, fmt.Errorf("nenhum scope selecionado")
		}

		for _, scopeID := range selectedScopes {
			scopesCreateList = append(scopesCreateList, ScopesCreate{
				ID: scopeID,
			})
		}
	}

	return scopesCreateList, nil
}
