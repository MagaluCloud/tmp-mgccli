package clients

import (
	"github.com/spf13/cobra"
)

func NilIfNotChanged[T any](
	cmd *cobra.Command,
	flag string,
	target **T,
	value T,
) {
	if !cmd.Flags().Changed(flag) {
		*target = nil
	} else {
		*target = &value
	}
}

const (
	ID                               = "id"
	Name                             = "name"
	Description                      = "description"
	RedirectURIs                     = "redirect-uris"
	BackchannelLogoutSessionEnabled  = "backchannel-logout-session"
	ClientTermsURL                   = "client-term-url"
	ClientPrivacyTermURL             = "client-privacy-term-url"
	Audiences                        = "audiences"
	Email                            = "email"
	Reason                           = "request-reason"
	Icon                             = "icon"
	AccessTokenExp                   = "access-token-expiration"
	AlwaysRequireLogin               = "always-require-login"
	BackchannelLogoutURI             = "backchannel-logout-uri"
	OidcAudience                     = "oidc-audience"
	RefreshTokenCustomExpiresEnabled = "refresh-token-custom-expires-enabled"
	RefreshTokenExp                  = "refresh-token-exp"
	SupportURL                       = "support-url"
	GrantTypes                       = "grant-types"
)
