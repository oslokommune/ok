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
		err := RunScript("does-not-exist.sh", nil)

		assert.ErrorContains(t, err, "reading script")
	})
}

func TestExecuteScript(t *testing.T) {
	t.Run("returns an error when the script fails", func(t *testing.T) {
		script := filepath.Join(t.TempDir(), "fail.sh")
		require.NoError(t, os.WriteFile(script, []byte("exit 3\n"), 0o600))

		assert.ErrorContains(t, executeScript(script, nil), "executing script")
	})

	t.Run("returns nil when the script succeeds", func(t *testing.T) {
		script := filepath.Join(t.TempDir(), "ok.sh")
		require.NoError(t, os.WriteFile(script, []byte("exit 0\n"), 0o600))

		assert.NoError(t, executeScript(script, nil))
	})
}
