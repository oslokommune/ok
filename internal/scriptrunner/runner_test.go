package scriptrunner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunScriptReturnsErrors(t *testing.T) {
	t.Run("unknown script", func(t *testing.T) {
		err := RunScript("does-not-exist.sh", nil)

		assert.ErrorContains(t, err, "reading script")
	})

	t.Run("script that fails", func(t *testing.T) {
		script := filepath.Join(t.TempDir(), "fail.sh")
		require.NoError(t, os.WriteFile(script, []byte("exit 3\n"), 0o600))

		err := executeScript(script, nil)

		assert.ErrorContains(t, err, "executing script")
	})

	t.Run("script that succeeds", func(t *testing.T) {
		script := filepath.Join(t.TempDir(), "ok.sh")
		require.NoError(t, os.WriteFile(script, []byte("exit 0\n"), 0o600))

		assert.NoError(t, executeScript(script, nil))
	})
}
