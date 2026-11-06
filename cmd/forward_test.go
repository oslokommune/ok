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
	all := map[string]string{"TASK_DEF_FAMILY": "ssm-portfwd-app", "RDS_ENDPOINT": "db.example", "LOCAL_PORT": "4812", "REMOTE_PORT": "5432"}
	some := map[string]string{"TASK_DEF_FAMILY": "ssm-portfwd-app", "RDS_ENDPOINT": "db.example"}
	family := map[string]string{"TASK_DEF_FAMILY": "ssm-portfwd-app"}

	// What --keep writes: the task definition's ARN, not its family
	saved := filepath.Join(t.TempDir(), ".ok-port-forward.conf")
	require.NoError(t, os.WriteFile(saved, []byte(
		"LOCAL_PORT=4812\nTASK_DEF_ARN=arn:aws:ecs:eu-west-1:123:task-definition/ssm-portfwd-app:3\nREMOTE_PORT=5432\nRDS_ENDPOINT=db.example\n"), 0o600))
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
		{name: "saved selections lack the task definition family", saved: saved, want: false},
		{name: "saved selections and TASK_DEF_FAMILY", env: family, saved: saved, want: true},
		{name: "saved selections cleared with --clear", env: family, saved: saved, clear: true, want: false},
		{name: "every value in the environment, with --clear", env: all, saved: saved, clear: true, want: true},
		{name: "nothing", saved: missing, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getenv := func(name string) string { return tt.env[name] }

			assert.Equal(t, tt.want, forwardHasSelection(getenv, tt.saved, tt.clear))
		})
	}
}

func TestForwardClearDeletesSavedSelections(t *testing.T) {
	// Given saved selections, and ok may not prompt so the script does not run
	home := t.TempDir()
	t.Setenv("HOME", home)
	saved := filepath.Join(home, ".ok-port-forward.conf")
	require.NoError(t, os.WriteFile(saved, []byte("RDS_ENDPOINT=db.example\n"), 0o600))
	prompt.SetMode(prompt.Mode{Reason: "CI=true is set"})
	t.Cleanup(func() { prompt.SetMode(prompt.Mode{Interactive: true}) })
	for _, name := range forwardEnv {
		t.Setenv(name, "")
	}
	require.NoError(t, forwardCommand.Flags().Set("clear", "true"))
	t.Cleanup(func() { _ = forwardCommand.Flags().Set("clear", "false") })

	// When
	_ = forwardCommand.RunE(forwardCommand, nil)

	// Then
	assert.NoFileExists(t, saved)
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
