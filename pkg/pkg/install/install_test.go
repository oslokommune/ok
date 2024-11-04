package install

import (
	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os/exec"
	"strings"
	"testing"
)

func TestInstall(t *testing.T) {
	testCases := []struct {
		testName                  string
		packageManifestFilename   string
		selectedPackages          []common.Package
		baseUrl                   string
		expectBoilerplateCommands []*exec.Cmd
	}{
		{
			testName:                "Should install all packages from packages.yml",
			packageManifestFilename: "packages.yml",
			expectBoilerplateCommands: []*exec.Cmd{
				exec.Command(
					"boilerplate",
					"--template-url", common.DefaultBaseUrl+"boilerplate/terraform/app?ref=app-v6.1.1",
					"--output-folder", "out/app-hello",
					"--non-interactive",
					"--var-file", "config/common-config.yml",
					"--var-file", "config/app-hello.yml",
				),
				exec.Command(
					"boilerplate",
					"--template-url", common.DefaultBaseUrl+"boilerplate/terraform/networking?ref=main",
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
			baseUrl:                 "git@github.com:oslokommune/SOMETHING_ELSE.git//",
			expectBoilerplateCommands: []*exec.Cmd{
				exec.Command(
					"boilerplate",
					"--template-url", "git@github.com:oslokommune/SOMETHING_ELSE.git//boilerplate/terraform/app?ref=app-v6.1.1",
					"--output-folder", "out/app-hello",
					"--non-interactive",
					"--var-file", "config/common-config.yml",
					"--var-file", "config/app-hello.yml",
				),
			},
		},
		{
			testName:                "Should support file path in BASE_URL",
			packageManifestFilename: "package.yml",
			baseUrl:                 "..",
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
			selectedPackages: []common.Package{
				{
					OutputFolder: "out/app-hello",
					Template:     "app",
					Ref:          "app-v6.1.1",
					VarFiles:     []string{"config/common-config.yml", "config/app-hello.yml"},
				},
			},
			expectBoilerplateCommands: []*exec.Cmd{
				exec.Command(
					"boilerplate",
					"--template-url", common.DefaultBaseUrl+"boilerplate/terraform/app?ref=app-v6.1.1",
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
					"--template-url", common.DefaultBaseUrl+"boilerplate/github-actions/terraform-on-changed-dirs?ref=main",
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
			inputFile, err := common.GetTestdataFilepath(tc.packageManifestFilename)
			require.Nil(t, err)

			manifest, err := common.LoadPackageManifest(inputFile)
			require.Nil(t, err)

			var packages []common.Package
			if len(tc.selectedPackages) > 0 {
				packages = tc.selectedPackages
			} else {
				packages = manifest.Packages
			}

			// When
			cmds, err := CreateBoilerplateCommands(packages, CreateBoilerPlateCommandsOpts{
				PackagePathPrefix: manifest.PackagePrefix(),
				BaseUrlOrPath:     tc.baseUrl,
			})

			// Then
			assert.Nil(t, err)

			for i, cmd := range cmds {
				assert.Equal(t, cmd.Path, tc.expectBoilerplateCommands[i].Path)
				assert.Equal(t,
					strings.Join(tc.expectBoilerplateCommands[i].Args, ","),
					strings.Join(cmd.Args, ","),
				)
			}

			assert.Equal(t, len(cmds), len(tc.expectBoilerplateCommands))
		})
	}
}

func TestIsUrl(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid http URL",
			input:    "http://example.com",
			expected: true,
		},
		{
			name:     "valid https URL",
			input:    "https://example.com",
			expected: true,
		},
		{
			name:     "valid git URL",
			input:    "git@github.com:example/repo.git",
			expected: true,
		},
		{
			name:     "invalid URL",
			input:    "ftp://example.com",
			expected: false,
		},
		{
			name:     "plain text",
			input:    "example.com",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isUrl(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}
