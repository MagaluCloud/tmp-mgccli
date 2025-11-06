package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// OAuthClient gerencia as requisições OAuth para autenticação
type OAuthClient struct {
	config       *Config
	httpClient   *http.Client
	codeVerifier *CodeVerifier
}

// NewOAuthClient cria uma nova instância do cliente OAuth
func NewOAuthClient(config *Config) (*OAuthClient, error) {
	verifier, err := NewVerifier()
	if err != nil {
		return nil, fmt.Errorf("failed to create code verifier: %w", err)
	}

	return &OAuthClient{
		config:       config,
		httpClient:   &http.Client{},
		codeVerifier: verifier,
	}, nil
}

// BuildAuthURL constrói a URL de autenticação com PKCE
func (c *OAuthClient) BuildAuthURL() (*url.URL, error) {
	authURL, err := url.Parse(c.config.AuthURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse auth URL: %w", err)
	}

	query := authURL.Query()
	query.Add("response_type", "code")
	query.Add("client_id", c.config.ClientID)
	query.Add("redirect_uri", c.config.RedirectURI)
	query.Add("code_challenge", c.codeVerifier.CodeChallengeS256())
	query.Add("code_challenge_method", "S256")
	query.Add("scope", strings.Join(c.config.Scopes, " "))
	query.Add("choose_tenants", "true")

	authURL.RawQuery = query.Encode()
	return authURL, nil
}

// ExchangeCodeForToken troca o código de autorização por tokens de acesso
func (c *OAuthClient) ExchangeCodeForToken(ctx context.Context, authCode string) (*TokenResponse, error) {
	if c.codeVerifier == nil {
		return nil, fmt.Errorf("code verifier not initialized")
	}

	data := url.Values{}
	data.Set("client_id", c.config.ClientID)
	data.Set("redirect_uri", c.config.RedirectURI)
	data.Set("grant_type", "authorization_code")
	data.Set("code", authCode)
	data.Set("code_verifier", c.codeVerifier.Value())

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.config.TokenURL,
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute token request: %w", err)
	}
	defer resp.Body.Close()

	// Ler o corpo da resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Verificar status da resposta
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Decodificar resposta
	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &tokenResp, nil
}
