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

func TestInstall(t *testing.T) {
	testCases := []struct {
		testName                  string
		packageManifestFilename   string
		expectBoilerplateCommands []*exec.Cmd
	}{
		{
			testName:                "Should run correct boilerplate commands from package manifest",
			packageManifestFilename: "packages.yml",
			expectBoilerplateCommands: []*exec.Cmd{
				exec.Command(
					"boilerplate",
					"--template-url", BaseUrl+"/app?ref=app-v6.1.1",
					"--output-folder", "out/app-hello",
					"--non-interactive",
					"--var-file", "config/common-config.yml",
					"--var-file", "config/app-hello.yml",
				),
				exec.Command(
					"boilerplate",
					"--template-url", BaseUrl+"/networking?ref=main",
					"--output-folder", "out/networking",
					"--non-interactive",
					"--var-file", "config/common-config.yml",
					"--var-file", "config/networking.yml",
				),
			},
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
			inputFile := filepath.Join(cwd, "testdata", tc.packageManifestFilename)

			// When
			cmds, err := CreateBoilerplateCommands(inputFile)

			// Then
			assert.NilError(t, err)

			for i, cmd := range cmds {
				assert.Equal(t, cmd.Path, tc.expectBoilerplateCommands[i].Path)
				assert.Equal(t,
					strings.Join(cmd.Args, ","),
					strings.Join(tc.expectBoilerplateCommands[i].Args, ","))
			}

			assert.Equal(t, len(cmds), len(tc.expectBoilerplateCommands))
		})
	}
}
