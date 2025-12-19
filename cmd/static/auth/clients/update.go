package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/charmbracelet/huh"
	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/common/auth"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type UpdateOptions struct {
	ID                               string
	Name                             string
	Description                      string
	RedirectURIs                     string
	BackchannelLogoutSessionEnabled  bool
	ClientTermsURL                   string
	ClientPrivacyTermURL             string
	Audiences                        string
	Reason                           string
	Icon                             string
	AccessTokenExp                   int
	AlwaysRequireLogin               bool
	BackchannelLogoutURI             string
	OidcAudience                     string
	RefreshTokenCustomExpiresEnabled bool
	RefreshTokenExp                  int
	SupportURL                       string
}

type UpdateClient struct {
	Name                             *string  `json:"name" jsonschema:"description=Name of new client,example=Client Name" mgc:"positional"`
	Description                      *string  `json:"description" jsonschema:"description=Description of new client,example=Client description" mgc:"positional"`
	RedirectURIs                     []string `json:"redirect_uris" jsonschema:"description=Redirect URIs (separated by space)" mgc:"positional"`
	Icon                             *string  `json:"icon,omitempty" jsonschema:"description=URL for client icon" mgc:"positional"`
	AccessTokenExp                   *int     `json:"access_token_exp,omitempty" jsonschema:"description=Access token expiration (in seconds),example=7200" mgc:"positional"`
	AlwaysRequireLogin               *bool    `json:"always_require_login,omitempty" jsonschema:"description=Must ignore active Magalu ID session and always require login,example=false" mgc:"positional"`
	ClientPrivacyTermUrl             *string  `json:"client_privacy_term_url" jsonschema:"description=URL to privacy term" mgc:"positional"`
	ClientTermUrl                    *string  `json:"client_term_url" jsonschema:"description=URL to terms of use" mgc:"positional"`
	Audience                         []string `json:"audience,omitempty" jsonschema:"description=Client audiences (separated by space),example=public" mgc:"positional"`
	OidcAudience                     []string `json:"oidc_audience,omitempty" jsonschema:"description=Audiences for ID token, should be the Client ID values" mgc:"positional"`
	BackchannelLogoutSessionEnabled  *bool    `json:"backchannel_logout_session_required,omitempty" jsonschema:"description=Client requires backchannel logout session,example=false" mgc:"positional"`
	BackchannelLogoutUri             *string  `json:"backchannel_logout_uri,omitempty" jsonschema:"description=Backchannel logout URI" mgc:"positional"`
	RefreshTokenCustomExpiresEnabled *bool    `json:"refresh_token_custom_expires_enabled,omitempty" jsonschema:"description=Use custom value for refresh token expiration,example=false" mgc:"positional"`
	RefreshTokenExp                  *int     `json:"refresh_token_exp,omitempty" jsonschema:"description=Custom refresh token expiration value (in seconds),example=15778476" mgc:"positional"`
	Reason                           *string  `json:"request_reason,omitempty" jsonschema:"description=Note to inform the reason for creating the client. Will help with the application approval process" mgc:"positional"`
	SupportUrl                       *string  `json:"support_url,omitempty" jsonschema:"description=URL for client support" mgc:"positional"`
}

type UpdateClientParams struct {
	ID                               string
	Name                             *string
	Description                      *string
	RedirectURIs                     *string
	BackchannelLogoutSessionEnabled  *bool
	ClientTermsURL                   *string
	ClientPrivacyTermURL             *string
	Audiences                        *string
	Reason                           *string
	Icon                             *string
	AccessTokenExp                   *int
	AlwaysRequireLogin               *bool
	BackchannelLogoutURI             *string
	OidcAudience                     *string
	RefreshTokenCustomExpiresEnabled *bool
	RefreshTokenExp                  *int
	SupportURL                       *string
}

type UpdateClientResult struct {
	UUID     string `json:"uuid,omitempty"`
	ClientID string `json:"client_id,omitempty"`
}

// UpdateCommand cria o comando de atualizar as informações de um cliente
func UpdateCommand(ctx context.Context) *cobra.Command {
	var opts UpdateOptions
	var params UpdateClientParams

	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:     "update",
		Short:   manager.T("cli.auth.clients.update.short"),
		Long:    manager.T("cli.auth.clients.update.long"),
		Example: `mgc auth clients update --access-token-expiration=7200 --audiences="public" --description="Client description" --name="Client Name" --refresh-token-exp=15552000`,
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			params = UpdateClientParams{
				ID: opts.ID,
			}

			NilIfNotChanged(cmd, Name, &params.Name, opts.Name)
			NilIfNotChanged(cmd, Description, &params.Description, opts.Description)
			NilIfNotChanged(cmd, RedirectURIs, &params.RedirectURIs, opts.RedirectURIs)
			NilIfNotChanged(cmd, BackchannelLogoutSessionEnabled, &params.BackchannelLogoutSessionEnabled, opts.BackchannelLogoutSessionEnabled)
			NilIfNotChanged(cmd, ClientTermsURL, &params.ClientTermsURL, opts.ClientTermsURL)
			NilIfNotChanged(cmd, ClientPrivacyTermURL, &params.ClientPrivacyTermURL, opts.ClientPrivacyTermURL)
			NilIfNotChanged(cmd, Audiences, &params.Audiences, opts.Audiences)
			NilIfNotChanged(cmd, Reason, &params.Reason, opts.Reason)
			NilIfNotChanged(cmd, Icon, &params.Icon, opts.Icon)
			NilIfNotChanged(cmd, AccessTokenExp, &params.AccessTokenExp, opts.AccessTokenExp)
			NilIfNotChanged(cmd, AlwaysRequireLogin, &params.AlwaysRequireLogin, opts.AlwaysRequireLogin)
			NilIfNotChanged(cmd, BackchannelLogoutURI, &params.BackchannelLogoutURI, opts.BackchannelLogoutURI)
			NilIfNotChanged(cmd, OidcAudience, &params.OidcAudience, opts.OidcAudience)
			NilIfNotChanged(cmd, RefreshTokenCustomExpiresEnabled, &params.RefreshTokenCustomExpiresEnabled, opts.RefreshTokenCustomExpiresEnabled)
			NilIfNotChanged(cmd, RefreshTokenExp, &params.RefreshTokenExp, opts.RefreshTokenExp)
			NilIfNotChanged(cmd, SupportURL, &params.SupportURL, opts.SupportURL)

			return runUpdate(ctx, params, raw)
		},
	}

	cmd.Flags().StringVar(&opts.ID, ID, "", manager.T("cli.auth.clients.update.uuid"))
	cmd.Flags().StringVar(&opts.Name, Name, "", manager.T("cli.auth.clients.create.name"))
	cmd.Flags().StringVar(&opts.Description, Description, "", manager.T("cli.auth.clients.create.description"))
	cmd.Flags().StringVar(&opts.RedirectURIs, RedirectURIs, "", manager.T("cli.auth.clients.create.redirect_uris"))
	cmd.Flags().BoolVar(&opts.BackchannelLogoutSessionEnabled, BackchannelLogoutSessionEnabled, false, manager.T("cli.auth.clients.create.backchannel_logout_session"))
	cmd.Flags().StringVar(&opts.ClientTermsURL, ClientTermsURL, "", manager.T("cli.auth.clients.create.client_term_url"))
	cmd.Flags().StringVar(&opts.ClientPrivacyTermURL, ClientPrivacyTermURL, "", manager.T("cli.auth.clients.create.client_privacy_term_url"))
	cmd.Flags().StringVar(&opts.Audiences, Audiences, "", manager.T("cli.auth.clients.create.audiences"))
	cmd.Flags().StringVar(&opts.Reason, Reason, "", manager.T("cli.auth.clients.create.request_reason"))
	cmd.Flags().StringVar(&opts.Icon, Icon, "", manager.T("cli.auth.clients.create.icon"))
	cmd.Flags().IntVar(&opts.AccessTokenExp, AccessTokenExp, 7200, manager.T("cli.auth.clients.create.access_token_expiration"))
	cmd.Flags().BoolVar(&opts.AlwaysRequireLogin, AlwaysRequireLogin, false, manager.T("cli.auth.clients.create.always_require_login"))
	cmd.Flags().StringVar(&opts.BackchannelLogoutURI, BackchannelLogoutURI, "", manager.T("cli.auth.clients.create.backchannel_logout_uri"))
	cmd.Flags().StringVar(&opts.OidcAudience, OidcAudience, "", manager.T("cli.auth.clients.create.oidc_audience"))
	cmd.Flags().BoolVar(&opts.RefreshTokenCustomExpiresEnabled, RefreshTokenCustomExpiresEnabled, false, manager.T("cli.auth.clients.create.refresh_token_custom_expires_enabled"))
	cmd.Flags().IntVar(&opts.RefreshTokenExp, RefreshTokenExp, 15552000, manager.T("cli.auth.clients.create.refresh_token_exp"))
	cmd.Flags().StringVar(&opts.SupportURL, SupportURL, "", manager.T("cli.auth.clients.create.support_url"))

	cmd.MarkFlagRequired(ID)

	return cmd
}

// runUpdate executa o processo de atualizar as informações de um cliente
func runUpdate(ctx context.Context, opts UpdateClientParams, rawMode bool) error {
	var confirm bool
	huh.NewConfirm().Title(fmt.Sprintf("This operation may disable your client %s until updates are approved by the ID Magalu. Do you wish to continue?", opts.ID)).
		Affirmative("Yes").
		Negative("No").Value(&confirm).Run()
	if !confirm {
		return nil
	}

	result, err := updateClient(ctx, opts)
	if err != nil {
		return fmt.Errorf("erro ao atualizar o cliente: %w", err)
	}

	beautiful.NewOutput(rawMode).PrintData(result)

	return nil
}

func updateClient(ctx context.Context, opts UpdateClientParams) (*UpdateClientResult, error) {
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

	clientPayload := UpdateClient{
		Name:                             opts.Name,
		Description:                      opts.Description,
		Icon:                             opts.Icon,
		ClientTermUrl:                    opts.ClientTermsURL,
		ClientPrivacyTermUrl:             opts.ClientPrivacyTermURL,
		AlwaysRequireLogin:               opts.AlwaysRequireLogin,
		BackchannelLogoutSessionEnabled:  opts.BackchannelLogoutSessionEnabled,
		BackchannelLogoutUri:             opts.BackchannelLogoutURI,
		AccessTokenExp:                   opts.AccessTokenExp,
		RefreshTokenCustomExpiresEnabled: opts.RefreshTokenCustomExpiresEnabled,
		RefreshTokenExp:                  opts.RefreshTokenExp,
		Reason:                           opts.Reason,
		SupportUrl:                       opts.SupportURL,
	}

	if opts.RedirectURIs != nil {
		clientPayload.RedirectURIs = cmdutils.StringToSlice(*opts.RedirectURIs, " ", true)
	}

	if opts.Audiences != nil {
		clientPayload.Audience = cmdutils.StringToSlice(*opts.Audiences, " ", true)
	}

	if opts.OidcAudience != nil {
		clientPayload.OidcAudience = cmdutils.StringToSlice(*opts.OidcAudience, " ", true)
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(clientPayload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s", config.PublicClientsURL, opts.ID)

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		url,
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

	if resp.StatusCode != http.StatusOK {
		return nil, cmdutils.NewHttpErrorFromResponse(resp, r)
	}

	defer resp.Body.Close()
	var result UpdateClientResult
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
