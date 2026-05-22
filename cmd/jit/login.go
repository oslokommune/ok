package jit

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	"github.com/charmbracelet/huh"
	"github.com/oslokommune/ok/pkg/jit"
	"github.com/spf13/cobra"
)

const accessRequestPath = "/prod/access-request"

func init() {
	LoginCommand.Flags().Bool("debug", false, "Print debug info about the token")
	LoginCommand.Flags().Bool("interactive", false, "Use interactive browser login instead of device code flow")
	LoginCommand.Flags().Float64("hours", 0, "Hours of access needed (1-4). Defaults to server default when unset.")
}

func loadConfig() (*jit.Config, error) {
	cfg, err := jit.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	if cfg == nil || cfg.TenantID == "" || cfg.ClientID == "" || cfg.BaseURL == "" {
		return nil, fmt.Errorf("JIT is not configured. Run 'ok jit configure' first")
	}

	return cfg, nil
}

var LoginCommand = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Azure Entra ID and select a group.",
	RunE: func(cmd *cobra.Command, args []string) error {
		debug, _ := cmd.Flags().GetBool("debug")
		interactive, _ := cmd.Flags().GetBool("interactive")
		hoursFlag, _ := cmd.Flags().GetFloat64("hours")

		var hours *float64
		if hoursFlag != 0 {
			if hoursFlag < 1 || hoursFlag > 4 {
				return fmt.Errorf("--hours must be between 1 and 4")
			}
			hours = &hoursFlag
		}

		cfg, err := loadConfig()
		if err != nil {
			return err
		}

		client, err := jit.NewClient(cfg.ClientID, cfg.TenantID)
		if err != nil {
			return fmt.Errorf("initializing auth client: %w", err)
		}

		ctx := context.Background()

		var result public.AuthResult
		if interactive {
			result, err = client.LoginInteractive(ctx)
		} else {
			result, err = client.Login(ctx)
		}
		if err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		if result.IDToken.Name != "" {
			fmt.Printf("Logged in as: %s (%s)\n\n", result.IDToken.Name, result.IDToken.PreferredUsername)
		}

		if debug {
			printDebugInfo(result)
		}

		groups := jit.ExtractGroupsFromJWT(result.AccessToken)
		if len(groups) == 0 {
			fmt.Println("No groups found in token.")
			return nil
		}

		options := make([]huh.Option[string], len(groups))
		for i, m := range groups {
			options[i] = huh.NewOption(m, m)
		}

		var selected string
		err = huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Options(options...).
					Title("Select group to escalate to production").
					Value(&selected),
			),
		).Run()
		if err != nil {
			if errors.Is(err, huh.ErrUserAborted) {
				return nil
			}
			return fmt.Errorf("selection failed: %w", err)
		}

		apiURL := cfg.BaseURL + accessRequestPath

		respBody, err := apiRequest(http.MethodPost, apiURL, result.AccessToken, selected, hours)
		if err != nil {
			return err
		}

		var apiResp struct {
			AccessExpiresAt time.Time `json:"access_expires_at"`
		}
		if json.Unmarshal(respBody, &apiResp) == nil && !apiResp.AccessExpiresAt.IsZero() {
			jit.SaveGrant(jit.Grant{
				Group:     selected,
				ExpiresAt: apiResp.AccessExpiresAt,
			})
			fmt.Printf("\nAccess granted for %s (expires %s).\n", selected,
				apiResp.AccessExpiresAt.Local().Format("15:04"))
		} else {
			fmt.Printf("\nAccess requested for %s.\n", selected)
		}

		return nil
	},
}

func printDebugInfo(result public.AuthResult) {
	now := time.Now()

	fmt.Println("--- Debug info ---")
	fmt.Printf("Account:          %s\n", result.Account.PreferredUsername)
	fmt.Printf("Tenant ID:        %s\n", result.IDToken.TenantID)
	fmt.Printf("Object ID:        %s\n", result.IDToken.Oid)
	fmt.Printf("Issuer:           %s\n", result.IDToken.Issuer)
	fmt.Printf("Audience:         %s\n", result.IDToken.Audience)
	fmt.Printf("Granted scopes:   %s\n", strings.Join(result.GrantedScopes, ", "))
	fmt.Printf("Token source:     %s\n", tokenSourceString(result.Metadata.TokenSource))
	fmt.Printf("Refresh token:    %v\n", jit.GetRefreshTokenInfo().Present)

	fmt.Printf("Access token expires: %s (in %s)\n",
		result.ExpiresOn.Format(time.RFC3339),
		result.ExpiresOn.Sub(now).Round(time.Second))

	if !result.Metadata.RefreshOn.IsZero() {
		fmt.Printf("Refresh recommended:  %s (in %s)\n",
			result.Metadata.RefreshOn.Format(time.RFC3339),
			result.Metadata.RefreshOn.Sub(now).Round(time.Second))
	}

	if result.IDToken.ExpirationTime > 0 {
		idExp := time.Unix(result.IDToken.ExpirationTime, 0)
		fmt.Printf("ID token expires:     %s (in %s)\n",
			idExp.Format(time.RFC3339),
			idExp.Sub(now).Round(time.Second))
	}

	fmt.Println("------------------")
	fmt.Println()
}

func tokenSourceString(s public.TokenSource) string {
	switch s {
	case public.TokenSourceCache:
		return "cache"
	case public.TokenSourceIdentityProvider:
		return "identity provider (fresh login)"
	default:
		return fmt.Sprintf("unknown (%d)", s)
	}
}

func apiRequest(method, apiURL, accessToken, group string, hours *float64) ([]byte, error) {
	body := map[string]interface{}{
		"group": group,
	}
	if hours != nil {
		body["hours"] = *hours
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshalling payload: %w", err)
	}

	req, err := http.NewRequest(method, apiURL, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not reach the JIT API at %s", apiURL)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr struct {
			Error string `json:"error"`
		}
		if json.Unmarshal(respBody, &apiErr) == nil && apiErr.Error != "" {
			return nil, fmt.Errorf("%s", apiErr.Error)
		}

		return nil, fmt.Errorf("JIT API returned an error (HTTP %d)", resp.StatusCode)
	}

	return respBody, nil
}
