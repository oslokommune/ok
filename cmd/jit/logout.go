package jit

import (
	"fmt"

	"github.com/oslokommune/ok/pkg/jit"
	"github.com/spf13/cobra"
)

var LogoutCommand = &cobra.Command{
	Use:   "logout",
	Short: "Clear cached tokens.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := jit.ClearCache(); err != nil {
			return fmt.Errorf("clearing cache: %w", err)
		}

		fmt.Println("Logged out successfully.")
		return nil
	},
}
