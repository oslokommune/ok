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

type UpdateTestData struct {
	TestData

	releases      map[string]string
	expectError   bool
	expectedFiles []string
}

func TestUpdateCommand(t *testing.T) {
	tests := []UpdateTestData{
		{
			TestData: TestData{
				name:            "Should bump the Ref field for the specified packages",
				args:            []string{"app-hello", "load-balancing-alb-main"},
				testdataRootDir: "testdata/update/bump-ref-field",
			},
			releases: map[string]string{
				"app":                "v9.0.0",
				"load-balancing-alb": "v4.0.0",
				"app-common":         "v7.0.0",
			},
			expectError: false,
			expectedFiles: []string{
				"packages.yml",
				"config/app-hello.yml",
				"config/common-config.yml",
			},
		},
		{
			TestData: TestData{
				name:            "Should bump the Ref field only for semver-version package Refs",
				args:            []string{},
				testdataRootDir: "testdata/update/bump-ref-field-semver-only",
			},
			releases: map[string]string{
				"app": "v9.0.0",
			},
			expectError: false,
			expectedFiles: []string{
				"packages.yml",
			},
		},
		{
			TestData: TestData{
				name:            "Should bump schema version in var files",
				args:            []string{"app-hello"},
				testdataRootDir: "testdata/update/bump-schema-version",
			},
			releases: map[string]string{
				"app": "v9.0.1",
			},
			expectError: false,
			expectedFiles: []string{
				"packages.yml",
				"config/app-hello.yml",
				"common-config.yml",
			},
		},
		{
			TestData: TestData{
				name:            "Should migrate schema declaration from dir based to HTTPS based",
				args:            []string{"app-hello"},
				testdataRootDir: "testdata/update/migrate-schema-declaration-format",
			},
			releases: map[string]string{
				"app": "v9.0.1",
			},
			expectError: false,
			expectedFiles: []string{
				"packages.yml",
				"config/app-hello.yml",
				"common-config.yml",
			},
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

			command := pkg.NewUpdateCommand(ghReleases)

			tempDir, err := os.MkdirTemp(os.TempDir(), "ok-"+tt.name)

			// Remove temp dir after test run
			//defer func(path string) {
			//	err := os.RemoveAll(path)
			//	require.NoError(t, err)
			//}(tempDir)

			require.NoError(t, err)

			fmt.Println("tempDir: ", tempDir)
			copyTestdataRootDirToTempDir(t, tt.TestData, tempDir)
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
			for _, file := range tt.expectedFiles {
				actualBytes, err := os.ReadFile(filepath.Join(tempDir, file))
				require.NoError(t, err)
				actual := string(actualBytes)

				expectedBytes, err := os.ReadFile(filepath.Join(tt.testdataRootDir, "expected", file))
				require.NoError(t, err)
				expected := string(expectedBytes)

				assert.Equal(t, expected, actual)
			}
		})
	}
}

type GitHubReleasesMock struct {
	LatestReleases map[string]string
}

func (g *GitHubReleasesMock) GetLatestReleases() (map[string]string, error) {
	return g.LatestReleases, nil
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
