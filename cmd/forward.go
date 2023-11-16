package cmd

import (
	"github.com/oslokommune/ok/scriptrunner"
	"github.com/spf13/cobra"
)

func newForwardCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "forward",
		Short: "Starts a port forwarding session to a database.",
		Run: func(cmd *cobra.Command, args []string) {
			scriptrunner.RunScript("port-forward.sh", args)
		},
	}
}
