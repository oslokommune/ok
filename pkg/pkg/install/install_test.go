package install

import (
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

import (
	"gotest.tools/v3/assert"
	"testing"
)

const baseUrl = "someBaseUrl"

func TestInstall(t *testing.T) {
	testCases := []struct {
		testName                 string
		inputFile                string
		expectBoilerplateCommand *exec.Cmd
	}{
		{
			testName:  "Should work",
			inputFile: "package.yml",
			expectBoilerplateCommand: exec.Command(
				"boilerplate",
				"--template-url", "%s/app/ref=v1.0.2",
				"--var-file", "common-config.yml",
				"--template-url", "app-hello.yml",
				"--output-folder", "app-hello",
				"--non-interactive",
			),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			cwd, err := os.Getwd()
			if err != nil {
				require.Nil(t, err)
			}

			// Construct the absolute path to testdata/packages.yml
			inputFile := filepath.Join(cwd, "testdata", "packages.yml")

			// When
			cmd, err := CreateBoilerplateCommand(inputFile)

			// Then
			assert.NilError(t, err)

			assert.Equal(t, cmd.Path, tc.expectBoilerplateCommand.Path)
			assert.Equal(t,
				strings.Join(cmd.Args, ","),
				strings.Join(tc.expectBoilerplateCommand.Args, ","))
		})
	}
}
