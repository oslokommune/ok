// Package prompt shows huh forms only when ok may ask, and otherwise says what ok needed.
package prompt

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/mattn/go-isatty"
)

// ExitCodeInputNeeded is ok's exit code when it needs input but may not prompt.
const ExitCodeInputNeeded = 2

// Mode says whether ok may prompt, and if not, why.
type Mode struct {
	Interactive bool
	Reason      string // why ok may not prompt, e.g. "CI=true is set"
}

var mode = Mode{Interactive: true}

// SetMode sets whether ok may prompt. The root command calls it before any command runs.
func SetMode(m Mode) {
	mode = m
}

// InputNeededError is returned instead of showing a prompt when ok may not prompt.
type InputNeededError struct {
	Prompt string // the prompt's title, e.g. "Select cluster"
	Reason string // why ok may not prompt
	Hint   string // how to give the answer without a prompt; empty when there is no way yet
}

func (e *InputNeededError) Error() string {
	return fmt.Sprintf("ok needs an answer to %q but cannot ask: %s", e.Prompt, e.Reason)
}

// Details is printed under the error.
func (e *InputNeededError) Details() string {
	return e.Hint
}

func (e *InputNeededError) ExitCode() int {
	return ExitCodeInputNeeded
}

// Allowed reports whether ok may prompt.
func Allowed() bool {
	return mode.Interactive
}

// Check returns *InputNeededError when ok may not prompt. Use it before a command does anything
// it cannot undo, when a prompt will follow.
func Check(title string, hint string) error {
	if !Allowed() {
		return &InputNeededError{Prompt: title, Reason: mode.Reason, Hint: hint}
	}
	return nil
}

// Run shows the form, or returns *InputNeededError without showing it when ok may not prompt.
func Run(form *huh.Form, title string, hint string) error {
	if err := Check(title, hint); err != nil {
		return err
	}
	return form.Run()
}

// HasTerminal reports whether there is a terminal to prompt on: stdin, or /dev/tty, which huh
// opens when stdin is redirected.
func HasTerminal() bool {
	fd := os.Stdin.Fd()
	if isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd) {
		return true
	}
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return false
	}
	_ = tty.Close()
	return true
}
