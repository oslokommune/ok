package scriptrunner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunScript(t *testing.T) {
	t.Run("returns an error for a script that is not embedded", func(t *testing.T) {
		err := RunScript("does-not-exist.sh", nil, true)

		assert.ErrorContains(t, err, "reading script")
	})
}

func TestExecuteScript(t *testing.T) {
	t.Run("returns an error when the script fails", func(t *testing.T) {
		script := writeScript(t, "exit 3\n")

		assert.ErrorContains(t, executeScript(script, nil, true), "executing script")
	})

	t.Run("returns nil when the script succeeds", func(t *testing.T) {
		script := writeScript(t, "exit 0\n")

		assert.NoError(t, executeScript(script, nil, true))
	})

	t.Run("gives the script an empty stdin when it may not prompt", func(t *testing.T) {
		// Given a line waiting on ok's stdin, and a script that records whether read got a line
		r, w, err := os.Pipe()
		require.NoError(t, err)
		_, err = w.WriteString("yes\n")
		require.NoError(t, err)
		require.NoError(t, w.Close())
		oldStdin := os.Stdin
		os.Stdin = r
		t.Cleanup(func() { os.Stdin = oldStdin; _ = r.Close() })

		result := filepath.Join(t.TempDir(), "result")
		script := writeScript(t, `if read -r line; then echo line > "$1"; else echo eof > "$1"; fi`+"\n")

		// When
		require.NoError(t, executeScript(script, []string{result}, false))

		// Then
		got, err := os.ReadFile(result)
		require.NoError(t, err)
		assert.Equal(t, "eof\n", string(got))
	})
}

func writeScript(t *testing.T, content string) string {
	t.Helper()
	script := filepath.Join(t.TempDir(), "script.sh")
	require.NoError(t, os.WriteFile(script, []byte(content), 0o600))
	return script
}
