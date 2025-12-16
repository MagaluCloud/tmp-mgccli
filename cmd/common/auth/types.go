package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	TenantIDWithType string `json:"tenant"`
	ScopesStr        string `json:"scope"`
	Email            string `json:"email"`
}

// TokenResponse representa a resposta do servidor de autenticação OAuth
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

// AuthResult representa o resultado de uma tentativa de autenticação
type AuthResult struct {
	Token *TokenResponse
	Error error
}

// LoginOptions configura opções para o processo de login
type LoginOptions struct {
	Headless bool // Login sem abrir navegador
	QRCode   bool // Exibir QR code para login
	Show     bool // Mostrar token de acesso após login
}

// AuthService define a interface para serviços de autenticação
type AuthService interface {
	// Login inicia o fluxo de autenticação OAuth
	Login(ctx context.Context, opts LoginOptions) (*TokenResponse, error)
}

// TemplateData representa os dados passados para o template HTML
type TemplateData struct {
	Title            string
	Lines            []string
	ErrorDescription string
}

// Tenant representa as informações de um tenant
type Tenant struct {
	UUID        string `json:"uuid"`
	Name        string `json:"legal_name"`
	Email       string `json:"email"`
	IsManaged   bool   `json:"is_managed"`
	IsDelegated bool   `json:"is_delegated"`
}

// TenantResult representa o retorno da alteração do tenant
type TenantResult struct {
	AccessToken  string `json:"access_token"`
	CreatedAt    int    `json:"created_at"`
	ExpiresIn    int    `json:"expires_in"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"scope_type"`
}

// TokenExchangeResult representa o resultado da troca de token
type TokenExchangeResult struct {
	TenantID     string    `json:"uuid"`
	CreatedAt    time.Time `json:"created_at"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Scope        []string  `json:"scope"`
}

// ApiKeys representa o retorno da requisição de API Keys
type ApiKeys struct {
	UUID          string  `json:"uuid"`
	Name          string  `json:"name"`
	ApiKey        string  `json:"api_key"`
	Description   string  `json:"description"`
	KeyPairID     string  `json:"key_pair_id"`
	KeyPairSecret string  `json:"key_pair_secret"`
	StartValidity string  `json:"start_validity"`
	EndValidity   *string `json:"end_validity,omitempty"`
	RevokedAt     *string `json:"revoked_at,omitempty"`
	TenantName    *string `json:"tenant_name,omitempty"`
	Tenant        struct {
		UUID      string `json:"uuid"`
		LegalName string `json:"legal_name"`
	} `json:"tenant"`
	Scopes []struct {
		UUID        string `json:"uuid"`
		Name        string `json:"name"`
		Title       string `json:"title"`
		ConsentText string `json:"consent_text"`
		Icon        string `json:"icon"`
		APIProduct  struct {
			UUID string `json:"uuid"`
			Name string `json:"name"`
		} `json:"api_product"`
	} `json:"scopes"`
}

type Scope struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	UUID  string `json:"uuid"`
}

type APIProduct struct {
	Name   string  `json:"name"`
	Scopes []Scope `json:"scopes"`
	UUID   string  `json:"uuid"`
}

type Platform struct {
	APIProducts []APIProduct `json:"api_products"`
	Name        string       `json:"name"`
	UUID        string       `json:"uuid"`
}

type PlatformsResponse []Platform

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
