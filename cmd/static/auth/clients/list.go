package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/common/auth"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type Clients struct {
	UUID        string `json:"uuid,omitempty"`
	ClientID    string `json:"client_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	Scopes      []struct {
		UUID string `json:"uuid"`
		Name string `json:"name"`
	} `json:"scopes,omitempty"`
	ScopesDefault []struct {
		UUID string `json:"uuid"`
		Name string `json:"name"`
	} `json:"scopes_default,omitempty"`
	TermOfUse                        string   `json:"client_term_url,omitempty"`
	ClientPrivacyTermUrl             string   `json:"client_privacy_term_url,omitempty"`
	Audience                         []string `json:"audience,omitempty"`
	OidcAudience                     []string `json:"oidc_audience,omitempty"`
	AlwaysRequireLogin               bool     `json:"always_require_login,omitempty"`
	BackchannelLogoutSessionRequired bool     `json:"backchannel_logout_session_required,omitempty"`
	BackchannelLogoutUri             string   `json:"backchannel_logout_uri,omitempty"`
	RefreshTokenCustomExpiresEnabled bool     `json:"refresh_token_custom_expires_enabled,omitempty"`
	RefreshTokenExp                  int      `json:"refresh_token_exp,omitempty"`
	AccessTokenExp                   int      `json:"access_token_exp,omitempty"`
	RedirectURIs                     []string `json:"redirect_uris,omitempty"`
	Icon                             string   `json:"icon,omitempty"`
}

type ListClientResult struct {
	UUID                             string   `json:"uuid,omitempty"`
	ClientID                         string   `json:"client_id,omitempty"`
	Name                             string   `json:"name,omitempty"`
	Description                      string   `json:"description,omitempty"`
	Status                           string   `json:"client_approval_status,omitempty"`
	Scopes                           []string `json:"scopes,omitempty"`
	ScopesDefault                    []string `json:"scopes_default,omitempty"`
	TermOfUse                        string   `json:"term_of_use,omitempty"`
	ClientPrivacyTermUrl             string   `json:"client_privacy_term_url,omitempty"`
	Audiences                        []string `json:"audiences,omitempty"`
	OidcAudiences                    []string `json:"oidc_audience,omitempty"`
	AlwaysRequireLogin               bool     `json:"always_require_login"`
	BackchannelLogoutSessionEnabled  bool     `json:"backchannel_logout_session_enabled"`
	BackchannelLogoutUri             string   `json:"backchannel_logout_uri,omitempty"`
	RefreshTokenCustomExpiresEnabled bool     `json:"refresh_token_custom_expires_enabled"`
	RefreshTokenExp                  int      `json:"refresh_token_expiration,omitempty"`
	AccessTokenExp                   int      `json:"access_token_expiration,omitempty"`
	RedirectURIs                     []string `json:"redirect_uris,omitempty"`
	Icon                             string   `json:"icon,omitempty"`
}

// ListCommand cria o comando de listar todos os clientes
func ListCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "list",
		Short: manager.T("cli.auth.clients.list.short"),
		Long:  manager.T("cli.auth.clients.list.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runList(ctx, raw)
		},
	}

	return cmd
}

// runList executa o processo de exibir todos clientes
func runList(ctx context.Context, rawMode bool) error {
	clients, err := listClients(ctx)
	if err != nil {
		return fmt.Errorf("erro ao listar os clientes: %w", err)
	}

	result := []*ListClientResult{}

	for _, client := range clients {
		clientInfo := &ListClientResult{
			UUID:                             client.UUID,
			ClientID:                         client.ClientID,
			Name:                             client.Name,
			Description:                      client.Description,
			Status:                           client.Status,
			TermOfUse:                        client.TermOfUse,
			ClientPrivacyTermUrl:             client.ClientPrivacyTermUrl,
			Audiences:                        client.Audience,
			AlwaysRequireLogin:               client.AlwaysRequireLogin,
			OidcAudiences:                    client.OidcAudience,
			BackchannelLogoutSessionEnabled:  client.BackchannelLogoutSessionRequired,
			BackchannelLogoutUri:             client.BackchannelLogoutUri,
			RefreshTokenCustomExpiresEnabled: client.RefreshTokenCustomExpiresEnabled,
			RefreshTokenExp:                  client.RefreshTokenExp,
			AccessTokenExp:                   client.AccessTokenExp,
			RedirectURIs:                     client.RedirectURIs,
			Icon:                             client.Icon,
		}

		for _, scope := range client.Scopes {
			clientInfo.Scopes = append(clientInfo.Scopes, scope.Name)
		}
		for _, scope := range client.ScopesDefault {
			clientInfo.ScopesDefault = append(clientInfo.ScopesDefault, scope.Name)
		}

		result = append(result, clientInfo)
	}

	if len(result) == 0 {
		fmt.Println("Nenhum cliente foi criado ainda.")
	} else {
		beautiful.NewOutput(rawMode).PrintData(result)
	}

	return nil
}

func listClients(ctx context.Context) ([]*Clients, error) {
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

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		config.ClientsURLV2,
		nil,
	)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	clients := []*Clients{}
	if resp.StatusCode == http.StatusNoContent {
		return clients, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, cmdutils.NewHttpErrorFromResponse(resp, r)
	}

	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&clients); err != nil {
		return nil, err
	}

	return clients, nil
}
