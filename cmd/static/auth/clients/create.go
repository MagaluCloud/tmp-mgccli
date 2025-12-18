package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/common/auth"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type CreateOptions struct {
	Name                             string
	Description                      string
	RedirectURIs                     string
	BackchannelLogoutSessionEnabled  bool
	ClientTermsURL                   string
	ClientPrivacyTermURL             string
	Audiences                        string
	Email                            string
	Reason                           string
	Icon                             string
	AccessTokenExp                   int
	AlwaysRequireLogin               bool
	BackchannelLogoutURI             string
	OidcAudience                     string
	RefreshTokenCustomExpiresEnabled bool
	RefreshTokenExp                  int
	SupportURL                       string
	GrantTypes                       string
}

type CreateClientParams struct {
	Name                             string
	Description                      string
	RedirectURIs                     string
	BackchannelLogoutSessionEnabled  *bool
	ClientTermsURL                   string
	ClientPrivacyTermURL             string
	Audiences                        *string
	Email                            *string
	Reason                           *string
	Icon                             *string
	AccessTokenExp                   int
	AlwaysRequireLogin               *bool
	BackchannelLogoutURI             *string
	OidcAudience                     *string
	RefreshTokenCustomExpiresEnabled *bool
	RefreshTokenExp                  int
	SupportURL                       *string
	GrantTypes                       *string
}

type createClientScopes struct {
	UUID     string `json:"id"`
	Reason   string `json:"request_reason"`
	Optional bool   `json:"optional"`
}

type CreateClient struct {
	Name                             string               `json:"name" jsonschema:"description=Name of new client,example=Client Name" mgc:"positional"`
	Description                      string               `json:"description" jsonschema:"description=Description of new client,example=Client description" mgc:"positional"`
	Scopes                           []createClientScopes `json:"scopes" jsonschema:"description=List of scopes (separated by space),example=openid profile" mgc:"positional"`
	RedirectURIs                     []string             `json:"redirect_uris" jsonschema:"description=Redirect URIs (separated by space)" mgc:"positional"`
	Icon                             *string              `json:"icon,omitempty" jsonschema:"description=URL for client icon" mgc:"positional"`
	AccessTokenExp                   int                  `json:"access_token_exp,omitempty" jsonschema:"description=Access token expiration (in seconds),example=7200" mgc:"positional"`
	AlwaysRequireLogin               *bool                `json:"always_require_login,omitempty" jsonschema:"description=Must ignore active Magalu ID session and always require login,example=false" mgc:"positional"`
	ClientPrivacyTermUrl             string               `json:"client_privacy_term_url" jsonschema:"description=URL to privacy term" mgc:"positional"`
	ClientTermUrl                    string               `json:"client_term_url" jsonschema:"description=URL to terms of use" mgc:"positional"`
	Audience                         []string             `json:"audience,omitempty" jsonschema:"description=Client audiences (separated by space),example=public" mgc:"positional"`
	BackchannelLogoutSessionEnabled  *bool                `json:"backchannel_logout_session_required,omitempty" jsonschema:"description=Client requires backchannel logout session,example=false" mgc:"positional"`
	BackchannelLogoutUri             *string              `json:"backchannel_logout_uri,omitempty" jsonschema:"description=Backchannel logout URI" mgc:"positional"`
	OidcAudience                     []string             `json:"oidc_audience,omitempty" jsonschema:"description=Audiences for ID token, should be the Client ID values" mgc:"positional"`
	RefreshTokenCustomExpiresEnabled *bool                `json:"refresh_token_custom_expires_enabled,omitempty" jsonschema:"description=Use custom value for refresh token expiration,example=false" mgc:"positional"`
	RefreshTokenExp                  int                  `json:"refresh_token_exp,omitempty" jsonschema:"description=Custom refresh token expiration value (in seconds),example=15778476" mgc:"positional"`
	Reason                           string               `json:"request_reason,omitempty" jsonschema:"description=Note to inform the reason for creating the client. Will help with the application approval process" mgc:"positional"`
	SupportUrl                       *string              `json:"support_url,omitempty" jsonschema:"description=URL for client support" mgc:"positional"`
	GrantTypes                       []string             `json:"grant_types,omitempty" jsonschema:"description=Grant types the client can request for token generation (separated by space)" mgc:"positional"`
	Email                            *string              `json:"email,omitempty" jsonschema:"description=Email for client support" mgc:"positional"`
}

type CreateClientResult struct {
	UUID         string `json:"uuid,omitempty"`
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
}

// CreateCommand cria o comando de criar um novo cliente
func CreateCommand(ctx context.Context) *cobra.Command {
	var opts CreateOptions
	var params CreateClientParams

	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "create",
		Short: manager.T("cli.auth.clients.create.short"),
		Long:  manager.T("cli.auth.clients.create.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			params = CreateClientParams{
				Name:                 opts.Name,
				Description:          opts.Description,
				ClientTermsURL:       opts.ClientTermsURL,
				ClientPrivacyTermURL: opts.ClientPrivacyTermURL,
				RedirectURIs:         opts.RedirectURIs,
				AccessTokenExp:       opts.AccessTokenExp,
				RefreshTokenExp:      opts.RefreshTokenExp,
			}

			NilIfNotChanged(cmd, BackchannelLogoutSessionEnabled, &params.BackchannelLogoutSessionEnabled, opts.BackchannelLogoutSessionEnabled)
			NilIfNotChanged(cmd, Audiences, &params.Audiences, opts.Audiences)
			NilIfNotChanged(cmd, Reason, &params.Reason, opts.Reason)
			NilIfNotChanged(cmd, Icon, &params.Icon, opts.Icon)
			NilIfNotChanged(cmd, AlwaysRequireLogin, &params.AlwaysRequireLogin, opts.AlwaysRequireLogin)
			NilIfNotChanged(cmd, BackchannelLogoutURI, &params.BackchannelLogoutURI, opts.BackchannelLogoutURI)
			NilIfNotChanged(cmd, OidcAudience, &params.OidcAudience, opts.OidcAudience)
			NilIfNotChanged(cmd, RefreshTokenCustomExpiresEnabled, &params.RefreshTokenCustomExpiresEnabled, opts.RefreshTokenCustomExpiresEnabled)
			NilIfNotChanged(cmd, SupportURL, &params.SupportURL, opts.SupportURL)
			NilIfNotChanged(cmd, Email, &params.Email, opts.Email)
			NilIfNotChanged(cmd, GrantTypes, &params.GrantTypes, opts.GrantTypes)

			return runCreate(ctx, params, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Name, Name, "", manager.T("cli.auth.clients.create.name"))
	cmd.Flags().StringVar(&opts.Description, Description, "", manager.T("cli.auth.clients.create.description"))
	cmd.Flags().StringVar(&opts.RedirectURIs, RedirectURIs, "", manager.T("cli.auth.clients.create.redirect_uris"))
	cmd.Flags().BoolVar(&opts.BackchannelLogoutSessionEnabled, BackchannelLogoutSessionEnabled, false, manager.T("cli.auth.clients.create.backchannel_logout_session"))
	cmd.Flags().StringVar(&opts.ClientTermsURL, ClientTermsURL, "", manager.T("cli.auth.clients.create.client_term_url"))
	cmd.Flags().StringVar(&opts.ClientPrivacyTermURL, ClientPrivacyTermURL, "", manager.T("cli.auth.clients.create.client_privacy_term_url"))
	cmd.Flags().StringVar(&opts.Audiences, Audiences, "", manager.T("cli.auth.clients.create.audiences"))
	cmd.Flags().StringVar(&opts.Email, Email, "", manager.T("cli.auth.clients.create.email"))
	cmd.Flags().StringVar(&opts.Reason, Reason, "", manager.T("cli.auth.clients.create.request_reason"))
	cmd.Flags().StringVar(&opts.Icon, Icon, "", manager.T("cli.auth.clients.create.icon"))
	cmd.Flags().IntVar(&opts.AccessTokenExp, AccessTokenExp, 7200, manager.T("cli.auth.clients.create.access_token_expiration"))
	cmd.Flags().BoolVar(&opts.AlwaysRequireLogin, AlwaysRequireLogin, false, manager.T("cli.auth.clients.create.always_require_login"))
	cmd.Flags().StringVar(&opts.BackchannelLogoutURI, BackchannelLogoutURI, "", manager.T("cli.auth.clients.create.backchannel_logout_uri"))
	cmd.Flags().StringVar(&opts.OidcAudience, OidcAudience, "", manager.T("cli.auth.clients.create.oidc_audience"))
	cmd.Flags().BoolVar(&opts.RefreshTokenCustomExpiresEnabled, RefreshTokenCustomExpiresEnabled, false, manager.T("cli.auth.clients.create.refresh_token_custom_expires_enabled"))
	cmd.Flags().IntVar(&opts.RefreshTokenExp, RefreshTokenExp, 15552000, manager.T("cli.auth.clients.create.refresh_token_exp"))
	cmd.Flags().StringVar(&opts.SupportURL, SupportURL, "", manager.T("cli.auth.clients.create.support_url"))
	cmd.Flags().StringVar(&opts.GrantTypes, GrantTypes, "", manager.T("cli.auth.clients.create.grant_types"))

	cmd.MarkFlagRequired(Name)
	cmd.MarkFlagRequired(ClientTermsURL)
	cmd.MarkFlagRequired(ClientPrivacyTermURL)
	cmd.MarkFlagRequired(Description)
	cmd.MarkFlagRequired(RedirectURIs)

	return cmd
}

// runCreate executa o processo de criar um novo cliente
func runCreate(ctx context.Context, opts CreateClientParams, rawMode bool) error {

	if opts.BackchannelLogoutSessionEnabled != nil && *opts.BackchannelLogoutSessionEnabled && opts.BackchannelLogoutURI == nil {
		return fmt.Errorf("backchannel-logout-uri is required when backchannel-logout-session is true")
	}

	result, err := createClient(ctx, opts)
	if err != nil {
		return fmt.Errorf("erro ao criar o cliente: %w", err)
	}

	fmt.Fprintln(os.Stderr, "Client created successfully! We'll analise your requisition and approve your client. You can check the approval status using client list command.")
	beautiful.NewOutput(rawMode).PrintData(result)

	return nil
}

func createClient(ctx context.Context, opts CreateClientParams) (*CreateClientResult, error) {
	authCtx := ctx.Value(cmdutils.CTX_AUTH_KEY).(auth.Auth)

	config := authCtx.GetConfig()

	client, err := auth.NewOAuthClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create OAuth client: %w", err)
	}

	httpClient := client.AuthenticatedHttpClientFromContext(ctx)
	if httpClient == nil {
		return nil, fmt.Errorf("programming error: unable to get HTTP Client from context")
	}

	clientPayload := CreateClient{
		Name:                             opts.Name,
		Description:                      opts.Description,
		ClientPrivacyTermUrl:             opts.ClientPrivacyTermURL,
		ClientTermUrl:                    opts.ClientTermsURL,
		BackchannelLogoutSessionEnabled:  opts.BackchannelLogoutSessionEnabled,
		AccessTokenExp:                   opts.AccessTokenExp,
		AlwaysRequireLogin:               opts.AlwaysRequireLogin,
		RefreshTokenCustomExpiresEnabled: opts.RefreshTokenCustomExpiresEnabled,
		RefreshTokenExp:                  opts.RefreshTokenExp,
		Icon:                             opts.Icon,
		BackchannelLogoutUri:             opts.BackchannelLogoutURI,
		SupportUrl:                       opts.SupportURL,
		Email:                            opts.Email,
	}
	clientPayload.RedirectURIs = cmdutils.StringToSlice(opts.RedirectURIs, " ", true)

	if opts.Reason == nil {
		opts.Reason = new(string)
		*opts.Reason = "Created by MGCCLI"
	}

	// Scopes fixos
	clientPayload.Scopes = []createClientScopes{{
		UUID:     config.PublicClientsScopeIDs["openid"],
		Reason:   *opts.Reason,
		Optional: true,
	}, {
		UUID:     config.PublicClientsScopeIDs["profile"],
		Reason:   *opts.Reason,
		Optional: true,
	}}

	if opts.Audiences != nil {
		clientPayload.Audience = cmdutils.StringToSlice(*opts.Audiences, " ", true)
	}

	if opts.OidcAudience != nil {
		clientPayload.OidcAudience = cmdutils.StringToSlice(*opts.OidcAudience, " ", true)
	}

	if opts.GrantTypes != nil {
		clientPayload.GrantTypes = cmdutils.StringToSlice(*opts.GrantTypes, " ", true)
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(clientPayload)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		config.PublicClientsURL,
		&buf,
	)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, cmdutils.NewHttpErrorFromResponse(resp, r)
	}

	defer resp.Body.Close()
	var result CreateClientResult
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
