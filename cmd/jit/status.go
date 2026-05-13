package jit

import (
	"fmt"
	"time"

	"github.com/oslokommune/ok/pkg/jit"
	"github.com/spf13/cobra"
)

var StatusCommand = &cobra.Command{
	Use:   "status",
	Short: "Show active access grants.",
	RunE: func(cmd *cobra.Command, args []string) error {
		grants, err := jit.LoadGrants()
		if err != nil {
			return fmt.Errorf("loading grants: %w", err)
		}

		if len(grants) == 0 {
			fmt.Println("No active access grants.")
			return nil
		}

		now := time.Now()
		fmt.Printf("Active grants (%d):\n\n", len(grants))
		for _, g := range grants {
			remaining := g.ExpiresAt.Sub(now).Round(time.Minute)
			fmt.Printf("  %s  (expires in %s)\n", g.Group, remaining)
		}

		return nil
	},
}
