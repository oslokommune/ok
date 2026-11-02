package prompt

import (
	"os"
	"testing"

	"github.com/charmbracelet/huh"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNonInteractive(t *testing.T) {
	t.Run("the flag turns prompts off", func(t *testing.T) {
		t.Setenv("CI", "")
		setFlag(t, true)

		assert.True(t, NonInteractive())
	})

	t.Run("CI=true turns prompts off", func(t *testing.T) {
		t.Setenv("CI", "true")
		setFlag(t, false)

		assert.True(t, NonInteractive())
	})

	t.Run("stdin that is not a terminal turns prompts off", func(t *testing.T) {
		t.Setenv("CI", "")
		setFlag(t, false)
		setStdinToPipe(t)

		assert.True(t, NonInteractive())
	})
}

func TestRun(t *testing.T) {
	t.Run("returns InputNeededError without showing the form", func(t *testing.T) {
		t.Setenv("CI", "true")

		err := Run(huh.NewForm(huh.NewGroup(huh.NewInput())), "Select cluster", "pass --cluster")

		var inputNeeded *InputNeededError
		require.ErrorAs(t, err, &inputNeeded)
		assert.Equal(t, "Select cluster", inputNeeded.Prompt)
		assert.Equal(t, "pass --cluster", inputNeeded.Hint)
	})
}

func TestInputNeededError(t *testing.T) {
	t.Run("names the prompt, the flag and the hint", func(t *testing.T) {
		err := &InputNeededError{Prompt: "Select cluster", Hint: "pass --cluster"}

		assert.Equal(t,
			`ok needs input ("Select cluster") but runs non-interactively (--non-interactive, OK_NON_INTERACTIVE, CI=true or no terminal): pass --cluster`,
			err.Error())
	})

	t.Run("leaves the hint out when there is none", func(t *testing.T) {
		err := &InputNeededError{Prompt: "Select cluster"}

		assert.Equal(t,
			`ok needs input ("Select cluster") but runs non-interactively (--non-interactive, OK_NON_INTERACTIVE, CI=true or no terminal)`,
			err.Error())
	})
}

func setFlag(t *testing.T, value bool) {
	t.Helper()
	viper.Set(ConfigKey, value)
	t.Cleanup(func() { viper.Set(ConfigKey, false) })
}

func setStdinToPipe(t *testing.T) {
	t.Helper()
	r, w, err := os.Pipe()
	require.NoError(t, err)
	oldStdin := os.Stdin
	os.Stdin = r
	t.Cleanup(func() {
		os.Stdin = oldStdin
		_ = r.Close()
		_ = w.Close()
	})
}
