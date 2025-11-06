package auth

import (
	"os"
	"time"
)

// Config contém todas as configurações necessárias para autenticação OAuth
type Config struct {
	// OAuth Configuration
	ClientID    string
	AuthURL     string
	TokenURL    string
	RedirectURI string
	Scopes      []string

	// Server Configuration
	ListenAddr string
	Timeout    time.Duration

	// External Links
	TermsURL   string
	PrivacyURL string
}

// DefaultConfig retorna a configuração padrão para autenticação
func DefaultConfig() *Config {
	return &Config{
		ClientID:    "cw9qpaUl2nBiC8PVjNFN5jZeb2vTd_1S5cYs1FhEXh0",
		AuthURL:     "https://id.magalu.com/login",
		TokenURL:    "https://id.magalu.com/oauth/token",
		RedirectURI: "http://localhost:8095/callback",
		Scopes: []string{
			"mke.write", "api-consulta.read", "openid", "mcr.read", "dbaas.write",
			"cpo:read", "cpo:write", "evt:event-tr", "network.read", "network.write",
			"object-storage.write", "object-storage.read", "block-storage.read",
			"block-storage.write", "mke.read", "virtual-machine.read",
			"virtual-machine.write", "dbaas.read", "mcr.write", "gdb:ssh-pkey-r",
			"gdb:ssh-pkey-w", "pa:sa:manage", "lba.loadbalancer.read",
			"lba.loadbalancer.write", "gdb:azs-r", "lbaas.read", "lbaas.write",
			"iam:read", "iam:write",
		},
		ListenAddr: getListenAddr(),
		Timeout:    500 * time.Millisecond,
		TermsURL:   "https://magalu.cloud/termos-legais/termos-de-uso-magalu-cloud/",
		PrivacyURL: "https://magalu.cloud/termos-legais/politica-de-privacidade/",
	}
}

// getListenAddr retorna o endereço de escuta do servidor de callback
// Verifica a variável de ambiente MGC_LISTEN_ADDRESS ou usa o padrão
func getListenAddr() string {
	if addr := os.Getenv("MGC_LISTEN_ADDRESS"); addr != "" {
		return addr
	}
	return "127.0.0.1:8095"
}
