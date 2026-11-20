package scriptrunner

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeAws answers the aws calls port-forward.sh makes, and writes the arguments of
// "aws ssm start-session" to $FAKE_AWS_SESSION.
const fakeAws = `#!/bin/bash
case "$*" in
    "ecs list-task-definitions "*) echo '"arn:aws:ecs:eu-west-1:111111111111:task-definition/ssm-portfwd-test:1"' ;;
    "ecs list-tags-for-resource "*) echo '{"tags":[{"key":"OkTool:SecurityGroup","value":"sg-1"},{"key":"OkTool:Subnets","value":"subnet-1/subnet-2"},{"key":"OkTool:ClusterName","value":"cluster"}]}' ;;
    "ecs run-task "*) echo 'arn:aws:ecs:eu-west-1:111111111111:task/cluster/task1' ;;
    "ecs describe-tasks "*runtimeId*) echo 'runtime-1' ;;
    "ecs describe-tasks "*containers*) echo '[{"lastStatus":"RUNNING","managedAgents":[{"lastStatus":"RUNNING"}]}]' ;;
    "ecs describe-tasks "*) echo '{"tasks":[{"lastStatus":"RUNNING"}]}' ;;
    "ssm start-session "*) echo "$*" > "$FAKE_AWS_SESSION" ;;
    "ecs stop-task "*) echo '{"task":{"desiredStatus":"STOPPED"}}' ;;
    *) echo "fake aws: unexpected call: $*" >&2; exit 1 ;;
esac
`

func TestPortForwardScript(t *testing.T) {
	if _, err := exec.LookPath("jq"); err != nil {
		t.Skip("port-forward.sh needs jq")
	}

	t.Run("a variable that is set wins over its saved value, and one that is not set is read from the saved file", func(t *testing.T) {
		// Given saved selections from an earlier ok forward --keep
		home := t.TempDir()
		saved := "LOCAL_PORT=1111\nTASK_DEF_ARN=arn:aws:ecs:eu-west-1:111111111111:task-definition/saved:1\nREMOTE_PORT=2222\nRDS_ENDPOINT=saved.example.com\n"
		require.NoError(t, os.WriteFile(filepath.Join(home, ".ok-port-forward.conf"), []byte(saved), 0o600))

		bin := t.TempDir()
		writeExecutable(t, filepath.Join(bin, "aws"), fakeAws)
		writeExecutable(t, filepath.Join(bin, "fzf"), "#!/bin/bash\necho 'fake fzf: no picker expected' >&2\nexit 1\n")
		writeExecutable(t, filepath.Join(bin, "session-manager-plugin"), "#!/bin/bash\n")

		session := filepath.Join(t.TempDir(), "session")
		t.Setenv("HOME", home)
		t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
		t.Setenv("FAKE_AWS_SESSION", session)

		// And the environment sets the task definition family, the endpoint and the local port, but not the remote port
		t.Setenv("TASK_DEF_FAMILY", "ssm-portfwd-test")
		t.Setenv("RDS_ENDPOINT", "env.example.com")
		t.Setenv("LOCAL_PORT", "4813")
		t.Setenv("REMOTE_PORT", "")

		// When
		require.NoError(t, RunScript("port-forward.sh", nil, false))

		// Then
		got, err := os.ReadFile(session)
		require.NoError(t, err)
		assert.Contains(t, string(got), "--parameters host=env.example.com,portNumber=2222,localPortNumber=4813")
	})
}

func writeExecutable(t *testing.T, path, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(path, []byte(content), 0o700))
}
