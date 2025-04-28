package add

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/stretchr/testify/require"
)

func TestCreateNewPackage(t *testing.T) {
	tests := []struct {
		name                         string
		manifest                     common.PackageManifest
		templateName                 string
		gitRef                       string
		outputFolder                 string
		consolidatedPackageStructure bool
		expected                     common.Package
	}{
		{
			name: "default package",
			manifest: common.PackageManifest{
				DefaultPackagePathPrefix: "",
			},
			templateName:                 "template1",
			gitRef:                       "template1-v1.0.0",
			outputFolder:                 "folder1",
			consolidatedPackageStructure: true,
			expected: common.Package{
				Template:     "template1",
				Ref:          "template1-v1.0.0",
				OutputFolder: "folder1",
				VarFiles:     []string{"_config/common-config.yml", "_config/folder1.yml"},
			},
		},
		{
			name: "GitHub Actions package",
			manifest: common.PackageManifest{
				DefaultPackagePathPrefix: common.BoilerplatePackageGitHubActionsPath,
			},
			templateName:                 "template2",
			gitRef:                       "template2-v2.0.0",
			outputFolder:                 "folder2",
			consolidatedPackageStructure: true,
			expected: common.Package{
				Template:     "template2",
				Ref:          "template2-v2.0.0",
				OutputFolder: "../..",
				VarFiles:     []string{"common-config.yml", "folder2.yml"},
			},
		},
		{
			name: "custom package prefix that doesn't exist should use the default",
			manifest: common.PackageManifest{
				DefaultPackagePathPrefix: "custom/prefix",
			},
			templateName:                 "template3",
			gitRef:                       "template3-v3.0.0",
			outputFolder:                 "folder3",
			consolidatedPackageStructure: true,
			expected: common.Package{
				Template:     "template3",
				Ref:          "template3-v3.0.0",
				OutputFolder: "folder3",
				VarFiles:     []string{"_config/common-config.yml", "_config/folder3.yml"},
			},
		},
		{
			name: "packages.yml in separate folder",
			manifest: common.PackageManifest{
				DefaultPackagePathPrefix: "",
			},
			templateName:                 "template4",
			gitRef:                       "template4-v4.0.0",
			outputFolder:                 "folder4",
			consolidatedPackageStructure: false,
			expected: common.Package{
				Template:     "template4",
				Ref:          "template4-v4.0.0",
				OutputFolder: ".",
				VarFiles:     []string{"../common-config.yml", "config.yml"},
			},
		},
		{
			name: "packages.yml in separate, nested folder",
			manifest: common.PackageManifest{
				DefaultPackagePathPrefix: "",
			},
			templateName:                 "template5",
			gitRef:                       "template5-v5.0.0",
			outputFolder:                 "dir/folder5",
			consolidatedPackageStructure: false,
			expected: common.Package{
				Template:     "template5",
				Ref:          "template5-v5.0.0",
				OutputFolder: ".",
				VarFiles:     []string{"../../common-config.yml", "config.yml"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := createNewPackage(tt.manifest, tt.templateName, tt.gitRef, tt.outputFolder, tt.consolidatedPackageStructure)
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestAllowDuplicateOutputFolder(t *testing.T) {
	tests := []struct {
		name          string
		manifest      common.PackageManifest
		newPackage    common.Package
		expectedError bool
	}{
		{
			name: "no duplicate output folder",
			manifest: common.PackageManifest{
				Packages: []common.Package{
					{OutputFolder: "folder1"},
					{OutputFolder: "folder2"},
				},
			},
			newPackage:    common.Package{OutputFolder: "folder3"},
			expectedError: false,
		},
		{
			name: "duplicate output folder",
			manifest: common.PackageManifest{
				Packages: []common.Package{
					{OutputFolder: "folder1"},
					{OutputFolder: "folder2"},
				},
			},
			newPackage:    common.Package{OutputFolder: "folder2"},
			expectedError: true,
		},
		{
			name: "GHA package prefix allows duplicates",
			manifest: common.PackageManifest{
				Packages: []common.Package{
					{OutputFolder: "folder1"},
					{OutputFolder: "folder2"},
				},
				DefaultPackagePathPrefix: common.BoilerplatePackageGitHubActionsPath,
			},
			newPackage:    common.Package{OutputFolder: "folder2"},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := allowDuplicateOutputFolder(tt.manifest, tt.newPackage)
			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUpdateSchemaConfig(t *testing.T) {
	// Create temporary test directory
	testDir, err := os.MkdirTemp("", "schema-config-test")
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	tests := []struct {
		name                         string
		consolidatedPackageStructure bool
		outputFolder                 string
		expectedFilePath             string
		expectedError                bool
	}{
		{
			name:                         "consolidated package structure",
			consolidatedPackageStructure: true,
			outputFolder:                 "my-output",
			expectedFilePath:             "_config/my-output.yml",
			expectedError:                false,
		},
		{
			name:                         "non-consolidated package structure",
			consolidatedPackageStructure: false,
			outputFolder:                 "my-output",
			expectedFilePath:             "my-output/config.yml",
			expectedError:                false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a subdirectory for each test case
			testCaseDir := filepath.Join(testDir, tt.name)
			err := os.MkdirAll(testCaseDir, 0755)
			require.NoError(t, err)

			// Change to test case directory
			originalDir, err := os.Getwd()
			require.NoError(t, err)
			err = os.Chdir(testCaseDir)
			require.NoError(t, err)
			defer os.Chdir(originalDir) // Restore original directory at end

			// Create test manifest
			manifest := common.PackageManifest{}

			// Create test package
			pkg := common.Package{
				Template: "test-template",
				Ref:      "test-template-v1.0.0",
			}

			// Create adder and run update
			adder := NewAdder()
			err = adder.updateSchemaConfig(manifest, pkg, tt.outputFolder, tt.consolidatedPackageStructure)

			// Check results
			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Verify file was created
				fullPath := filepath.Join(testCaseDir, tt.expectedFilePath)
				fileInfo, err := os.Stat(fullPath)
				require.NoError(t, err)
				require.False(t, fileInfo.IsDir())

				// Read file content and verify it contains schema reference
				content, err := os.ReadFile(fullPath)
				require.NoError(t, err)
				require.Contains(t, string(content), "schema")
				require.Contains(t, string(content), pkg.Ref)
			}
		})
	}
}
