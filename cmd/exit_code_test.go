package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/oslokommune/ok/pkg/prompt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExitCode(t *testing.T) {
	inputNeeded := &prompt.InputNeededError{Prompt: "Select cluster"}

	tests := []struct {
		name string
		err  error
		want int
	}{
		{name: "input needed", err: inputNeeded, want: 2},
		{name: "input needed, wrapped", err: fmt.Errorf("selecting packages: %w", inputNeeded), want: 2},
		{name: "any other error", err: errors.New("boom"), want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, exitCode(tt.err))
		})
	}
}

func TestLoadConfigurationKeepsStdoutClean(t *testing.T) {
	// Given a config file
	configFile := filepath.Join(t.TempDir(), "config.yml")
	require.NoError(t, os.WriteFile(configFile, []byte("enable_experimental: false\n"), 0o600))
	viper.SetConfigFile(configFile)
	t.Cleanup(viper.Reset)

	// When the configuration is loaded
	stdout, stderr := captureOutput(t, loadConfiguration)

	// Then the message goes to stderr, not stdout
	assert.Empty(t, stdout)
	assert.Contains(t, stderr, "Using config file: "+configFile)
}

func captureOutput(t *testing.T, f func()) (stdout string, stderr string) {
	t.Helper()

	outR, outW, err := os.Pipe()
	require.NoError(t, err)
	errR, errW, err := os.Pipe()
	require.NoError(t, err)

	oldOut, oldErr := os.Stdout, os.Stderr
	os.Stdout, os.Stderr = outW, errW
	defer func() { os.Stdout, os.Stderr = oldOut, oldErr }()

	f()

	require.NoError(t, outW.Close())
	require.NoError(t, errW.Close())
	out, err := io.ReadAll(outR)
	require.NoError(t, err)
	errOut, err := io.ReadAll(errR)
	require.NoError(t, err)

	return string(out), string(errOut)
}
