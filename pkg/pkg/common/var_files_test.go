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
		name          string
		files         []string
		varFiles      []string
		noRepo        bool
		absVarFiles   []string
		expectedError string
	}{
		{
			name:     "should return no error when all var files exist",
			files:    []string{"stacks/prod/common-config.yml", "stacks/prod/apps/my-app/package-config.yml"},
			varFiles: []string{"../../common-config.yml", "package-config.yml"},
		},
		{
			name:        "should return no error when an absolute var file exists",
			files:       []string{"stacks/prod/common-config.yml"},
			absVarFiles: []string{"stacks/prod/common-config.yml"},
		},
		{
			name:          "should suggest a var file in a parent directory",
			files:         []string{"stacks/prod/common-config.yml"},
			varFiles:      []string{"../common-config.yml"},
			expectedError: `stacks/prod/apps/my-app/packages.yml: var file "../common-config.yml" does not exist. Did you mean "../../common-config.yml"?`,
		},
		{
			name:          "should suggest the nearest var file",
			files:         []string{"common-config.yml", "stacks/common-config.yml"},
			varFiles:      []string{"common-config.yml"},
			expectedError: `stacks/prod/apps/my-app/packages.yml: var file "common-config.yml" does not exist. Did you mean "../../../common-config.yml"?`,
		},
		{
			name:          "should suggest a var file in the working directory",
			files:         []string{"stacks/prod/apps/my-app/common-config.yml"},
			varFiles:      []string{"../../common-config.yml"},
			expectedError: `stacks/prod/apps/my-app/packages.yml: var file "../../common-config.yml" does not exist. Did you mean "common-config.yml"?`,
		},
		{
			name:          "should keep the subdirectory of the var file",
			files:         []string{"stacks/prod/_config/app.yml", "stacks/prod/apps/app.yml"},
			varFiles:      []string{"_config/app.yml"},
			expectedError: `stacks/prod/apps/my-app/packages.yml: var file "_config/app.yml" does not exist. Did you mean "../../_config/app.yml"?`,
		},
		{
			name:          "should not suggest a var file when no file with the same name exists",
			varFiles:      []string{"../common-config.yml"},
			expectedError: `stacks/prod/apps/my-app/packages.yml: var file "../common-config.yml" does not exist`,
		},
		{
			name:          "should not suggest a var file outside the repository",
			files:         []string{"../outside.yml"},
			varFiles:      []string{"outside.yml"},
			expectedError: `stacks/prod/apps/my-app/packages.yml: var file "outside.yml" does not exist`,
		},
		{
			name:          "should not suggest a var file when not in a repository",
			files:         []string{"stacks/prod/common-config.yml"},
			varFiles:      []string{"../common-config.yml"},
			noRepo:        true,
			expectedError: `packages.yml: var file "../common-config.yml" does not exist`,
		},
		{
			name:     "should report each missing var file once",
			files:    []string{"stacks/prod/common-config.yml"},
			varFiles: []string{"../common-config.yml", "../common-config.yml", "missing.yml"},
			expectedError: `stacks/prod/apps/my-app/packages.yml: var file "../common-config.yml" does not exist. Did you mean "../../common-config.yml"?` + "\n" +
				`stacks/prod/apps/my-app/packages.yml: var file "missing.yml" does not exist`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			tempDir, err := resolvePath(t.TempDir())
			require.NoError(t, err)

			repoRoot := filepath.Join(tempDir, "repo")
			workingDirectory := filepath.Join(repoRoot, "stacks", "prod", "apps", "my-app")
			require.NoError(t, os.MkdirAll(workingDirectory, 0755))

			for _, file := range tt.files {
				path := filepath.Join(repoRoot, file)
				require.NoError(t, os.MkdirAll(filepath.Dir(path), 0755))
				require.NoError(t, os.WriteFile(path, []byte(""), 0644))
			}

			manifestFile := filepath.Join(workingDirectory, PackagesManifestFilename)
			require.NoError(t, os.WriteFile(manifestFile, []byte(""), 0644))

			if tt.noRepo {
				repoRoot = ""
				manifestFile = PackagesManifestFilename
			}

			varFiles := tt.varFiles
			for _, file := range tt.absVarFiles {
				varFiles = append(varFiles, filepath.Join(repoRoot, file))
			}

			packages := []Package{{OutputFolder: "out", Template: "app", Ref: "app-v1.0.0", VarFiles: varFiles}}

			// When
			err = checkVarFiles(manifestFile, packages, workingDirectory, repoRoot)

			// Then
			if tt.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func TestGitRepoRoot(t *testing.T) {
	tempDir, err := resolvePath(t.TempDir())
	require.NoError(t, err)

	subDir := filepath.Join(tempDir, "stacks", "prod")
	require.NoError(t, os.MkdirAll(subDir, 0755))
	require.NoError(t, exec.Command("git", "init", "--quiet", tempDir).Run())

	require.Equal(t, tempDir, gitRepoRoot(subDir))
}
