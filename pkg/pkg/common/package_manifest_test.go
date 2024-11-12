package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackageManifestPackagePrefix(t *testing.T) {
	tests := []struct {
		name     string
		manifest PackageManifest
		expected string
	}{
		{
			name: "default package path prefix",
			manifest: PackageManifest{
				DefaultPackagePathPrefix: "",
			},
			expected: DefaultPackagePathPrefix,
		},
		{
			name: "custom package path prefix",
			manifest: PackageManifest{
				DefaultPackagePathPrefix: "custom/prefix",
			},
			expected: "custom/prefix",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.manifest.PackagePrefix()
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestPackageManifestPackageConfigPrefix(t *testing.T) {
	tests := []struct {
		name     string
		manifest PackageManifest
		expected string
	}{
		{
			name: "default package config prefix when no known package path prefix is configured",
			manifest: PackageManifest{
				DefaultPackagePathPrefix: "boilerplate/unknown",
			},
			expected: DefaultPackageConfigPrefix,
		},
		{
			name: "GitHub Actions package config prefix",
			manifest: PackageManifest{
				DefaultPackagePathPrefix: BoilerplatePackageGitHubActionsPath,
			},
			expected: BoilerplatePackageGitHubActionsConfigPrefix,
		},
		{
			name: "Terraform package config prefix",
			manifest: PackageManifest{
				DefaultPackagePathPrefix: BoilerplatePackageTerraformPath,
			},
			expected: BoilerplatePackageTerraformConfigPrefix,
		},
		{
			name: "custom package config prefix",
			manifest: PackageManifest{
				DefaultPackagePathPrefix: "custom/prefix",
			},
			expected: DefaultPackageConfigPrefix,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.manifest.PackageConfigPrefix()
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestPackageManifestPackageOutputFolder(t *testing.T) {
	tests := []struct {
		name         string
		manifest     PackageManifest
		outputFolder string
		expected     string
	}{
		{
			name: "default output folder",
			manifest: PackageManifest{
				DefaultPackagePathPrefix: "",
			},
			outputFolder: "output/folder",
			expected:     "output/folder",
		},
		{
			name: "GitHub Actions output folder",
			manifest: PackageManifest{
				DefaultPackagePathPrefix: BoilerplatePackageGitHubActionsPath,
			},
			outputFolder: "output/folder",
			expected:     BoilerplatePackageGitHubActionsOutputFolder,
		},
		{
			name: "custom output folder",
			manifest: PackageManifest{
				DefaultPackagePathPrefix: "custom/prefix",
			},
			outputFolder: "output/folder",
			expected:     "output/folder",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.manifest.PackageOutputFolder(tt.outputFolder)
			require.Equal(t, tt.expected, result)
		})
	}
}
