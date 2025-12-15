package auth

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"syscall"
	"time"

	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
)

type authRoundTripper struct {
	parent   http.RoundTripper
	auth     Auth
	attempts int
}

func (rt *authRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	token := rt.auth.GetAccessToken(req.Context())

	if token == "" {
		return nil, fmt.Errorf("não foi possível obter o token de acesso. Você esqueceu de fazer login?")
	}

	req.Header.Set("Authorization", "Bearer "+token)

	waitBeforeRetry := 100 * time.Millisecond
	var res *http.Response
	var err error

	if req.Body != nil {
		defer req.Body.Close()
	}

	for i := 0; i < rt.attempts; i++ {
		reqCopy := rt.cloneRequest(req)
		res, err = rt.parent.RoundTrip(reqCopy)

		if err != nil {
			var sysErr *os.SyscallError

			if os.IsTimeout(err) {
				fmt.Println("Request timeout, retrying...\n", "attempt", i+1, "\n ")
				time.Sleep(waitBeforeRetry)
				waitBeforeRetry = waitBeforeRetry * 2
				continue
			}

			if errors.As(err, &sysErr) {
				if sysErr.Err == syscall.ECONNRESET {
					fmt.Println("\n\n\nConn reset by peer! THIS IS A SERVER PROBLEM!!!\n\n\n", "attempt", i+1, "")
					time.Sleep(waitBeforeRetry)
					waitBeforeRetry = waitBeforeRetry * 2
					continue
				}
			}
			return res, err
		}
		if res.StatusCode >= 500 {
			fmt.Println("\n\n\nServer responded with fail, retrying...\n\n\n", "attempt", i+1, "status code", res.StatusCode, "")
			time.Sleep(waitBeforeRetry)
			waitBeforeRetry = waitBeforeRetry * 2
			continue
		}

		return res, err
	}

	return rt.parent.RoundTrip(req)
}

func (c *OAuthClient) AuthenticatedHttpClientFromContext(ctx context.Context) *http.Client {
	auth := ctx.Value(cmdutils.CTX_AUTH_KEY).(Auth)

	transport := c.httpClient.Transport

	if transport == nil {
		transport = http.DefaultTransport
	}

	transport = &authRoundTripper{parent: transport, auth: auth, attempts: 5}

	return &http.Client{
		Transport: transport,
		Timeout:   c.httpClient.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Don't follow redirects, return an error to use the last response.
			return http.ErrUseLastResponse
		}}
}

func (rt *authRoundTripper) cloneRequest(req *http.Request) *http.Request {
	var body io.Reader
	var err error

	if req.Body != nil {
		body, err = rt.cloneRequestBody(req)
		if err != nil {
			fmt.Println("Erro: %w", err)
			return req
		}
	}
	clonedRequest, err := http.NewRequestWithContext(req.Context(), req.Method, req.URL.String(), body)
	if err != nil {
		return req
	}
	clonedRequest.Header = req.Header
	return clonedRequest
}

func (rt *authRoundTripper) cloneRequestBody(req *http.Request) (io.Reader, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("erro: %w", err)
	}

	req.Body = io.NopCloser(bytes.NewBuffer(body))
	return bytes.NewReader(body), nil
}
