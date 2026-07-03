package update

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/stretchr/testify/require"
)

func TestUpdatedPackages(t *testing.T) {
	tests := []struct {
		name             string
		manifest         common.PackageManifest
		packagesToUpdate []common.Package
		latestReleases   map[string]string
		expected         []common.Package
		expectError      bool
	}{
		{
			name: "update all packages",
			manifest: common.PackageManifest{
				Packages: []common.Package{
					{Template: "template1", Ref: "template1-v1.0.0", OutputFolder: "folder1"},
					{Template: "template2", Ref: "template2-v1.0.0", OutputFolder: "folder2"},
				},
			},
			packagesToUpdate: []common.Package{
				{Template: "template1", Ref: "template1-v1.0.0", OutputFolder: "folder1"},
				{Template: "template2", Ref: "template2-v1.0.0", OutputFolder: "folder2"},
			},
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
			packagesToUpdate: []common.Package{
				{Template: "template1", Ref: "template1-v1.0.0", OutputFolder: "folder1"},
			},
			latestReleases: map[string]string{
				"template1": "v1.1.0",
				"template2": "v1.2.0",
			},
			expected: []common.Package{
				{Template: "template1", Ref: "template1-v1.1.0", OutputFolder: "folder1"},
				{Template: "template2", Ref: "template2-v1.0.0", OutputFolder: "folder2"},
			},
			expectError: false,
		},
		{
			name: "no latest release found",
			manifest: common.PackageManifest{
				Packages: []common.Package{
					{Template: "template1", Ref: "template1-v1.0.0", OutputFolder: "folder1"},
				},
			},
			packagesToUpdate: []common.Package{
				{Template: "template1", Ref: "template1-v1.0.0", OutputFolder: "folder1"},
				{Template: "template2", Ref: "template2-v1.0.0", OutputFolder: "folder2"},
			},
			latestReleases: map[string]string{
				"template2": "v1.2.0",
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "don't update non-semver package Refs",
			manifest: common.PackageManifest{
				Packages: []common.Package{
					{Template: "app", Ref: "main", OutputFolder: "app-hello"},
				},
			},
			packagesToUpdate: []common.Package{
				{Template: "app", Ref: "main", OutputFolder: "app-hello"},
			},
			latestReleases: map[string]string{
				"app": "v2.0.0",
			},
			expected: []common.Package{
				{Template: "app", Ref: "main", OutputFolder: "app-hello"},
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			updatedManifest, _, err := updatePackages(tc.manifest, tc.packagesToUpdate, tc.latestReleases)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, updatedManifest.Packages)
			}
		})
	}
}

// schemaLine returns a JSON schema declaration line as scaffolded by `ok pkg add`,
// for a package ref such as "load-balancing-alb-v5.1.0".
func schemaLine(ref string) string {
	return fmt.Sprintf(
		"# yaml-language-server: $schema=https://raw.githubusercontent.com/oslokommune/golden-path-boilerplate-schemas/refs/heads/main/schemas/%s.schema.json",
		ref)
}

// captureStderr runs fn and returns everything fn wrote to os.Stderr.
func captureStderr(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stderr
	r, w, err := os.Pipe()
	require.NoError(t, err)

	os.Stderr = w
	defer func() { os.Stderr = old }()

	fn()

	require.NoError(t, w.Close())
	data, err := io.ReadAll(r)
	require.NoError(t, err)

	return string(data)
}

func TestSetJsonSchemaDeclarationInVarFiles(t *testing.T) {
	// File contents modeled on a real stack layout:
	//   dev/common-config.yml                          (shared across stacks, no schema declaration)
	//   dev/load-balancing-alb-nhn/package-config.yml  (scaffolded schema declaration)
	// with workingDirectory = dev/load-balancing-alb-nhn.
	commonConfig := "AccountId: \"123456789012\"\nRegion: \"eu-north-1\"\nTeam: \"legevaktmottak\"\n"
	packageConfigOld := schemaLine("load-balancing-alb-v5.1.0") + "\nStackName: \"load-balancing-alb-nhn\"\n"
	packageConfigNew := schemaLine("load-balancing-alb-v5.2.0") + "\nStackName: \"load-balancing-alb-nhn\"\n"

	tests := []struct {
		name string
		// files maps a path relative to the repo root to its content. workingDirectory is dev/load-balancing-alb-nhn.
		files           map[string]string
		packages        []common.Package
		expectedFiles   map[string]string
		expectedWarning string
	}{
		{
			name: "conventional order writes declaration into the declared package config",
			files: map[string]string{
				"dev/common-config.yml":                         commonConfig,
				"dev/load-balancing-alb-nhn/package-config.yml": packageConfigOld,
			},
			packages: []common.Package{
				{
					Template:     "load-balancing-alb",
					Ref:          "load-balancing-alb-v5.2.0",
					OutputFolder: ".",
					VarFiles:     []string{"../common-config.yml", "package-config.yml"},
				},
			},
			expectedFiles: map[string]string{
				"dev/common-config.yml":                         commonConfig,
				"dev/load-balancing-alb-nhn/package-config.yml": packageConfigNew,
			},
		},
		{
			name: "reversed order does not write into the shared common config listed last",
			files: map[string]string{
				"dev/common-config.yml":                         commonConfig,
				"dev/load-balancing-alb-nhn/package-config.yml": packageConfigOld,
			},
			packages: []common.Package{
				{
					Template:     "load-balancing-alb",
					Ref:          "load-balancing-alb-v5.2.0",
					OutputFolder: ".",
					VarFiles:     []string{"package-config.yml", "../common-config.yml"},
				},
			},
			expectedFiles: map[string]string{
				"dev/common-config.yml":                         commonConfig,
				"dev/load-balancing-alb-nhn/package-config.yml": packageConfigNew,
			},
		},
		{
			name: "multiple var files declaring the same template: last match wins",
			files: map[string]string{
				"dev/load-balancing-alb-nhn/base-config.yml":    schemaLine("load-balancing-alb-v5.0.0") + "\nBase: true\n",
				"dev/load-balancing-alb-nhn/package-config.yml": packageConfigOld,
			},
			packages: []common.Package{
				{
					Template:     "load-balancing-alb",
					Ref:          "load-balancing-alb-v5.2.0",
					OutputFolder: ".",
					VarFiles:     []string{"base-config.yml", "package-config.yml"},
				},
			},
			expectedFiles: map[string]string{
				"dev/load-balancing-alb-nhn/base-config.yml":    schemaLine("load-balancing-alb-v5.0.0") + "\nBase: true\n",
				"dev/load-balancing-alb-nhn/package-config.yml": packageConfigNew,
			},
		},
		{
			name: "var file declaring a different template is never written to",
			files: map[string]string{
				"dev/common-config.yml":                         schemaLine("app-v1.0.0") + "\n" + commonConfig,
				"dev/load-balancing-alb-nhn/package-config.yml": "StackName: \"load-balancing-alb-nhn\"\n",
			},
			packages: []common.Package{
				{
					Template:     "load-balancing-alb",
					Ref:          "load-balancing-alb-v5.2.0",
					OutputFolder: ".",
					VarFiles:     []string{"package-config.yml", "../common-config.yml"},
				},
			},
			expectedFiles: map[string]string{
				"dev/common-config.yml":                         schemaLine("app-v1.0.0") + "\n" + commonConfig,
				"dev/load-balancing-alb-nhn/package-config.yml": "StackName: \"load-balancing-alb-nhn\"\n",
			},
			expectedWarning: "no var file of package '.' declares a JSON schema for template 'load-balancing-alb'",
		},
		{
			name: "no var file with a declaration: warn and skip instead of guessing by position",
			files: map[string]string{
				"dev/common-config.yml":                         commonConfig,
				"dev/load-balancing-alb-nhn/package-config.yml": "StackName: \"load-balancing-alb-nhn\"\n",
			},
			packages: []common.Package{
				{
					Template:     "load-balancing-alb",
					Ref:          "load-balancing-alb-v5.2.0",
					OutputFolder: ".",
					VarFiles:     []string{"package-config.yml", "../common-config.yml"},
				},
			},
			expectedFiles: map[string]string{
				"dev/common-config.yml":                         commonConfig,
				"dev/load-balancing-alb-nhn/package-config.yml": "StackName: \"load-balancing-alb-nhn\"\n",
			},
			expectedWarning: "no var file of package '.' declares a JSON schema for template 'load-balancing-alb'",
		},
		{
			name: "skipping one package does not abort the remaining packages",
			files: map[string]string{
				"dev/load-balancing-alb-nhn/app-config.yml":     "AppName: \"hello\"\n",
				"dev/load-balancing-alb-nhn/package-config.yml": packageConfigOld,
			},
			packages: []common.Package{
				{
					Template:     "app",
					Ref:          "app-v9.0.0",
					OutputFolder: "app-hello",
					VarFiles:     []string{"app-config.yml"},
				},
				{
					Template:     "load-balancing-alb",
					Ref:          "load-balancing-alb-v5.2.0",
					OutputFolder: ".",
					VarFiles:     []string{"package-config.yml"},
				},
			},
			expectedFiles: map[string]string{
				"dev/load-balancing-alb-nhn/app-config.yml":     "AppName: \"hello\"\n",
				"dev/load-balancing-alb-nhn/package-config.yml": packageConfigNew,
			},
			expectedWarning: "no var file of package 'app-hello' declares a JSON schema for template 'app'",
		},
		{
			name: "declaration already up to date is left unchanged",
			files: map[string]string{
				"dev/common-config.yml":                         commonConfig,
				"dev/load-balancing-alb-nhn/package-config.yml": packageConfigNew,
			},
			packages: []common.Package{
				{
					Template:     "load-balancing-alb",
					Ref:          "load-balancing-alb-v5.2.0",
					OutputFolder: ".",
					VarFiles:     []string{"../common-config.yml", "package-config.yml"},
				},
			},
			expectedFiles: map[string]string{
				"dev/common-config.yml":                         commonConfig,
				"dev/load-balancing-alb-nhn/package-config.yml": packageConfigNew,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootDir := t.TempDir()
			for relPath, content := range tt.files {
				fullPath := filepath.Join(rootDir, relPath)
				require.NoError(t, os.MkdirAll(filepath.Dir(fullPath), 0o755))
				require.NoError(t, os.WriteFile(fullPath, []byte(content), 0o644))
			}
			workingDirectory := filepath.Join(rootDir, "dev", "load-balancing-alb-nhn")

			var err error
			stderr := captureStderr(t, func() {
				err = Updater{}.setJsonSchemaDeclarationInVarFiles(tt.packages, workingDirectory)
			})
			require.NoError(t, err)

			if tt.expectedWarning == "" {
				require.NotContains(t, stderr, "Warning")
			} else {
				require.Contains(t, stderr, tt.expectedWarning)
			}

			for relPath, expectedContent := range tt.expectedFiles {
				content, err := os.ReadFile(filepath.Join(rootDir, relPath))
				require.NoError(t, err)
				require.Equal(t, expectedContent, string(content), "unexpected content in %s", relPath)
			}
		})
	}
}

func TestFindSchemaVarFile(t *testing.T) {
	t.Run("package without var files", func(t *testing.T) {
		_, _, found := findSchemaVarFile(common.Package{Template: "app", VarFiles: []string{}}, t.TempDir())
		require.False(t, found)
	})

	t.Run("missing var file is skipped without error", func(t *testing.T) {
		pkg := common.Package{Template: "app", VarFiles: []string{"does-not-exist.yml"}}
		_, _, found := findSchemaVarFile(pkg, t.TempDir())
		require.False(t, found)
	})

	t.Run("matching declaration is found and parsed", func(t *testing.T) {
		workingDirectory := t.TempDir()
		varFile := filepath.Join(workingDirectory, "package-config.yml")
		require.NoError(t, os.WriteFile(varFile, []byte(schemaLine("app-v1.2.3")+"\nAppName: \"hello\"\n"), 0o644))

		pkg := common.Package{Template: "app", VarFiles: []string{"package-config.yml"}}
		selected, existingSchema, found := findSchemaVarFile(pkg, workingDirectory)
		require.True(t, found)
		require.Equal(t, varFile, selected)
		require.Equal(t, "app-v1.2.3", existingSchema.Ref())
	})
}
