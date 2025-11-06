package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

const (
	// verifierLength define o tamanho do code verifier em bytes
	// RFC 7636 recomenda entre 43-128 caracteres após encoding
	verifierLength = 32
)

// CodeVerifier representa um PKCE code verifier usado no fluxo OAuth
// Ref: https://datatracker.ietf.org/doc/html/rfc7636
type CodeVerifier struct {
	value string
}

// NewVerifier cria um novo code verifier aleatório usando geração criptográfica segura
func NewVerifier() (*CodeVerifier, error) {
	b := make([]byte, verifierLength)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return newCodeVerifierFromBytes(b), nil
}

// newCodeVerifierFromBytes cria um code verifier a partir de bytes fornecidos
func newCodeVerifierFromBytes(b []byte) *CodeVerifier {
	return &CodeVerifier{
		value: base64URLEncode(b),
	}
}

// Value retorna o valor do code verifier
func (v *CodeVerifier) Value() string {
	return v.value
}

// CodeChallengeS256 gera o code challenge usando o método S256 (SHA256)
// Este é o método recomendado pela RFC 7636
func (v *CodeVerifier) CodeChallengeS256() string {
	h := sha256.New()
	h.Write([]byte(v.value))
	return base64URLEncode(h.Sum(nil))
}

// base64URLEncode codifica bytes em Base64 URL-safe sem padding
// conforme especificado na RFC 7636
func base64URLEncode(data []byte) string {
	encoded := base64.StdEncoding.EncodeToString(data)
	encoded = strings.ReplaceAll(encoded, "+", "-")
	encoded = strings.ReplaceAll(encoded, "/", "_")
	encoded = strings.ReplaceAll(encoded, "=", "")
	return encoded
}
