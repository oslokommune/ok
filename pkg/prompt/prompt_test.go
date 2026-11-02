package prompt

import (
	"testing"

	"github.com/charmbracelet/huh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Run("returns InputNeededError without showing the form when ok may not prompt", func(t *testing.T) {
		// Given
		setMode(t, Mode{Reason: "CI=true is set"})

		// When
		err := Run(huh.NewForm(huh.NewGroup(huh.NewInput())), "Select cluster", "Pass --cluster")

		// Then
		var inputNeeded *InputNeededError
		require.ErrorAs(t, err, &inputNeeded)
		assert.Equal(t, &InputNeededError{Prompt: "Select cluster", Reason: "CI=true is set", Hint: "Pass --cluster"}, inputNeeded)
	})
}

func TestCheck(t *testing.T) {
	t.Run("passes when ok may prompt", func(t *testing.T) {
		setMode(t, Mode{Interactive: true})

		assert.NoError(t, Check("Select cluster", ""))
	})

	t.Run("returns InputNeededError when ok may not prompt", func(t *testing.T) {
		setMode(t, Mode{Reason: "there is no terminal"})

		var inputNeeded *InputNeededError
		assert.ErrorAs(t, Check("Select cluster", ""), &inputNeeded)
	})
}

func TestInputNeededError(t *testing.T) {
	err := &InputNeededError{Prompt: "Select cluster", Reason: "CI=true is set", Hint: "Pass --cluster"}

	assert.Equal(t, `ok needs an answer to "Select cluster" but cannot ask: CI=true is set`, err.Error())
	assert.Equal(t, "Pass --cluster", err.Details())
	assert.Equal(t, 2, err.ExitCode())
}

func setMode(t *testing.T, m Mode) {
	t.Helper()
	old := mode
	SetMode(m)
	t.Cleanup(func() { SetMode(old) })
}
