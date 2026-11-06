package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/oslokommune/ok/pkg/doctor"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
	"github.com/spf13/cobra"
)

var doctorJSON bool

func init() {
	doctorCmd.Flags().BoolVar(&doctorJSON, "json", false, "Print the report as JSON")
	rootCmd.AddCommand(doctorCmd)
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Checks that your machine has the tools the Golden Path needs.",
	Long: `Checks that your machine has the tools the Golden Path needs, and tells you how to fix what is missing.

It checks the tools from https://km.oslo.systems/set-up-your-machine/ plus git, that Terraform satisfies
` + "`" + doctor.TerraformConstraint + "`" + ` (the version the Golden Path templates require), and whether a newer ok is out.
It never changes anything.

With --json, the report is printed to stdout as
` + "`" + `{"ok": bool, "checks": [{"name", "status", "found", "required", "fix"}]}` + "`" + `,
where status is pass, warn or fail.

Exits 1 if any check fails. Warnings never fail.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		latest := func() (string, error) {
			v, err := githubreleases.GetLatestOkVersion()
			if err != nil {
				return "", err
			}
			return v.String(), nil
		}
		report := doctor.New(VersionData.Version, latest).Run(cmd.Context())

		if doctorJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			if err := enc.Encode(report); err != nil {
				return err
			}
		} else {
			doctor.Print(os.Stdout, report)
		}

		if failed := report.Failed(); len(failed) > 0 {
			return fmt.Errorf("%d of %d checks failed", len(failed), len(report.Checks))
		}
		return nil
	},
}
