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
