package update

import (
	"testing"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/stretchr/testify/require"
)

func TestUpdatedPackages(t *testing.T) {
	tests := []struct {
		name           string
		manifest       common.PackageManifest
		packageName    string
		latestReleases map[string]string
		expected       []common.Package
		expectError    bool
	}{
		{
			name: "update all packages",
			manifest: common.PackageManifest{
				Packages: []common.Package{
					{Template: "template1", Ref: "template1-v1.0.0", OutputFolder: "folder1"},
					{Template: "template2", Ref: "template2-v1.0.0", OutputFolder: "folder2"},
				},
			},
			packageName: "",
			latestReleases: map[string]string{
				"template1": "v1.1.0",
				"template2": "v1.2.0",
			},
			expected: []common.Package{
				{Template: "template1", Ref: "template1-v1.1.0", OutputFolder: "folder1"},
				{Template: "template2", Ref: "template2-v1.2.0", OutputFolder: "folder2"},
			},
			expectError: false,
		},
		{
			name: "update specific package",
			manifest: common.PackageManifest{
				Packages: []common.Package{
					{Template: "template1", Ref: "template1-v1.0.0", OutputFolder: "folder1"},
					{Template: "template2", Ref: "template2-v1.0.0", OutputFolder: "folder2"},
				},
			},
			packageName: "folder1",
			latestReleases: map[string]string{
				"template1": "v1.1.0",
				"template2": "v1.2.0",
			},
			expected: []common.Package{
				{Template: "template1", Ref: "template1-v1.1.0", OutputFolder: "folder1"},
			},
			expectError: false,
		},
		{
			name: "package not found",
			manifest: common.PackageManifest{
				Packages: []common.Package{
					{Template: "template1", Ref: "template1-v1.0.0", OutputFolder: "folder1"},
				},
			},
			packageName: "folder2",
			latestReleases: map[string]string{
				"template1": "v1.1.0",
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "no latest release found",
			manifest: common.PackageManifest{
				Packages: []common.Package{
					{Template: "template1", Ref: "template1-v1.0.0", OutputFolder: "folder1"},
				},
			},
			packageName: "",
			latestReleases: map[string]string{
				"template2": "v1.2.0",
			},
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := updatePackages(tt.manifest, tt.packageName, tt.latestReleases)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestUpdateSpecificPackage(t *testing.T) {
	tests := []struct {
		name           string
		manifest       common.PackageManifest
		packageName    string
		latestReleases map[string]string
		expected       []common.Package
		expectUpdated  bool
		expectError    bool
	}{
		{
			name: "update specific package",
			manifest: common.PackageManifest{
				Packages: []common.Package{
					{Template: "template1", Ref: "template1-v1.0.0", OutputFolder: "folder1"},
					{Template: "template2", Ref: "template2-v1.0.0", OutputFolder: "folder2"},
				},
			},
			packageName: "folder1",
			latestReleases: map[string]string{
				"template1": "v1.1.0",
				"template2": "v1.2.0",
			},
			expected: []common.Package{
				{Template: "template1", Ref: "template1-v1.1.0", OutputFolder: "folder1"},
			},
			expectUpdated: true,
			expectError:   false,
		},
		{
			name: "package not found",
			manifest: common.PackageManifest{
				Packages: []common.Package{
					{Template: "template1", Ref: "template1-v1.0.0", OutputFolder: "folder1"},
				},
			},
			packageName: "folder2",
			latestReleases: map[string]string{
				"template1": "v1.1.0",
			},
			expected:      []common.Package{},
			expectUpdated: false,
			expectError:   false,
		},
		{
			name: "no latest release found",
			manifest: common.PackageManifest{
				Packages: []common.Package{
					{Template: "template1", Ref: "template1-v1.0.0", OutputFolder: "folder1"},
				},
			},
			packageName: "folder1",
			latestReleases: map[string]string{
				"template2": "v1.2.0",
			},
			expected:      []common.Package{},
			expectUpdated: false,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedPackages := []common.Package{}
			updated, err := updatePackage(tt.manifest, tt.packageName, tt.latestReleases, &updatedPackages)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectUpdated, updated)
				require.Equal(t, tt.expected, updatedPackages)
			}
		})
	}
}

func TestGetLastConfigFile(t *testing.T) {
	tests := []struct {
		name     string
		pkg      common.Package
		expected string
		ok       bool
	}{
		{
			name: "package with var files",
			pkg: common.Package{
				VarFiles: []string{"config1.yaml", "config2.yaml"},
			},
			expected: "config2.yaml",
			ok:       true,
		},
		{
			name: "package with one var files",
			pkg: common.Package{
				VarFiles: []string{"config1.yaml"},
			},
			expected: "config1.yaml",
			ok:       true,
		},
		{
			name: "package without var files",
			pkg: common.Package{
				VarFiles: []string{},
			},
			expected: "",
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := getLastConfigFile(tt.pkg)
			require.Equal(t, tt.ok, ok)
			require.Equal(t, tt.expected, result)
		})
	}
}
