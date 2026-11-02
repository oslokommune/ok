package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/oslokommune/ok/internal/scriptrunner"
	"github.com/oslokommune/ok/pkg/prompt"
	"github.com/spf13/cobra"
)

// forwardEnv are the values port-forward.sh needs to run without prompts. The script also reads
// them from the selections saved with --keep.
var forwardEnv = []string{"TASK_DEF_FAMILY", "RDS_ENDPOINT", "LOCAL_PORT", "REMOTE_PORT"}

func init() {
	rootCmd.AddCommand(forwardCommand)
	forwardCommand.Flags().BoolP("keep", "k", false, "Saves the selection you make for the next time")
	forwardCommand.Flags().BoolP("clear", "c", false, "Clears the saved selections")
}

// In order to forward the flags to the script, we need to collect them here and pass them as arguments to the script.
func getForwardCommandScriptFlags(cmd *cobra.Command) []string {
	scriptArgs := []string{}
	if ok, _ := cmd.Flags().GetBool("keep"); ok {
		scriptArgs = append(scriptArgs, "--keep")
	}
	if ok, _ := cmd.Flags().GetBool("clear"); ok {
		scriptArgs = append(scriptArgs, "--clear")
	}
	return scriptArgs
}

// forwardHasSelection reports whether port-forward.sh has every value it would otherwise ask for:
// from the environment, or from the selections saved in savedSelectionsFile.
func forwardHasSelection(getenv func(string) string, savedSelectionsFile string, clear bool) bool {
	if !clear {
		if _, err := os.Stat(savedSelectionsFile); err == nil {
			return true
		}
	}
	for _, name := range forwardEnv {
		if getenv(name) == "" {
			return false
		}
	}
	return true
}

var forwardCommand = &cobra.Command{
	Use:   "forward",
	Short: "Starts a port forwarding session to a database.",
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptArgs := getForwardCommandScriptFlags(cmd)

		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("finding home directory: %w", err)
		}
		clear, _ := cmd.Flags().GetBool("clear")
		if !forwardHasSelection(os.Getenv, filepath.Join(home, ".ok-port-forward.conf"), clear) {
			hint := fmt.Sprintf("Set %s to run without prompts.", strings.Join(forwardEnv, ", "))
			if err := prompt.Check("Select the database to forward to", hint); err != nil {
				return err
			}
		}

		return scriptrunner.RunScript("port-forward.sh", scriptArgs, prompt.Allowed())
	},
}
