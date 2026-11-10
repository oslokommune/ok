package cmd

import (
	"errors"
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

// forwardSavedSelections is where port-forward.sh saves selections with --keep, under the home directory.
const forwardSavedSelections = "~/.ok-port-forward.conf"

func init() {
	rootCmd.AddCommand(forwardCommand)
	forwardCommand.Flags().BoolP("keep", "k", false, "Saves the database endpoint and ports in "+forwardSavedSelections+" for next time")
	forwardCommand.Flags().BoolP("clear", "c", false, "Deletes the saved selections in "+forwardSavedSelections)
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

// forwardHasSelection reports whether port-forward.sh has every value it would otherwise ask for,
// each from the environment or from the selections saved in savedSelectionsFile.
func forwardHasSelection(getenv func(string) string, savedSelectionsFile string, clear bool) bool {
	saved := map[string]string{}
	if !clear {
		saved = readSavedSelections(savedSelectionsFile)
	}
	for _, name := range forwardEnv {
		if getenv(name) == "" && saved[name] == "" {
			return false
		}
	}
	return true
}

// readSavedSelections reads the NAME=value lines that port-forward.sh writes with --keep.
func readSavedSelections(path string) map[string]string {
	selections := map[string]string{}
	content, err := os.ReadFile(path)
	if err != nil {
		return selections
	}
	for _, line := range strings.Split(string(content), "\n") {
		if name, value, ok := strings.Cut(strings.TrimSpace(line), "="); ok {
			selections[name] = value
		}
	}
	return selections
}

var forwardCommand = &cobra.Command{
	Use:   "forward",
	Short: "Starts a port forwarding session to a database.",
	Long: `Starts a port forwarding session to a database, through a port forwarding task in ECS.

ok forward asks for the task definition, the database and the ports. To run it without
prompts, set all four and pass --non-interactive:

  TASK_DEF_FAMILY  the port forwarding task's definition family (it starts with ssm-portfwd)
  RDS_ENDPOINT     the database's endpoint address
  LOCAL_PORT       the port on this machine, e.g. 4812
  REMOTE_PORT      the database's port, e.g. 5432

--keep saves the database endpoint and the ports in ` + forwardSavedSelections + `, and ok forward
uses them next time. It does not save the task definition: set TASK_DEF_FAMILY for that.
--clear deletes the file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptArgs := getForwardCommandScriptFlags(cmd)

		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("finding home directory: %w", err)
		}
		savedSelectionsFile := filepath.Join(home, strings.TrimPrefix(forwardSavedSelections, "~/"))
		clear, _ := cmd.Flags().GetBool("clear")
		if clear {
			if err := os.Remove(savedSelectionsFile); err != nil && !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("deleting saved selections: %w", err)
			}
		}
		if !forwardHasSelection(os.Getenv, savedSelectionsFile, clear) {
			hint := fmt.Sprintf("Set %s to run without prompts. See ok forward --help.", strings.Join(forwardEnv, ", "))
			if err := prompt.Check("Select the database to forward to", hint); err != nil {
				return err
			}
		}

		return scriptrunner.RunScript("port-forward.sh", scriptArgs, prompt.Allowed())
	},
}
