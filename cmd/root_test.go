package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/oslokommune/ok/pkg/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIterate(t *testing.T) {
	t.Run("test errors", func(t *testing.T) {

		a := fmt.Errorf("deepest error a")
		b := fmt.Errorf("middle b that wraps a: %w", a)
		c := fmt.Errorf("middle c that wraps b: %w", b)
		d := fmt.Errorf("top level d that wraps c: %w", c)

		fmt.Println("--------")
		prettyPrintError(d)
		fmt.Println("--------")
	})
}

func TestExitCode(t *testing.T) {
	inputNeeded := &prompt.InputNeededError{Prompt: "Select cluster", Reason: "CI=true is set"}

	// A program ok runs that exits 2 must not make ok exit 2: 2 means ok needs input.
	var programFailed *exec.ExitError
	require.ErrorAs(t, exec.Command("sh", "-c", "exit 2").Run(), &programFailed)
	require.Equal(t, 2, programFailed.ExitCode())

	tests := []struct {
		name string
		err  error
		want int
	}{
		{name: "input needed", err: inputNeeded, want: 2},
		{name: "input needed, wrapped", err: fmt.Errorf("selecting packages: %w", inputNeeded), want: 2},
		{name: "input needed, joined", err: errors.Join(errors.New("other"), inputNeeded), want: 2},
		{name: "any other error", err: errors.New("boom"), want: 1},
		{name: "a program that exits 2, wrapped", err: fmt.Errorf("running boilerplate command: %w", programFailed), want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, exitCode(tt.err))
		})
	}
}

func TestPromptMode(t *testing.T) {
	on, off := true, false

	tests := []struct {
		name        string
		flag        *bool
		env         map[string]string
		hasTerminal bool
		want        prompt.Mode
	}{
		{name: "a terminal and nothing set", hasTerminal: true,
			want: prompt.Mode{Interactive: true}},
		{name: "no terminal", hasTerminal: false,
			want: prompt.Mode{Reason: "there is no terminal"}},
		{name: "--non-interactive", flag: &on, hasTerminal: true,
			want: prompt.Mode{Reason: "--non-interactive is set"}},
		{name: "--non-interactive=false wins over CI and no terminal", flag: &off, env: map[string]string{"CI": "true"}, hasTerminal: false,
			want: prompt.Mode{Interactive: true}},
		{name: "OK_NON_INTERACTIVE=true", env: map[string]string{"OK_NON_INTERACTIVE": "true"}, hasTerminal: true,
			want: prompt.Mode{Reason: "OK_NON_INTERACTIVE=true is set"}},
		{name: "OK_NON_INTERACTIVE=1", env: map[string]string{"OK_NON_INTERACTIVE": "1"}, hasTerminal: true,
			want: prompt.Mode{Reason: "OK_NON_INTERACTIVE=1 is set"}},
		{name: "OK_NON_INTERACTIVE=false wins over CI", env: map[string]string{"OK_NON_INTERACTIVE": "false", "CI": "true"}, hasTerminal: true,
			want: prompt.Mode{Interactive: true}},
		{name: "CI=true", env: map[string]string{"CI": "true"}, hasTerminal: true,
			want: prompt.Mode{Reason: "CI=true is set"}},
		{name: "CI=1", env: map[string]string{"CI": "1"}, hasTerminal: true,
			want: prompt.Mode{Reason: "CI=1 is set"}},
		{name: "CI=false falls back to the terminal", env: map[string]string{"CI": "false"}, hasTerminal: false,
			want: prompt.Mode{Reason: "there is no terminal"}},
		{name: "a value that is not a boolean is ignored", env: map[string]string{"OK_NON_INTERACTIVE": "yes"}, hasTerminal: true,
			want: prompt.Mode{Interactive: true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getenv := func(name string) string { return tt.env[name] }

			assert.Equal(t, tt.want, promptMode(tt.flag, getenv, tt.hasTerminal))
		})
	}
}

func TestNonInteractiveFlagReachesPrompts(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		wantInputNeeded bool
	}{
		{name: "--non-interactive", args: []string{"--non-interactive"}, wantInputNeeded: true},
		{name: "--non-interactive=false", args: []string{"--non-interactive=false"}, wantInputNeeded: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given a command that checks whether it may prompt
			var checkErr error
			probe := &cobra.Command{Use: "probe", Run: func(*cobra.Command, []string) {
				checkErr = prompt.Check("Select cluster", "")
			}}
			rootCmd.AddCommand(probe)
			resetRootCmdAfter(t)
			t.Cleanup(func() {
				rootCmd.RemoveCommand(probe)
				prompt.SetMode(prompt.Mode{Interactive: true})
			})

			// When ok runs it with the flag
			rootCmd.SetArgs(append(tt.args, "probe"))
			require.NoError(t, rootCmd.Execute())

			// Then
			var inputNeeded *prompt.InputNeededError
			assert.Equal(t, tt.wantInputNeeded, errors.As(checkErr, &inputNeeded))
		})
	}
}

func TestHelpShowsExitCodesAfterTheFlags(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		lastFlags string
	}{
		{name: "ok", args: []string{"--help"}, lastFlags: "Flags:"},
		{name: "a subcommand", args: []string{"pkg", "install", "--help"}, lastFlags: "Global Flags:"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			var out bytes.Buffer
			resetRootCmdAfter(t)
			rootCmd.SetOut(&out)
			rootCmd.SetArgs(tt.args)

			// When
			require.NoError(t, rootCmd.Execute())

			// Then
			help := out.String()
			flags := strings.Index(help, tt.lastFlags)
			exitCodes := strings.Index(help, "Exit Codes:\n  0  success\n  1  error\n  2  input needed")
			require.NotEqual(t, -1, exitCodes, "no Exit Codes section in:\n%s", help)
			assert.Less(t, flags, exitCodes, "Exit Codes should come after %s", tt.lastFlags)
			assert.NotContains(t, help[:flags], "Exit codes", "the description should not list the exit codes")
		})
	}
}

func TestLoadConfigurationKeepsStdoutClean(t *testing.T) {
	// Given a config file
	configFile := filepath.Join(t.TempDir(), "config.yml")
	require.NoError(t, os.WriteFile(configFile, []byte("enable_experimental: false\n"), 0o600))
	previous := viper.ConfigFileUsed()
	viper.SetConfigFile(configFile)
	t.Cleanup(func() { viper.SetConfigFile(previous) })

	// When the configuration is loaded
	stdout, stderr := captureOutput(t, loadConfiguration)

	// Then the message goes to stderr, not stdout
	assert.Empty(t, stdout)
	assert.Contains(t, stderr, "Using config file: "+configFile)
}

// resetRootCmdAfter undoes, when the test ends, what running rootCmd leaves behind: the arguments,
// the output, and the parsed flags (cobra keeps their values and Changed between runs).
func resetRootCmdAfter(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		rootCmd.SetArgs([]string{}) // not nil: then cobra reads os.Args, which hold go test's flags
		rootCmd.SetOut(nil)
		resetFlags(rootCmd)
	})
}

func resetFlags(c *cobra.Command) {
	reset := func(f *pflag.Flag) {
		_ = f.Value.Set(f.DefValue)
		f.Changed = false
	}
	c.Flags().VisitAll(reset)
	c.PersistentFlags().VisitAll(reset)
	for _, sub := range c.Commands() {
		resetFlags(sub)
	}
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
