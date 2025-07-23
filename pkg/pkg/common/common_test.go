package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigFile(t *testing.T) {
	tests := []struct {
		name       string
		prefix     string
		configName string
		expected   string
	}{
		{
			name:       "no prefix",
			prefix:     "",
			configName: "config",
			expected:   "config.yml",
		},
		{
			name:       "with prefix",
			prefix:     "prefix",
			configName: "config",
			expected:   "prefix/config.yml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := VarFile(tt.prefix, tt.configName)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestUseConsolidatedPackageStructure(t *testing.T) {
	testDir, err := os.MkdirTemp("", "package-structure-test")
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	tests := []struct {
		name           string
		setup          func(baseDir string, t *testing.T)
		expectedResult bool
		expectedError  bool
	}{
		{
			name: "should use existing packages.yml when packages.yml and _config found in dir",
			setup: func(baseDir string, t *testing.T) {
				os.WriteFile(filepath.Join(baseDir, PackagesManifestFilename), []byte(""), 0644)
				require.NoError(t, err)
				os.MkdirAll(filepath.Join(baseDir, BoilerplatePackageTerraformConfigPrefix), 0755)
				require.NoError(t, err)
			},
			expectedResult: true,
			expectedError:  false,
		},
		{
			name: "should use existing packages.yml for github actions",
			setup: func(baseDir string, t *testing.T) {
				content := "DefaultPackagePathPrefix: " + BoilerplatePackageGitHubActionsPath
				os.WriteFile(filepath.Join(baseDir, PackagesManifestFilename), []byte(content), 0644)
				require.NoError(t, err)
			},
			expectedResult: true,
			expectedError:  false,
		},
		{
			name:           "should not use consolidated packages.yml by default",
			setup:          func(baseDir string, t *testing.T) {},
			expectedResult: false,
			expectedError:  false,
		},
		{
			name:           "should not use consolidated packages.yml when only packages.yml is present in dir",
			setup:          func(baseDir string, t *testing.T) {},
			expectedResult: false,
			expectedError:  false,
		},
		{
			name: "should not use consolidated packages.yml when only _config is present in dir",
			setup: func(baseDir string, t *testing.T) {
				os.MkdirAll(filepath.Join(baseDir, BoilerplatePackageTerraformConfigPrefix), 0755)
				require.NoError(t, err)
			},
			expectedResult: false,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a subdirectory for each test case
			testCaseDir := filepath.Join(testDir, tt.name)
			err := os.MkdirAll(testCaseDir, 0755)
			require.NoError(t, err)

			tt.setup(testCaseDir, t)

			result, err := UseOldPackageStructure(testCaseDir)

			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.expectedResult, result)
		})
	}
}
