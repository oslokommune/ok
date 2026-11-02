// Package prompt decides whether ok may ask the user anything, and shows huh forms only when it may.
package prompt

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/mattn/go-isatty"
	"github.com/spf13/viper"
)

const (
	// FlagName is the persistent flag that turns prompts off. OK_NON_INTERACTIVE=true does the same.
	FlagName = "non-interactive"

	// ConfigKey is the viper key bound to FlagName and OK_NON_INTERACTIVE.
	ConfigKey = "non_interactive"

	// ExitCodeInputNeeded is the exit code when ok needs input it is not allowed to ask for.
	ExitCodeInputNeeded = 2
)

// InputNeededError is returned instead of showing a prompt when ok runs non-interactively.
type InputNeededError struct {
	Prompt string // what ok would have asked, e.g. "Select cluster"
	Hint   string // how to give the answer without a prompt; empty when there is no way yet
}

func (e *InputNeededError) Error() string {
	msg := fmt.Sprintf("ok needs input (%q) but runs non-interactively (--%s, OK_NON_INTERACTIVE, CI=true or no terminal)", e.Prompt, FlagName)
	if e.Hint != "" {
		msg += ": " + e.Hint
	}
	return msg
}

// NonInteractive reports whether prompts are off: --non-interactive, OK_NON_INTERACTIVE=true,
// CI=true, or stdin is not a terminal.
func NonInteractive() bool {
	if viper.GetBool(ConfigKey) || os.Getenv("CI") == "true" {
		return true
	}
	fd := os.Stdin.Fd()
	return !isatty.IsTerminal(fd) && !isatty.IsCygwinTerminal(fd)
}

// Run shows the form, or returns *InputNeededError without showing it when prompts are off.
func Run(form *huh.Form, prompt string, hint string) error {
	if NonInteractive() {
		return &InputNeededError{Prompt: prompt, Hint: hint}
	}
	return form.Run()
}
