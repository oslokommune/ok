package install

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstall(t *testing.T) {
	testCases := []struct {
		testName                  string
		packageManifestFilename   string
		expectBoilerplateCommands []*exec.Cmd
		outputFolders             []string
		baseUrl                   string
	}{
		{
			testName:                "Should install all packages from packages.yml",
			packageManifestFilename: "packages.yml",
			expectBoilerplateCommands: []*exec.Cmd{
				exec.Command(
					"boilerplate",
					"--template-url", DefaultBaseUrl+"boilerplate/terraform/app?ref=app-v6.1.1",
					"--output-folder", "out/app-hello",
					"--non-interactive",
					"--var-file", "config/common-config.yml",
					"--var-file", "config/app-hello.yml",
				),
				exec.Command(
					"boilerplate",
					"--template-url", DefaultBaseUrl+"boilerplate/terraform/networking?ref=main",
					"--output-folder", "out/networking",
					"--non-interactive",
					"--var-file", "config/common-config.yml",
					"--var-file", "config/networking.yml",
				),
			},
		},
		{
			testName:                "Should support URL in BASE_URL",
			packageManifestFilename: "package.yml",
			baseUrl:                 DefaultBaseUrl,
			expectBoilerplateCommands: []*exec.Cmd{
				exec.Command(
					"boilerplate",
					"--template-url", DefaultBaseUrl+"boilerplate/terraform/app?ref=app-v6.1.1",
					"--output-folder", "out/app-hello",
					"--non-interactive",
					"--var-file", "config/common-config.yml",
					"--var-file", "config/app-hello.yml",
				),
			},
		}, {
			testName:                "Should support file path in BASE_URL",
			packageManifestFilename: "package.yml",
			baseUrl:                 "../",
			expectBoilerplateCommands: []*exec.Cmd{
				exec.Command(
					"boilerplate",
					"--template-url", "../boilerplate/terraform/app",
					"--output-folder", "out/app-hello",
					"--non-interactive",
					"--var-file", "config/common-config.yml",
					"--var-file", "config/app-hello.yml",
				),
			},
		},
		{
			testName:                "Should install package with specified output folder",
			packageManifestFilename: "packages.yml",
			outputFolders:           []string{"out/app-hello"},
			expectBoilerplateCommands: []*exec.Cmd{
				exec.Command(
					"boilerplate",
					"--template-url", DefaultBaseUrl+"boilerplate/terraform/app?ref=app-v6.1.1",
					"--output-folder", "out/app-hello",
					"--non-interactive",
					"--var-file", "config/common-config.yml",
					"--var-file", "config/app-hello.yml",
				),
			},
		},
		{
			testName:                "Should install package correctly when template path prefix set",
			packageManifestFilename: "packages-github-actions.yml",
			expectBoilerplateCommands: []*exec.Cmd{
				exec.Command(
					"boilerplate",
					"--template-url", DefaultBaseUrl+"boilerplate/github-actions/terraform-on-changed-dirs?ref=main",
					"--output-folder", "out/.github/workflows",
					"--non-interactive",
					"--var-file", "config/common-config.yml",
					"--var-file", "config/terraform-on-changed-dirs.yml",
				),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.testName, func(t *testing.T) {
			// Given
			inputFile, err := getTestdataFilepath(tc.packageManifestFilename)
			require.Nil(t, err)

			// When
			cmds, err := CreateBoilerplateCommands(inputFile, tc.outputFolders, tc.baseUrl)

			// Then
			assert.Nil(t, err)

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

func getTestdataFilepath(testDataFilename string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting current directory: %w", err)
	}

	return filepath.Join(cwd, "testdata", testDataFilename), nil
}
