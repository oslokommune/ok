package add

import (
	"testing"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/stretchr/testify/require"
)

func TestCreateNewPackage(t *testing.T) {
	tests := []struct {
		name         string
		manifest     common.PackageManifest
		templateName string
		gitRef       string
		outputFolder string
		expected     common.Package
	}{
		{
			name: "default package",
			manifest: common.PackageManifest{
				DefaultPackagePathPrefix: "",
			},
			templateName: "template1",
			gitRef:       "template1-v1.0.0",
			outputFolder: "folder1",
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
			templateName: "template2",
			gitRef:       "template2-v2.0.0",
			outputFolder: "folder2",
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
			templateName: "template3",
			gitRef:       "template3-v3.0.0",
			outputFolder: "folder3",
			expected: common.Package{
				Template:     "template3",
				Ref:          "template3-v3.0.0",
				OutputFolder: "folder3",
				VarFiles:     []string{"_config/common-config.yml", "_config/folder3.yml"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := createNewPackage(tt.manifest, tt.templateName, tt.gitRef, tt.outputFolder)
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
