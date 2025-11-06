package auth

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// CallbackServer gerencia o servidor HTTP de callback OAuth
type CallbackServer struct {
	config   *Config
	client   *OAuthClient
	template *template.Template
	listener net.Listener
	server   *http.Server
	resultCh chan *AuthResult
	cancelCh chan struct{}
}

// NewCallbackServer cria uma nova instância do servidor de callback
func NewCallbackServer(config *Config, client *OAuthClient, tmpl *template.Template) (*CallbackServer, error) {
	listener, err := net.Listen("tcp", config.ListenAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on %s: %w", config.ListenAddr, err)
	}

	s := &CallbackServer{
		config:   config,
		client:   client,
		template: tmpl,
		listener: listener,
		resultCh: make(chan *AuthResult, 1),
		cancelCh: make(chan struct{}, 1),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", s.handleCallback)
	mux.HandleFunc("/term", s.handleTermsRedirect)
	mux.HandleFunc("/privacy", s.handlePrivacyRedirect)

	s.server = &http.Server{
		Addr:    config.ListenAddr,
		Handler: mux,
	}

	return s, nil
}

// Start inicia o servidor e retorna um canal para o resultado
func (s *CallbackServer) Start(ctx context.Context) <-chan *AuthResult {
	go s.serve(ctx)
	return s.resultCh
}

// Cancel cancela o servidor
func (s *CallbackServer) Cancel() {
	select {
	case s.cancelCh <- struct{}{}:
		// Sinal enviado
	default:
		// Canal já fechado ou sem receptor
	}
}

// serve executa o loop principal do servidor
func (s *CallbackServer) serve(ctx context.Context) {
	serverErrCh := make(chan error, 1)
	signalCh := make(chan os.Signal, 1)
	callbackCh := make(chan *AuthResult, 1)

	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signalCh)

	// Goroutine para aguardar eventos e fazer shutdown
	go func() {
		var result *AuthResult

		select {
		case err := <-serverErrCh:
			result = &AuthResult{Error: fmt.Errorf("server error: %w", err)}

		case sig := <-signalCh:
			result = &AuthResult{Error: fmt.Errorf("canceled by signal: %v", sig)}

		case <-s.cancelCh:
			result = &AuthResult{Error: fmt.Errorf("canceled by user")}

		case result = <-callbackCh:
			// Resultado do callback
		}

		s.shutdown(ctx)
		s.resultCh <- result
		close(s.resultCh)
	}()

	// Servir HTTP
	if err := s.server.Serve(s.listener); err != nil && err != http.ErrServerClosed {
		serverErrCh <- err
	}
}

// shutdown encerra o servidor graciosamente
func (s *CallbackServer) shutdown(ctx context.Context) {
	shutdownCtx, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()

	if err := s.server.Shutdown(shutdownCtx); err != nil {
		// Tentar fechar forçadamente
		_ = s.server.Close()
	}
}

// handleCallback processa o callback OAuth
func (s *CallbackServer) handleCallback(w http.ResponseWriter, r *http.Request) {
	authCode := r.URL.Query().Get("code")
	if authCode == "" {
		s.showErrorPage(w, "Código de autorização ausente", fmt.Errorf("no authorization code received"))
		s.resultCh <- &AuthResult{Error: fmt.Errorf("no authorization code in callback")}
		return
	}

	// Trocar código por token
	token, err := s.client.ExchangeCodeForToken(r.Context(), authCode)
	if err != nil {
		s.showErrorPage(w, "Falha ao obter token", err)
		s.resultCh <- &AuthResult{Error: fmt.Errorf("failed to exchange code for token: %w", err)}
		return
	}

	// Exibir página de sucesso
	if err := s.showSuccessPage(w); err != nil {
		// Log do erro, mas não falha a autenticação
		fmt.Fprintf(os.Stderr, "Warning: failed to show success page: %v\n", err)
	}

	s.resultCh <- &AuthResult{Token: token}
}

// handleTermsRedirect redireciona para os termos de uso
func (s *CallbackServer) handleTermsRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, s.config.TermsURL, http.StatusPermanentRedirect)
}

// handlePrivacyRedirect redireciona para a política de privacidade
func (s *CallbackServer) handlePrivacyRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, s.config.PrivacyURL, http.StatusPermanentRedirect)
}

// showErrorPage renderiza a página de erro
func (s *CallbackServer) showErrorPage(w http.ResponseWriter, description string, err error) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)

	data := TemplateData{
		Title:            "Error",
		ErrorDescription: fmt.Sprintf("%s: %v", description, err),
	}

	if err := s.renderTemplate(w, data); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to render error page: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// showSuccessPage renderiza a página de sucesso
func (s *CallbackServer) showSuccessPage(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := TemplateData{
		Title: "Success",
		Lines: []string{
			"Você fez login com sucesso na Magalu Cloud.",
			"Esta página pode ser fechada agora.",
		},
	}

	return s.renderTemplate(w, data)
}

// renderTemplate renderiza o template HTML
func (s *CallbackServer) renderTemplate(w io.Writer, data TemplateData) error {
	buf := bytes.NewBuffer(nil)
	if err := s.template.Execute(buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if _, err := io.Copy(w, buf); err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}

	return nil
}
