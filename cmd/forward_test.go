package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/oslokommune/ok/pkg/prompt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestForwardHasSelection(t *testing.T) {
	all := map[string]string{"TASK_DEF_FAMILY": "app", "RDS_ENDPOINT": "db.example", "LOCAL_PORT": "4812", "REMOTE_PORT": "5432"}
	some := map[string]string{"TASK_DEF_FAMILY": "app", "RDS_ENDPOINT": "db.example"}

	saved := filepath.Join(t.TempDir(), ".ok-port-forward.conf")
	require.NoError(t, os.WriteFile(saved, []byte("TASK_DEF_FAMILY=app\n"), 0o600))
	missing := filepath.Join(t.TempDir(), ".ok-port-forward.conf")

	tests := []struct {
		name  string
		env   map[string]string
		saved string
		clear bool
		want  bool
	}{
		{name: "every value in the environment", env: all, saved: missing, want: true},
		{name: "some values in the environment", env: some, saved: missing, want: false},
		{name: "saved selections", saved: saved, want: true},
		{name: "saved selections, cleared with --clear", saved: saved, clear: true, want: false},
		{name: "nothing", saved: missing, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getenv := func(name string) string { return tt.env[name] }

			assert.Equal(t, tt.want, forwardHasSelection(getenv, tt.saved, tt.clear))
		})
	}
}

func TestForwardWithoutPrompting(t *testing.T) {
	// Given ok may not prompt, and nothing tells the script what to forward
	prompt.SetMode(prompt.Mode{Reason: "CI=true is set"})
	t.Cleanup(func() { prompt.SetMode(prompt.Mode{Interactive: true}) })
	t.Setenv("HOME", t.TempDir())
	for _, name := range forwardEnv {
		t.Setenv(name, "")
	}

	// When
	err := forwardCommand.RunE(forwardCommand, nil)

	// Then ok stops before the script runs, and says what to set
	var inputNeeded *prompt.InputNeededError
	require.ErrorAs(t, err, &inputNeeded)
	assert.Contains(t, inputNeeded.Hint, "TASK_DEF_FAMILY, RDS_ENDPOINT, LOCAL_PORT, REMOTE_PORT")
}
