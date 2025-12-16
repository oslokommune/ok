package cmd

import (
	"github.com/oslokommune/ok/internal/scriptrunner"
	"github.com/spf13/cobra"
)

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

var forwardCommand = &cobra.Command{
	Use:   "forward",
	Short: "Starts a port forwarding session to a database.",
	Run: func(cmd *cobra.Command, args []string) {
		scriptArgs := getForwardCommandScriptFlags(cmd)
		scriptrunner.RunScript("port-forward.sh", scriptArgs)
	},
}
