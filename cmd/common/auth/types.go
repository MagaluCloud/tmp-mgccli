package auth

import (
	"context"

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
