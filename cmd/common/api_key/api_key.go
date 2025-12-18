package apikey

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/magaluCloud/mgccli/cmd/common/auth"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
)

type ApiKey interface {
	List(ctx context.Context, showInvalidKeys bool) ([]*ApiKeys, error)
}

type apiKey struct {
	auth auth.Auth
}

func NewApiKey(a auth.Auth) ApiKey {
	return &apiKey{auth: a}
}

func (a *apiKey) List(ctx context.Context, showInvalidKeys bool) ([]*ApiKeys, error) {
	config := a.auth.GetConfig()

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
		config.ApiKeysURLV1,
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

	var apiKeys []*ApiKeys
	if resp.StatusCode == http.StatusNoContent {
		return apiKeys, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, cmdutils.NewHttpErrorFromResponse(resp, r)
	}

	defer resp.Body.Close()
	var result []*ApiKeys
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	for _, key := range result {
		if !showInvalidKeys && key.RevokedAt != nil {
			continue
		}

		if !showInvalidKeys && key.EndValidity != nil {
			expDate, _ := time.Parse(time.RFC3339, *key.EndValidity)
			if expDate.Before(time.Now()) {
				continue
			}
		}

		apiKeys = append(apiKeys, key)
	}

	return apiKeys, nil
}
