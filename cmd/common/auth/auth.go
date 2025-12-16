package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/golang-jwt/jwt/v5"
	"github.com/magaluCloud/mgccli/cmd/common/structs"
	"github.com/magaluCloud/mgccli/cmd/common/workspace"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"gopkg.in/yaml.v3"
)

// access_key_id: ""
// access_token: ""
// current_environment: ""
// refresh_token: ""
// secret_access_key: ""

type AuthFile struct {
	AccessKeyID     string `yaml:"access_key_id"`
	AccessToken     string `yaml:"access_token"`
	RefreshToken    string `yaml:"refresh_token"`
	SecretAccessKey string `yaml:"secret_access_key"`
}

type Auth interface {
	GetAccessKeyID() string
	GetAccessToken(ctx context.Context) string
	GetRefreshToken() string
	GetSecretAccessKey() string
	GetService() *Service
	GetCurrentTenantID() (string, error)
	GetCurrentTenant(ctx context.Context) (*Tenant, error)
	GetScopes() (string, error)
	TokenClaims() (*TokenClaims, error)

	SetAccessToken(token string) error
	SetRefreshToken(token string) error
	SetSecretAccessKey(key string) error
	SetAccessKeyID(key string) error
	SetTenant(ctx context.Context, id string) (*TokenExchangeResult, error)

	ValidateToken() error
	RefreshToken(ctx context.Context) error

	Logout() error

	ListTenants(ctx context.Context) ([]*Tenant, error)
	ListApiKeys(ctx context.Context, showInvalidKeys bool) ([]*ApiKeys, error)

	CreateApiKey(ctx context.Context, name, description, expiration string, scopes []string) (*ApiKeyResult, error)
	RevokeApiKey(ctx context.Context, ID string) error
}

type authValue struct {
	authValue AuthFile
	workspace workspace.Workspace
	service   *Service
}

func NewAuth(workspace workspace.Workspace) Auth {
	authFile := path.Join(workspace.Dir(), "auth.yaml")
	authContent, err := structs.LoadFileToStruct[AuthFile](authFile)

	config := DefaultConfig()
	service := NewService(config)

	if err != nil {
		//TODO: Handle error
		panic(err)
	}
	return &authValue{workspace: workspace, authValue: authContent, service: service}
}

func (a *authValue) GetService() *Service {
	return a.service
}

func (a *authValue) GetAccessKeyID() string {
	return a.authValue.AccessKeyID
}

func (a *authValue) GetAccessToken(ctx context.Context) string {
	if a.authValue.AccessToken == "" {
		return ""
	}
	err := a.ValidateToken()
	if err != nil {
		if a.authValue.RefreshToken != "" {
			err := a.RefreshToken(ctx)
			if err != nil {
				return ""
			}
			return a.authValue.AccessToken
		}
	}
	return a.authValue.AccessToken
}

func (a *authValue) GetRefreshToken() string {
	return a.authValue.RefreshToken
}

func (a *authValue) GetSecretAccessKey() string {
	return a.authValue.SecretAccessKey
}

func (a *authValue) GetCurrentTenantID() (string, error) {
	claims, err := a.TokenClaims()
	if err != nil {
		return "", err
	}

	tenantId := claims.TenantIDWithType

	// Dot is a separator, Tenant will be <TenantType>.<ID>. We only want the ID
	dotIndex := strings.Index(tenantId, ".")

	if dotIndex != -1 {
		tenantId = tenantId[dotIndex+1:]
	}

	return tenantId, nil
}

func (a *authValue) GetCurrentTenant(ctx context.Context) (*Tenant, error) {
	currentTenantId, err := a.GetCurrentTenantID()
	if err != nil {
		return nil, err
	}

	tenants, err := a.ListTenants(ctx)
	if err != nil || len(tenants) == 0 {
		fmt.Printf("Não foi possível pegar as informações sobre o tenant atual, retornando apenas o ID.\nErro: %v\n\n", err)
		return &Tenant{UUID: currentTenantId}, nil
	}

	for _, tenant := range tenants {
		if tenant.UUID == currentTenantId {
			return tenant, nil
		}
	}

	return nil, fmt.Errorf("o ID (%s) do tenant atual não foi encontrado na lista de tenants", currentTenantId)
}

func (a *authValue) GetScopes() (string, error) {
	tokenClaims, err := a.TokenClaims()
	if err != nil {
		return "", err
	}

	return tokenClaims.ScopesStr, nil
}

func (a *authValue) SetAccessToken(token string) error {
	a.authValue.AccessToken = token
	return a.Write()
}

func (a *authValue) SetRefreshToken(token string) error {
	a.authValue.RefreshToken = token
	return a.Write()
}

func (a *authValue) SetSecretAccessKey(key string) error {
	a.authValue.SecretAccessKey = key
	return a.Write()
}

func (a *authValue) SetAccessKeyID(key string) error {
	a.authValue.AccessKeyID = key
	return a.Write()
}

func (a *authValue) Logout() error {
	a.SetAccessToken("")
	a.SetRefreshToken("")
	a.SetSecretAccessKey("")
	a.SetAccessKeyID("")
	return a.Write()
}

func (a *authValue) Write() error {
	data, err := yaml.Marshal(a.authValue)
	if err != nil {
		return err
	}

	err = os.MkdirAll(a.workspace.Dir(), 0744)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(a.workspace.Dir(), "auth.yaml"), data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (a *authValue) TokenClaims() (*TokenClaims, error) {
	if a.authValue.AccessToken == "" {
		return nil, fmt.Errorf("access token is not set")
	}

	tokenClaims := &TokenClaims{}
	tokenParser := jwt.NewParser()

	_, _, err := tokenParser.ParseUnverified(a.authValue.AccessToken, tokenClaims)
	if err != nil {
		return nil, err
	}

	return tokenClaims, nil
}

func (a *authValue) ValidateToken() error {
	//extract iat from token, if expires in less than 30 sec, run refresh operation
	tokenClaims, err := a.TokenClaims()
	if err != nil {
		return err
	}
	iat := tokenClaims.ExpiresAt.Time.Unix()
	if iat < time.Now().Unix()-60 {
		return fmt.Errorf("token expired")
	}
	return nil
}

func (a *authValue) RefreshToken(ctx context.Context) error {
	token, err := a.service.RefreshToken(ctx, a.authValue.RefreshToken)
	if err != nil {
		return err
	}
	a.authValue.AccessToken = token.AccessToken
	a.authValue.RefreshToken = token.RefreshToken
	return a.Write()
}

func (a *authValue) ListTenants(ctx context.Context) ([]*Tenant, error) {
	client, err := NewOAuthClient(a.service.config)
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
		a.service.config.TenantsListURL,
		nil,
	)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, cmdutils.NewHttpErrorFromResponse(resp, r)
	}

	defer resp.Body.Close()
	var result []*Tenant
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (a *authValue) SetTenant(ctx context.Context, id string) (*TokenExchangeResult, error) {
	return a.runTokenExchange(ctx, id)
}

func (a *authValue) runTokenExchange(
	ctx context.Context, tenantId string,
) (*TokenExchangeResult, error) {
	client, err := NewOAuthClient(a.service.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create OAuth client: %w", err)
	}

	httpClient := client.AuthenticatedHttpClientFromContext(ctx)
	if httpClient == nil {
		return nil, fmt.Errorf("programming error: unable to get HTTP Client from context")
	}

	scopes, err := a.GetScopes()
	if err != nil {
		return nil, err
	}

	data := map[string]any{
		"tenant": tenantId,
		"scopes": scopes,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(jsonData)
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, a.service.config.TokenExchangeURL, bodyReader)
	r.Header.Set("Content-Type", "application/json")

	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, cmdutils.NewHttpErrorFromResponse(resp, r)
	}

	payload := &TenantResult{}
	if err = json.NewDecoder(resp.Body).Decode(payload); err != nil {
		return nil, err
	}

	err = a.SetAccessToken(payload.AccessToken)
	if err != nil {
		return nil, err
	}

	err = a.SetRefreshToken(payload.RefreshToken)
	if err != nil {
		return nil, err
	}

	createdAt := time.Time(time.Unix(int64(payload.CreatedAt), 0))

	return &TokenExchangeResult{
		AccessToken:  payload.AccessToken,
		CreatedAt:    createdAt,
		TenantID:     tenantId,
		RefreshToken: payload.RefreshToken,
		Scope:        strings.Split(payload.Scope, " "),
	}, nil
}

func (a *authValue) ListApiKeys(ctx context.Context, showInvalidKeys bool) ([]*ApiKeys, error) {
	client, err := NewOAuthClient(a.service.config)
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
		a.service.config.ApiKeysURLV1,
		nil,
	)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	var apiKeys []*ApiKeys
	if resp.StatusCode == http.StatusNoContent {
		return apiKeys, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, cmdutils.NewHttpErrorFromResponse(resp, r)
	}

	defer resp.Body.Close()
	var result []*ApiKeys
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	for _, key := range result {
		if !showInvalidKeys && key.RevokedAt != nil {
			continue
		}

		if !showInvalidKeys && key.EndValidity != nil {
			expDate, _ := time.Parse(time.RFC3339, *key.EndValidity)
			if expDate.Before(time.Now()) {
				continue
			}
		}

		apiKeys = append(apiKeys, key)
	}

	return apiKeys, nil
}

func (a *authValue) CreateApiKey(ctx context.Context, name, description, expiration string, scopes []string) (*ApiKeyResult, error) {
	client, err := NewOAuthClient(a.service.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create OAuth client: %w", err)
	}

	httpClient := client.AuthenticatedHttpClientFromContext(ctx)
	if httpClient == nil {
		return nil, fmt.Errorf("programming error: unable to get HTTP Client from context")
	}

	scopesCreateList, err := a.processScopes(ctx, scopes)
	if err != nil {
		return nil, err
	}

	currentTenantID, err := a.GetCurrentTenantID()
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
		a.service.config.ApiKeysURLV2,
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

func (a *authValue) processScopes(ctx context.Context, scopes []string) ([]ScopesCreate, error) {
	client, err := NewOAuthClient(a.service.config)
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
		a.service.config.ScopesURL,
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

func (a *authValue) RevokeApiKey(ctx context.Context, ID string) error {
	client, err := NewOAuthClient(a.service.config)
	if err != nil {
		return fmt.Errorf("failed to create OAuth client: %w", err)
	}

	httpClient := client.AuthenticatedHttpClientFromContext(ctx)
	if httpClient == nil {
		return fmt.Errorf("programming error: unable to get HTTP Client from context")
	}

	url := fmt.Sprintf("%s/%s/revoke", a.service.config.ApiKeysURLV1, ID)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return cmdutils.NewHttpErrorFromResponse(resp, req)
	}

	return nil
}
