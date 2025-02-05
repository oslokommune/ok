package pkg_test

import (
	"context"
	"fmt"
	"github.com/oslokommune/ok/cmd/pkg"
	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestUpdateCommand(t *testing.T) {
	tests := []TestData{
		{
			name:            "Should bump the Ref field for the specified packages",
			args:            []string{"app-hello", "load-balancing-alb-main"},
			jsonSchemasDir:  "testdata/bump-ref-field/input/json-schemas",
			packageManifest: "testdata/bump-ref-field/input/packages.yml",
			configDir:       "testdata/bump-ref-field/input/config",
			releases: map[string]string{
				"app":                "v9.0.0",
				"load-balancing-alb": "v4.0.0",
				"app-common":         "v7.0.0",
			},
			expectError:             false,
			expectedPackageManifest: "testdata/bump-ref-field/expected/packages.yml",
			expectedConfigDir:       "testdata/bump-ref-field/expected/config",
		},
		{
			name:            "Should bump the Ref field only for semver-version package Refs",
			args:            []string{},
			jsonSchemasDir:  "testdata/bump-ref-field-semver-only/input/json-schemas",
			packageManifest: "testdata/bump-ref-field-semver-only/input/packages.yml",
			configDir:       "testdata/bump-ref-field-semver-only/input/config",
			releases: map[string]string{
				"app": "v9.0.0",
			},
			expectError:             false,
			expectedPackageManifest: "testdata/bump-ref-field-semver-only/expected/packages.yml",
		},
		{
			name:            "Should bump schema version in var files",
			args:            []string{"app-hello"},
			jsonSchemasDir:  "testdata/bump-schema-version/input/json-schemas",
			packageManifest: "testdata/bump-schema-version/input/packages.yml",
			configDir:       "testdata/bump-schema-version/input/config",
			releases: map[string]string{
				"app": "v9.0.0",
			},
			expectError:             false,
			expectedPackageManifest: "testdata/bump-schema-version/expected/packages.yml",
			expectedConfigDir:       "testdata/bump-schema-version/expected/config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			testDir, err := os.Getwd()
			require.NoError(t, err)

			ghReleases := &GitHubReleasesMock{
				LatestReleases: tt.releases,
			}
			schemaGenerator := NewSchemaGeneratorMock(filepath.Join(testDir, tt.jsonSchemasDir))

			command := pkg.NewUpdateCommand(ghReleases, schemaGenerator)

			tempDir, err := os.MkdirTemp(os.TempDir(), "ok-"+tt.name)
			defer func(path string) {
				err := os.RemoveAll(path)
				require.NoError(t, err)
			}(tempDir)

			require.NoError(t, err)

			fmt.Println("tempDir: ", tempDir)
			copyTestdataToTempDir(t, tt, tempDir)
			command.SetArgs(tt.args)

			err = os.Chdir(tempDir) // Works, but disables the possibility for parallel tests.
			require.NoError(t, err)

			// When
			err = command.Execute()

			// Then
			if tt.expectError {
				assert.Error(t, err, "expected an error")
				return
			}
			require.NoError(t, err)

			err = os.Chdir(testDir)
			require.NoError(t, err)

			// Compare package manifest file
			actualBytes, err := os.ReadFile(filepath.Join(tempDir, "packages.yml"))
			require.NoError(t, err)
			actual := string(actualBytes)

			expectedBytes, err := os.ReadFile(tt.expectedPackageManifest)
			require.NoError(t, err)
			expected := string(expectedBytes)

			assert.Equal(t, expected, actual)

			// Compare var files:
			// Given
			// testadat/some-test/expected/config/app-hello.yml, we want to compare it to
			//                     tempDir/config/app-hello.yml
			if tt.expectedConfigDir == "" {
				return
			}

			err = filepath.Walk(tt.expectedConfigDir, func(path string, fileInfo os.FileInfo, err error) error {
				require.NoError(t, err)

				if fileInfo.IsDir() {
					return nil
				}

				actualFilename := filepath.Join(tempDir, "config", fileInfo.Name())
				actualVarFile := actualFilename

				varFileBytes, err := os.ReadFile(actualVarFile)
				require.NoError(t, err)
				varFile := string(varFileBytes)

				expectedFilename := filepath.Join(tt.expectedConfigDir, fileInfo.Name())
				expectedVarFileBytes, err := os.ReadFile(expectedFilename)
				require.NoError(t, err)
				expectedVarFile := string(expectedVarFileBytes)

				assert.Equal(t, expectedVarFile, varFile)

				return nil
			})
		})
	}
}

type TestData struct {
	name                    string
	args                    []string
	jsonSchemasDir          string
	packageManifest         string
	configDir               string
	releases                map[string]string
	expectError             bool
	expectedPackageManifest string
	expectedConfigDir       string
}

type GitHubReleasesMock struct {
	LatestReleases map[string]string
}

func (g *GitHubReleasesMock) GetLatestReleases() (map[string]string, error) {
	return g.LatestReleases, nil
}

func copyTestdataToTempDir(t *testing.T, tt TestData, rootDir string) {
	if tt.packageManifest != "" {
		srcPath := tt.packageManifest
		dstPath := filepath.Join(rootDir, filepath.Base(tt.packageManifest))

		err := copyFile(srcPath, dstPath)
		require.NoError(t, err)
	}

	if tt.configDir != "" {
		srcDir := os.DirFS(tt.configDir)
		dstDir := filepath.Join(rootDir, filepath.Base(tt.configDir))

		err := os.CopyFS(dstDir, srcDir)
		require.NoError(t, err)
	}
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	err = os.WriteFile(dst, input, 0644)
	if err != nil {
		return err
	}

	return nil
}

type SchemaGeneratorMock struct {
	jsonSchemasDir string
}

func NewSchemaGeneratorMock(jsonSchemasDir string) SchemaGeneratorMock {
	return SchemaGeneratorMock{
		jsonSchemasDir: jsonSchemasDir,
	}
}

// CreateJsonSchemaFile emulates creating JSON schema file from a Boilerplate template configuration.
// Instead of generating the schema, it just copies a pre-generated file.
func (s SchemaGeneratorMock) CreateJsonSchemaFile(
	ctx context.Context, manifestPackagePrefix string, pkg common.Package) ([]byte, error) {

	// Example: app-v8.0.5.schema.json
	schemaFilePath := fmt.Sprintf("%s.schema.json", pkg.Ref)

	// Example: testdata/json-schemas/.schemas/app-v8.0.5.schema.json
	filePath := filepath.Join(s.jsonSchemasDir, schemaFilePath)

	schemaData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", filePath, err)
	}

	return schemaData, nil
}
