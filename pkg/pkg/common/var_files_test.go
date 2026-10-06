package common

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckVarFilesExist(t *testing.T) {
	tests := []struct {
		name string
		// manifestPath is the path of the package manifest, relative to the repository root.
		manifestPath string
		packages     []Package
		// files are the files that exist in addition to the package manifest, relative to the repository root.
		files         []string
		noGitRepo     bool
		expectedError string
	}{
		{
			name:         "should return no error when all var files exist",
			manifestPath: "stacks/dev/packages.yml",
			packages: []Package{
				{OutputFolder: "databases", Template: "databases", Ref: "databases-v4.0.0", VarFiles: []string{"_config/common-config.yml", "_config/databases.yml"}},
				{OutputFolder: "app-hello", Template: "app", Ref: "app-v6.1.1", VarFiles: []string{"_config/common-config.yml", "_config/app-hello.yml"}},
			},
			files: []string{"stacks/dev/_config/common-config.yml", "stacks/dev/_config/databases.yml", "stacks/dev/_config/app-hello.yml"},
		},
		{
			name:         "should suggest a var file in a parent directory",
			manifestPath: "stacks/prod/apps/my-app/packages.yml",
			packages: []Package{
				{OutputFolder: ".", Template: "app", Ref: "app-v1.0.0", VarFiles: []string{"../common-config.yml", "package-config.yml"}},
			},
			files:         []string{"stacks/prod/common-config.yml", "stacks/prod/apps/my-app/package-config.yml"},
			expectedError: `stacks/prod/apps/my-app/packages.yml: var file "../common-config.yml" does not exist. Did you mean "../../common-config.yml"?`,
		},
		{
			name:         "should suggest the nearest var file, starting in the directory of the package manifest",
			manifestPath: "stacks/prod/apps/my-app/packages.yml",
			packages: []Package{
				{OutputFolder: ".", Template: "app", Ref: "app-v1.0.0", VarFiles: []string{"../../common-config.yml"}},
			},
			files:         []string{"stacks/prod/apps/my-app/common-config.yml", "common-config.yml"},
			expectedError: `stacks/prod/apps/my-app/packages.yml: var file "../../common-config.yml" does not exist. Did you mean "common-config.yml"?`,
		},
		{
			name:         "should report each missing var file one time, and keep its subdirectory in the suggestion",
			manifestPath: "stacks/dev/packages.yml",
			packages: []Package{
				{OutputFolder: "databases", Template: "databases", Ref: "databases-v4.0.0", VarFiles: []string{"_config/common-config.yml", "_config/databases.yml"}},
				{OutputFolder: "app-hello", Template: "app", Ref: "app-v6.1.1", VarFiles: []string{"_config/common-config.yml", "_config/app-hello.yml"}},
			},
			files: []string{"stacks/_config/common-config.yml", "stacks/dev/common-config.yml", "stacks/dev/_config/databases.yml"},
			expectedError: `stacks/dev/packages.yml: var file "_config/common-config.yml" does not exist. Did you mean "../_config/common-config.yml"?` + "\n" +
				`stacks/dev/packages.yml: var file "_config/app-hello.yml" does not exist`,
		},
		{
			name:         "should stop the search at the repository root",
			manifestPath: "packages.yml",
			packages: []Package{
				{OutputFolder: ".", Template: "app", Ref: "app-v1.0.0", VarFiles: []string{"common-config.yml"}},
			},
			// The only common-config.yml is in the directory above the repository root.
			files:         []string{"../common-config.yml"},
			expectedError: `packages.yml: var file "common-config.yml" does not exist`,
		},
		{
			name:         "should not suggest a var file when not in a git repository",
			manifestPath: "stacks/prod/apps/my-app/packages.yml",
			packages: []Package{
				{OutputFolder: ".", Template: "app", Ref: "app-v1.0.0", VarFiles: []string{"../common-config.yml"}},
			},
			files:     []string{"stacks/prod/common-config.yml"},
			noGitRepo: true,
			// Without a git repository, the path of the package manifest shows as given.
			expectedError: `packages.yml: var file "../common-config.yml" does not exist`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			repoRoot := createRepo(t, !tt.noGitRepo)
			require.NoError(t, SavePackageManifest(filepath.Join(repoRoot, tt.manifestPath), PackageManifest{Packages: tt.packages}))
			for _, file := range tt.files {
				createFile(t, filepath.Join(repoRoot, file), "")
			}

			// Run from the directory of the package manifest, as "ok pkg install" and "ok pkg update" do.
			t.Chdir(filepath.Join(repoRoot, filepath.Dir(tt.manifestPath)))
			manifestFile := filepath.Base(tt.manifestPath)

			// When
			err := CheckVarFilesExist(manifestFile, tt.packages, ".")

			// Then
			if tt.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func TestCheckVarFilesExistWithAbsolutePath(t *testing.T) {
	// Given
	repoRoot := createRepo(t, true)
	existing := filepath.Join(repoRoot, "stacks/prod/common-config.yml")
	missing := filepath.Join(repoRoot, "stacks/prod/apps/common-config.yml")
	manifestPath := "stacks/prod/apps/my-app/packages.yml"

	packages := []Package{
		{OutputFolder: ".", Template: "app", Ref: "app-v1.0.0", VarFiles: []string{existing, missing}},
	}

	createFile(t, existing, "")
	require.NoError(t, SavePackageManifest(filepath.Join(repoRoot, manifestPath), PackageManifest{Packages: packages}))

	t.Chdir(filepath.Join(repoRoot, filepath.Dir(manifestPath)))

	// When
	err := CheckVarFilesExist("packages.yml", packages, ".")

	// Then
	require.EqualError(t, err, fmt.Sprintf(`%s: var file %q does not exist`, manifestPath, missing))
}

// TestCheckVarFilesExistFromRepositoryRoot runs the check from the repository root, as "ok pkg install --recursive" and
// "ok pkg update --recursive" do.
func TestCheckVarFilesExistFromRepositoryRoot(t *testing.T) {
	// Given
	repoRoot := createRepo(t, true)
	manifestPath := "stacks/prod/apps/my-app/packages.yml"

	packages := []Package{
		{OutputFolder: ".", Template: "app", Ref: "app-v1.0.0", VarFiles: []string{"../common-config.yml"}},
	}

	createFile(t, filepath.Join(repoRoot, "stacks/prod/common-config.yml"), "")
	require.NoError(t, SavePackageManifest(filepath.Join(repoRoot, manifestPath), PackageManifest{Packages: packages}))

	// When
	err := CheckVarFilesExist(manifestPath, packages, filepath.Dir(manifestPath))

	// Then
	require.EqualError(t, err, `stacks/prod/apps/my-app/packages.yml: var file "../common-config.yml" does not exist. Did you mean "../../common-config.yml"?`)
}

// createRepo creates a repository root directory, makes it the current directory, and returns its path. Git does not
// look for a repository above the directory, so the result does not depend on where the test runs.
func createRepo(t *testing.T, gitRepo bool) string {
	tempDir := t.TempDir()

	repoRoot := filepath.Join(tempDir, "repo")
	require.NoError(t, os.MkdirAll(repoRoot, 0755))

	t.Setenv("GIT_CEILING_DIRECTORIES", tempDir)
	t.Chdir(repoRoot)

	if gitRepo {
		require.NoError(t, exec.Command("git", "init", "--quiet").Run())
	}

	return repoRoot
}

func createFile(t *testing.T, path string, content string) {
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0755))
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))
}
