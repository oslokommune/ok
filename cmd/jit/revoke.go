package jit

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/oslokommune/ok/pkg/jit"
	"github.com/spf13/cobra"
)

var RevokeCommand = &cobra.Command{
	Use:   "revoke",
	Short: "Revoke access to a group.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}

		grants, err := jit.LoadGrants()
		if err != nil {
			return fmt.Errorf("loading grants: %w", err)
		}

		if len(grants) == 0 {
			fmt.Println("No active access grants.")
			return nil
		}

		now := time.Now()
		options := make([]huh.Option[string], len(grants))
		for i, g := range grants {
			remaining := g.ExpiresAt.Sub(now).Round(time.Minute)
			label := fmt.Sprintf("%s (expires in %s)", g.Group, remaining)
			options[i] = huh.NewOption(label, g.Group)
		}

		var selected string
		err = huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Options(options...).
					Title("Select grant to revoke").
					Value(&selected),
			),
		).Run()
		if err != nil {
			if errors.Is(err, huh.ErrUserAborted) {
				return nil
			}
			return fmt.Errorf("selection failed: %w", err)
		}

		client, err := jit.NewClient(cfg.ClientID, cfg.TenantID)
		if err != nil {
			return fmt.Errorf("initializing auth client: %w", err)
		}

		ctx := context.Background()

		result, err := client.Login(ctx)
		if err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		apiURL := cfg.BaseURL + accessRequestPath

		if _, err := apiRequest(http.MethodDelete, apiURL, result.AccessToken, selected, nil); err != nil {
			return err
		}

		jit.RemoveGrant(selected)

		fmt.Printf("\nAccess revoked for %s.\n", selected)

		return nil
	},
}
