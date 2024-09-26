package add

import (
	"testing"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/stretchr/testify/require"
)

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
