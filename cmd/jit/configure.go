package jit

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/oslokommune/ok/pkg/jit"
	"github.com/spf13/cobra"
)

var ConfigureCommand = &cobra.Command{
	Use:   "configure",
	Short: "Configure JIT settings (tenant ID, client ID, base URL).",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := jit.LoadConfig()
		if cfg == nil {
			cfg = &jit.Config{}
		}

		err := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Tenant ID").
					Value(&cfg.TenantID),
				huh.NewInput().
					Title("Client ID").
					Value(&cfg.ClientID),
				huh.NewInput().
					Title("Base URL").
					Value(&cfg.BaseURL),
			),
		).Run()
		if err != nil {
			return fmt.Errorf("configuration cancelled: %w", err)
		}

		if cfg.TenantID == "" || cfg.ClientID == "" || cfg.BaseURL == "" {
			return fmt.Errorf("all fields are required")
		}

		if err := jit.SaveConfig(cfg); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		fmt.Println("Config saved.")
		return nil
	},
}
