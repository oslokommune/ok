package common

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckVarFiles(t *testing.T) {
	tests := []struct {
		name string
		// manifest is the path of the package manifest, relative to the repository root.
		manifest string
		varFiles []string
		// files are the files that exist, relative to the repository root.
		files         []string
		noGitRepo     bool
		expectedError string
	}{
		// Package manifest in the stack directory
		{
			name:     "stack directory: should return no error when all var files exist",
			manifest: "stacks/prod/apps/my-app/packages.yml",
			varFiles: []string{"../../common-config.yml", "package-config.yml"},
			files:    []string{"stacks/prod/common-config.yml", "stacks/prod/apps/my-app/package-config.yml"},
		},
		{
			name:          "stack directory: should suggest a var file in a parent directory",
			manifest:      "stacks/prod/apps/my-app/packages.yml",
			varFiles:      []string{"../common-config.yml"},
			files:         []string{"stacks/prod/common-config.yml"},
			expectedError: `stacks/prod/apps/my-app/packages.yml: var file "../common-config.yml" does not exist. Did you mean "../../common-config.yml"?`,
		},
		{
			name:          "stack directory: should suggest the nearest var file, starting in the stack directory",
			manifest:      "stacks/prod/apps/my-app/packages.yml",
			varFiles:      []string{"../../common-config.yml"},
			files:         []string{"stacks/prod/apps/my-app/common-config.yml", "common-config.yml"},
			expectedError: `stacks/prod/apps/my-app/packages.yml: var file "../../common-config.yml" does not exist. Did you mean "common-config.yml"?`,
		},

		// Package manifest for GitHub Actions workflows
		{
			name:     "github actions: should return no error when all var files exist",
			manifest: ".github/workflows/_config/dev/packages.yml",
			varFiles: []string{"common-config.yml", "docker-build-push.yml"},
			files:    []string{".github/workflows/_config/dev/common-config.yml", ".github/workflows/_config/dev/docker-build-push.yml"},
		},
		{
			name:          "github actions: should suggest a var file in a parent directory",
			manifest:      ".github/workflows/_config/dev/packages.yml",
			varFiles:      []string{"common-config.yml"},
			files:         []string{".github/workflows/_config/common-config.yml"},
			expectedError: `.github/workflows/_config/dev/packages.yml: var file "common-config.yml" does not exist. Did you mean "../common-config.yml"?`,
		},

		// Central package manifest with var files in a _config directory
		{
			name:     "central manifest: should return no error when all var files exist",
			manifest: "stacks/dev/packages.yml",
			varFiles: []string{"_config/common-config.yml", "_config/databases.yml"},
			files:    []string{"stacks/dev/_config/common-config.yml", "stacks/dev/_config/databases.yml"},
		},
		{
			name:          "central manifest: should keep the subdirectory of the var file in the suggestion",
			manifest:      "stacks/dev/packages.yml",
			varFiles:      []string{"_config/common-config.yml"},
			files:         []string{"stacks/_config/common-config.yml", "stacks/dev/common-config.yml"},
			expectedError: `stacks/dev/packages.yml: var file "_config/common-config.yml" does not exist. Did you mean "../_config/common-config.yml"?`,
		},

		// Cases without a suggestion
		{
			name:          "should not suggest a var file when no file with the same name exists",
			manifest:      "stacks/prod/apps/my-app/packages.yml",
			varFiles:      []string{"../common-config.yml"},
			expectedError: `stacks/prod/apps/my-app/packages.yml: var file "../common-config.yml" does not exist`,
		},
		{
			name:          "should not suggest a var file outside the repository",
			manifest:      "stacks/prod/apps/my-app/packages.yml",
			varFiles:      []string{"common-config.yml"},
			files:         []string{"../common-config.yml"},
			expectedError: `stacks/prod/apps/my-app/packages.yml: var file "common-config.yml" does not exist`,
		},
		{
			name:          "should not suggest a var file when not in a git repository",
			manifest:      "stacks/prod/apps/my-app/packages.yml",
			varFiles:      []string{"../common-config.yml"},
			files:         []string{"stacks/prod/common-config.yml"},
			noGitRepo:     true,
			expectedError: `stacks/prod/apps/my-app/packages.yml: var file "../common-config.yml" does not exist`,
		},

		// Several missing var files
		{
			name:     "should report each missing var file one time",
			manifest: "stacks/prod/apps/my-app/packages.yml",
			varFiles: []string{"../common-config.yml", "../common-config.yml", "package-config.yml"},
			files:    []string{"stacks/prod/common-config.yml"},
			expectedError: `stacks/prod/apps/my-app/packages.yml: var file "../common-config.yml" does not exist. Did you mean "../../common-config.yml"?` + "\n" +
				`stacks/prod/apps/my-app/packages.yml: var file "package-config.yml" does not exist`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			repoRoot := createRepo(t, !tt.noGitRepo)
			createFiles(t, repoRoot, append(tt.files, tt.manifest))

			packages := []Package{{OutputFolder: ".", Template: "app", Ref: "app-v1.0.0", VarFiles: tt.varFiles}}

			// When
			err := CheckVarFiles(tt.manifest, packages, filepath.Dir(tt.manifest))

			// Then
			if tt.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func TestCheckVarFilesWithAbsolutePath(t *testing.T) {
	// Given
	repoRoot := createRepo(t, true)
	manifest := "stacks/prod/apps/my-app/packages.yml"
	createFiles(t, repoRoot, []string{manifest, "stacks/prod/common-config.yml"})

	existing := filepath.Join(repoRoot, "stacks/prod/common-config.yml")
	missing := filepath.Join(repoRoot, "stacks/prod/apps/common-config.yml")
	packages := []Package{{OutputFolder: ".", Template: "app", Ref: "app-v1.0.0", VarFiles: []string{existing, missing}}}

	// When
	err := CheckVarFiles(manifest, packages, filepath.Dir(manifest))

	// Then
	require.EqualError(t, err, manifest+`: var file "`+missing+`" does not exist`)
}

// createRepo creates a repository root directory, makes it the current directory, and returns its path. Git does not
// look for a repository above the directory, so the result does not depend on where the test runs.
func createRepo(t *testing.T, gitRepo bool) string {
	tempDir, err := resolvePath(t.TempDir())
	require.NoError(t, err)

	repoRoot := filepath.Join(tempDir, "repo")
	require.NoError(t, os.MkdirAll(repoRoot, 0755))

	t.Setenv("GIT_CEILING_DIRECTORIES", tempDir)
	t.Chdir(repoRoot)

	if gitRepo {
		require.NoError(t, exec.Command("git", "init", "--quiet").Run())
	}

	return repoRoot
}

func createFiles(t *testing.T, dir string, files []string) {
	for _, file := range files {
		path := filepath.Join(dir, file)
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0755))
		require.NoError(t, os.WriteFile(path, []byte(""), 0644))
	}
}
