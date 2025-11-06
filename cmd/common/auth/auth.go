package auth

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/magaluCloud/mgccli/cmd/common/structs"
	"github.com/magaluCloud/mgccli/cmd/common/workspace"
	"gopkg.in/yaml.v3"
)

// access_key_id: ""
// access_token: ""
// current_environment: ""
// refresh_token: ""
// secret_access_key: ""

type AuthFile struct {
	AccessKeyID     string `yaml:"access_key_id"`
	AccessToken     string `yaml:"access_token"`
	RefreshToken    string `yaml:"refresh_token"`
	SecretAccessKey string `yaml:"secret_access_key"`
}

type Auth interface {
	GetAccessKeyID() string
	GetAccessToken(ctx context.Context) string
	GetRefreshToken() string
	GetSecretAccessKey() string
	GetService() *Service
	TokenClaims() (*TokenClaims, error)

	SetAccessToken(token string) error
	SetRefreshToken(token string) error
	SetSecretAccessKey(key string) error
	SetAccessKeyID(key string) error

	ValidateToken() error
	RefreshToken(ctx context.Context) error
}

type authValue struct {
	authValue AuthFile
	workspace workspace.Workspace
	service   *Service
}

func NewAuth(workspace workspace.Workspace) Auth {
	authFile := path.Join(workspace.Dir(), "auth.yaml")
	authContent, err := structs.LoadFileToStruct[AuthFile](authFile)

	config := DefaultConfig()
	service := NewService(config)

	if err != nil {
		//TODO: Handle error
		panic(err)
	}
	return &authValue{workspace: workspace, authValue: authContent, service: service}
}

func (a *authValue) GetService() *Service {
	return a.service
}

func (a *authValue) GetAccessKeyID() string {
	return a.authValue.AccessKeyID
}

func (a *authValue) GetAccessToken(ctx context.Context) string {
	if a.authValue.AccessToken == "" {
		return ""
	}
	err := a.ValidateToken()
	if err != nil {
		if a.authValue.RefreshToken != "" {
			err := a.RefreshToken(ctx)
			if err != nil {
				return ""
			}
			return a.authValue.AccessToken
		}
	}
	return a.authValue.AccessToken
}

func (a *authValue) GetRefreshToken() string {
	return a.authValue.RefreshToken
}

func (a *authValue) GetSecretAccessKey() string {
	return a.authValue.SecretAccessKey
}

func (a *authValue) SetAccessToken(token string) error {
	a.authValue.AccessToken = token
	return a.Write()
}

func (a *authValue) SetRefreshToken(token string) error {
	a.authValue.RefreshToken = token
	return a.Write()
}

func (a *authValue) SetSecretAccessKey(key string) error {
	a.authValue.SecretAccessKey = key
	return a.Write()
}

func (a *authValue) SetAccessKeyID(key string) error {
	a.authValue.AccessKeyID = key
	return a.Write()
}

func (a *authValue) Logout(name string) error {
	a.SetAccessToken("")
	a.SetRefreshToken("")
	a.SetSecretAccessKey("")
	a.SetAccessKeyID("")
	return a.Write()
}

func (a *authValue) Write() error {
	data, err := yaml.Marshal(a.authValue)
	if err != nil {
		return err
	}
	err = os.WriteFile(path.Join(a.workspace.Dir(), "auth.yaml"), data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (a *authValue) TokenClaims() (*TokenClaims, error) {
	if a.authValue.AccessToken == "" {
		return nil, fmt.Errorf("access token is not set")
	}

	tokenClaims := &TokenClaims{}
	tokenParser := jwt.NewParser()

	_, _, err := tokenParser.ParseUnverified(a.authValue.AccessToken, tokenClaims)
	if err != nil {
		return nil, err
	}

	return tokenClaims, nil
}

func (a *authValue) ValidateToken() error {
	//extract iat from token, if expires in less than 30 sec, run refresh operation
	tokenClaims, err := a.TokenClaims()
	if err != nil {
		return err
	}
	iat := tokenClaims.ExpiresAt.Time.Unix()
	if iat < time.Now().Unix()-60 {
		return fmt.Errorf("token expired")
	}
	return nil
}

func (a *authValue) RefreshToken(ctx context.Context) error {
	token, err := a.service.RefreshToken(ctx, a.authValue.RefreshToken)
	if err != nil {
		return err
	}
	a.authValue.AccessToken = token.AccessToken
	a.authValue.RefreshToken = token.RefreshToken
	return a.Write()
}
