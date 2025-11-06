package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"

	_ "embed"

	"github.com/pkg/browser"
)

var (
	//go:embed html.template
	htmlTemplateContent string
)

// Service implementa a interface AuthService para autenticação OAuth
type Service struct {
	config *Config
}

// NewService cria uma nova instância do serviço de autenticação
func NewService(config *Config) *Service {
	return &Service{
		config: config,
	}
}

// Login executa o fluxo de autenticação OAuth com as opções fornecidas
func (s *Service) Login(ctx context.Context, opts LoginOptions) (*TokenResponse, error) {
	if opts.QRCode {
		return s.qrCodeLogin(ctx)
	}

	if opts.Headless {
		return s.headlessLogin(ctx)
	}

	return s.browserLogin(ctx, opts.Show)
}

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	if refreshToken == "" {
		return nil, fmt.Errorf("RefreshToken is not set")
	}

	data := url.Values{}
	data.Set("client_id", s.config.ClientID)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, s.config.TokenURL, strings.NewReader(data.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(r)
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

// browserLogin executa o fluxo de login padrão abrindo o navegador
func (s *Service) browserLogin(ctx context.Context, showToken bool) (*TokenResponse, error) {
	// Preparar template HTML
	tmpl, err := template.New("html").Parse(htmlTemplateContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML template: %w", err)
	}

	// Criar cliente OAuth
	client, err := NewOAuthClient(s.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create OAuth client: %w", err)
	}

	// Iniciar servidor de callback
	server, err := NewCallbackServer(s.config, client, tmpl)
	if err != nil {
		return nil, fmt.Errorf("failed to start callback server: %w", err)
	}

	resultCh := server.Start(ctx)
	defer server.Cancel()

	// Construir URL de autenticação
	authURL, err := client.BuildAuthURL()
	if err != nil {
		return nil, fmt.Errorf("failed to build auth URL: %w", err)
	}

	// Abrir navegador
	fmt.Printf("Abrindo navegador em: %s://%s\n", authURL.Scheme, authURL.Host)
	if err := browser.OpenURL(authURL.String()); err != nil {
		fmt.Printf("Não foi possível abrir o navegador automaticamente.\n")
		fmt.Printf("Por favor, abra manualmente: %s\n", authURL.String())
	}

	// Aguardar resultado
	result := <-resultCh
	if result.Error != nil {
		return nil, result.Error
	}

	// Exibir token se solicitado
	if showToken && result.Token != nil {
		fmt.Printf("\nAccess Token: %s\n", result.Token.AccessToken)
		if result.Token.RefreshToken != "" {
			fmt.Printf("Refresh Token: %s\n", result.Token.RefreshToken)
		}
	}

	return result.Token, nil
}

// headlessLogin executa o fluxo de login sem abrir navegador
// TODO: Implementar fluxo headless (device flow)
func (s *Service) headlessLogin(ctx context.Context) (*TokenResponse, error) {
	return nil, fmt.Errorf("headless login not implemented yet")
}

// qrCodeLogin executa o fluxo de login exibindo um QR code
// TODO: Implementar fluxo com QR code
func (s *Service) qrCodeLogin(ctx context.Context) (*TokenResponse, error) {
	return nil, fmt.Errorf("QR code login not implemented yet")
}
