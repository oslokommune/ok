package cmd

import (
	"fmt"
	"testing"
)

func TestRewordCompletionCommand(t *testing.T) {
	// If cobra renames the completion command or a shell subcommand, the help text
	// falls back to cobra's default without an error. This test catches that.
	for name, long := range completionLongs {
		args := []string{"completion", name}
		if name == "completion" {
			args = args[:1]
		}

		cmd, _, err := rootCmd.Find(args)
		if err != nil {
			t.Fatalf("find %v: %v", args, err)
		}

		want := fmt.Sprintf(long, rootCmd.Name())
		if cmd.Long != want {
			t.Errorf("%v has help text %q, want %q", args, cmd.Long, want)
		}
	}
}
