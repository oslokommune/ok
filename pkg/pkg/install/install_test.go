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
		}, {
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
			outputFolders:           []string{"out/app-hello"},
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

func TestFilterPackages(t *testing.T) {
	tests := []struct {
		name          string
		packages      []common.Package
		outputFolders []string
		expected      []common.Package
	}{
		{
			name: "no output folders specified",
			packages: []common.Package{
				{OutputFolder: "out/folder1"},
				{OutputFolder: "out/folder2"},
			},
			outputFolders: []string{},
			expected:      []common.Package{},
		},
		{
			name: "single output folder specified",
			packages: []common.Package{
				{OutputFolder: "out/folder1"},
				{OutputFolder: "out/folder2"},
			},
			outputFolders: []string{"out/folder1"},
			expected: []common.Package{
				{OutputFolder: "out/folder1"},
			},
		},
		{
			name: "multiple output folders specified",
			packages: []common.Package{
				{OutputFolder: "out/folder1"},
				{OutputFolder: "out/folder2"},
				{OutputFolder: "out/folder3"},
			},
			outputFolders: []string{"out/folder1", "out/folder3"},
			expected: []common.Package{
				{OutputFolder: "out/folder1"},
				{OutputFolder: "out/folder3"},
			},
		},
		{
			name: "no matching output folders",
			packages: []common.Package{
				{OutputFolder: "out/folder1"},
				{OutputFolder: "out/folder2"},
			},
			outputFolders: []string{"out/folder3"},
			expected:      []common.Package{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterPackages(tt.packages, tt.outputFolders)
			require.Equal(t, tt.expected, result)
		})
	}
}
