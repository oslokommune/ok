package jit

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
)

// Client wraps an MSAL public client for device code authentication.
type Client struct {
	app      public.Client
	clientID string
}

// NewClient creates a new MSAL public client with persistent token cache.
func NewClient(clientID, tenantID string) (*Client, error) {
	cache, err := newTokenCache()
	if err != nil {
		return nil, fmt.Errorf("creating token cache: %w", err)
	}

	authority := fmt.Sprintf("https://login.microsoftonline.com/%s", tenantID)
	app, err := public.New(clientID, public.WithAuthority(authority), public.WithCache(cache))
	if err != nil {
		return nil, fmt.Errorf("creating MSAL client: %w", err)
	}

	return &Client{app: app, clientID: clientID}, nil
}

// Login acquires a token, using the cache first, falling back to device code flow.
func (c *Client) Login(ctx context.Context) (public.AuthResult, error) {
	return c.login(ctx, false)
}

// LoginInteractive acquires a token, using the cache first, falling back to the
// interactive browser flow with a localhost redirect.
func (c *Client) LoginInteractive(ctx context.Context) (public.AuthResult, error) {
	return c.login(ctx, true)
}

func (c *Client) login(ctx context.Context, interactive bool) (public.AuthResult, error) {
	scopes := []string{fmt.Sprintf("%s/.default", c.clientID)}

	// Try silent acquisition first (cached/refresh token)
	accounts, err := c.app.Accounts(ctx)
	if err == nil && len(accounts) > 0 {
		result, err := c.app.AcquireTokenSilent(ctx, scopes, public.WithSilentAccount(accounts[0]))
		if err == nil {
			return result, nil
		}
	}

	if interactive {
		return c.app.AcquireTokenInteractive(ctx, scopes, public.WithRedirectURI("http://localhost"))
	}

	// Fall back to device code flow
	dc, err := c.app.AcquireTokenByDeviceCode(ctx, scopes)
	if err != nil {
		return public.AuthResult{}, fmt.Errorf("acquiring device code: %w", err)
	}

	copyToClipboard(dc.Result.UserCode)
	fmt.Printf("\nCode %s copied to clipboard.\n", dc.Result.UserCode)
	fmt.Printf("Opening browser at: %s\n\n", dc.Result.VerificationURL)
	openURL(dc.Result.VerificationURL)

	return dc.AuthenticationResult(ctx)
}

func copyToClipboard(text string) {
	var cmd string
	switch runtime.GOOS {
	case "windows":
		cmd = "clip"
	case "darwin":
		cmd = "pbcopy"
	default:
		cmd = "xclip"
	}
	c := exec.Command(cmd)
	c.Stdin = strings.NewReader(text)
	c.Run()
}

func openURL(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	exec.Command(cmd, args...).Start()
}
