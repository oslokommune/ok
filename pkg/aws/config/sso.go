package config

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
)

const (
	timeout    = 3 * time.Minute
	clientName = "ok-aws-config-generator"
	pkceLength = 64
)

//go:embed assets/callback_result.html
var resultHTML string

// callbackServer manages the local HTTP server used during the OAuth2 authentication flow.
type callbackServer struct {
	server *http.Server
	code   chan string
	state  string
}

// GetAccounts retrieves all accessible AWS accounts and their roles through SSO.
// It handles the browser-based authentication flow and token management.
func GetAccounts(ctx context.Context, ssoClient *sso.Client, ssooidcClient *ssooidc.Client, startUrl, region string) ([]Account, error) {
	server, err := startCallbackServer()
	if err != nil {
		return nil, fmt.Errorf("starting callback server: %w", err)
	}
	defer server.close()

	token, err := authorize(ctx, ssooidcClient, server, startUrl, region)
	if err != nil {
		return nil, err
	}

	return listAccounts(ctx, ssoClient, *token.AccessToken)
}

func startCallbackServer() (*callbackServer, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("creating listener: %w", err)
	}

	state, err := generateState()
	if err != nil {
		listener.Close()
		return nil, fmt.Errorf("generating state: %w", err)
	}

	code := make(chan string, 1)
	server := &callbackServer{
		code:  code,
		state: state,
	}

	mux := http.NewServeMux()
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		tmpl := template.Must(template.New("result").Parse(resultHTML))
		data := struct {
			Class   string
			Icon    string
			Message string
		}{}

		if gotState := r.URL.Query().Get("state"); gotState != state {
			w.WriteHeader(http.StatusBadRequest)
			data.Class = "error"
			data.Icon = "✕"
			data.Message = "Authorization error"
			if err := tmpl.Execute(w, data); err != nil {
			}
			w.(http.Flusher).Flush()
			code <- ""
			server.close()
			return
		}

		authCode := r.URL.Query().Get("code")
		if authCode == "" {
			w.WriteHeader(http.StatusBadRequest)
			data.Class = "error"
			data.Icon = "✕"
			data.Message = "Authorization error"
			if err := tmpl.Execute(w, data); err != nil {
			}
			w.(http.Flusher).Flush()
			code <- ""
			server.close()
			return
		}

		data.Class = "success"
		data.Icon = "✓"
		data.Message = "Authorization success"
		if err := tmpl.Execute(w, data); err != nil {
		}
		w.(http.Flusher).Flush()
		code <- authCode
		server.close()
	}

	mux.HandleFunc("/", handler)

	server.server = &http.Server{
		Handler: mux,
		Addr:    listener.Addr().String(),
	}

	fmt.Fprintf(os.Stderr, "Started local callback server at: http://%s\n", server.server.Addr)

	go server.server.Serve(listener)

	return server, nil
}

func (s *callbackServer) waitForCode(ctx context.Context) (string, error) {
	select {
	case code := <-s.code:
		if code == "" {
			return "", fmt.Errorf("authorization failed or was denied")
		}
		return code, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func (s *callbackServer) close() error {
	if err := s.server.Close(); err != nil {
		return fmt.Errorf("closing server: %w", err)
	}
	return nil
}

func (s *callbackServer) getCallbackURL() string {
	return fmt.Sprintf("http://%s", s.server.Addr)
}

// generateState creates a random 16-byte state parameter for OAuth2 CSRF protection.
// The state is hex-encoded to make it URL-safe.
// Returns the hex-encoded state string and any error encountered.
func generateState() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// generateCodeVerifier creates a random code verifier string for PKCE (Proof Key for Code Exchange).
// The verifier is a random string of length pkceLength using characters from the allowed charset.
// The charset follows RFC 7636 requirements: characters from A-Z, a-z, 0-9, and "-._~".
// Returns the code verifier string and any error encountered.
func generateCodeVerifier() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~"

	b := make([]byte, pkceLength)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generating random bytes: %w", err)
	}

	result := make([]byte, pkceLength)
	for i := range b {
		result[i] = charset[int(b[i])%len(charset)]
	}
	return string(result), nil
}

// generateCodeChallenge creates a code challenge from a code verifier for PKCE.
// The challenge is created by:
// 1. Taking the SHA256 hash of the verifier
// 2. Base64URL-encoding the hash (without padding)
// This follows the S256 transformation method specified in RFC 7636.
func generateCodeChallenge(verifier string) (string, error) {
	if len(verifier) < 43 || len(verifier) > 128 {
		return "", fmt.Errorf("code verifier length must be between 43 and 128 characters")
	}
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:]), nil
}

func openURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

// authorize performs the OAuth2 authorization flow with AWS SSO.
// It opens the user's browser for authentication and waits for the callback.
func authorize(ctx context.Context, ssooidcClient *ssooidc.Client, server *callbackServer, startUrl, region string) (*ssooidc.CreateTokenOutput, error) {
	clientCreds, err := registerClient(ctx, ssooidcClient, server.getCallbackURL(), startUrl)
	if err != nil {
		return nil, err
	}

	verifier, err := generateCodeVerifier()
	if err != nil {
		return nil, fmt.Errorf("generating code verifier: %w", err)
	}
	challenge, err := generateCodeChallenge(verifier)
	if err != nil {
		return nil, fmt.Errorf("generating code challenge: %w", err)
	}

	authURL := fmt.Sprintf(
		"https://oidc.%s.amazonaws.com/authorize?%s",
		region,
		url.Values{
			"client_id":             {*clientCreds.ClientId},
			"response_type":         {"code"},
			"redirect_uri":          {server.getCallbackURL()},
			"code_challenge":        {challenge},
			"code_challenge_method": {"S256"},
			"scope":                 {"sso:account:access"},
			"state":                 {server.state},
		}.Encode(),
	)

	fmt.Fprintf(os.Stderr, "\nOpening browser at URL: %s\n", authURL)
	if err := openURL(authURL); err != nil {
		fmt.Fprintf(os.Stderr, "\nFailed to open browser. Please open the URL manually.\n")
	}
	fmt.Fprintf(os.Stderr, "\nWaiting for browser authorization\n")

	authCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	code, err := server.waitForCode(authCtx)
	if err != nil {
		return nil, fmt.Errorf("waiting for authorization code: %w", err)
	}

	return createToken(ctx, ssooidcClient, clientCreds, code, server.getCallbackURL(), verifier)
}

func registerClient(ctx context.Context, client *ssooidc.Client, callbackURL, startUrl string) (*ssooidc.RegisterClientOutput, error) {
	return client.RegisterClient(ctx, &ssooidc.RegisterClientInput{
		ClientName:   aws.String(clientName),
		ClientType:   aws.String("public"),
		Scopes:       []string{"sso:account:access"},
		GrantTypes:   []string{"authorization_code", "refresh_token"},
		IssuerUrl:    aws.String(startUrl),
		RedirectUris: []string{callbackURL},
	})
}

func createToken(ctx context.Context, client *ssooidc.Client, clientCreds *ssooidc.RegisterClientOutput, code, callbackURL, verifier string) (*ssooidc.CreateTokenOutput, error) {
	return client.CreateToken(ctx, &ssooidc.CreateTokenInput{
		ClientId:     clientCreds.ClientId,
		ClientSecret: clientCreds.ClientSecret,
		GrantType:    aws.String("authorization_code"),
		Code:         aws.String(code),
		RedirectUri:  aws.String(callbackURL),
		CodeVerifier: aws.String(verifier),
	})
}

// listAccounts retrieves all AWS accounts accessible to the authenticated user.
// For each account, it also fetches the available IAM roles.
func listAccounts(ctx context.Context, ssoClient *sso.Client, accessToken string) ([]Account, error) {
	var accounts []Account
	paginator := sso.NewListAccountsPaginator(ssoClient, &sso.ListAccountsInput{
		AccessToken: &accessToken,
	})

	fmt.Fprintf(os.Stderr, "\nFetching available accounts and roles")
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("listing accounts: %w", err)
		}

		for _, acc := range page.AccountList {
			account := Account{
				ID:   *acc.AccountId,
				Name: *acc.AccountName,
			}

			roles, err := listAccountRoles(ctx, ssoClient, accessToken, *acc.AccountId)
			if err != nil {
				return nil, fmt.Errorf("listing roles for account %s: %w", *acc.AccountId, err)
			}
			account.Roles = roles

			accounts = append(accounts, account)
			fmt.Fprintf(os.Stderr, ".")
		}
	}
	fmt.Fprintln(os.Stderr)

	return accounts, nil
}

func listAccountRoles(ctx context.Context, ssoClient *sso.Client, accessToken string, accountID string) ([]string, error) {
	var roles []string
	paginator := sso.NewListAccountRolesPaginator(ssoClient, &sso.ListAccountRolesInput{
		AccessToken: &accessToken,
		AccountId:   &accountID,
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("listing roles: %w", err)
		}

		for _, role := range page.RoleList {
			roles = append(roles, *role.RoleName)
		}
	}

	return roles, nil
}
