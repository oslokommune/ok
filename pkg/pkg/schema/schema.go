package schema

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const yamlLanguageServerComment = "# yaml-language-server:"

func getVarFileDir(varfilePath string) string {
	return filepath.Dir(varfilePath)
}

// SetSchemaDeclarationInVarFile sets the first line of a var file to include a $schema reference to the schema file.
// schemaName should be a ok package Ref, for instance "app-v9.0.0"
func SetSchemaDeclarationInVarFile(varfilePath string, schemaName string) error {
	varFileDir := getVarFileDir(varfilePath)

	if err := os.MkdirAll(varFileDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	schemaUri := fmt.Sprintf(
		"https://raw.githubusercontent.com/oslokommune/golden-path-boilerplate-schemas/refs/heads/main/schemas/%s.schema.json",
		schemaName)

	varFileData, err := os.ReadFile(varfilePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading file: %w", err)
	}

	varFileWithoutJsonSchemaDeclaration := stripYamlLanguageServerComment(string(varFileData))
	updatedVarFile := fmt.Sprintf("%s $schema=%s\n%s", yamlLanguageServerComment, schemaUri, varFileWithoutJsonSchemaDeclaration)

	err = os.WriteFile(varfilePath, []byte(updatedVarFile), 0644)
	if err != nil {
		return fmt.Errorf("overwriting config file %s: %w", varfilePath, err)
	}

	return nil
}

func stripYamlLanguageServerComment(varFileData string) string {
	first, rest, ok := strings.Cut(varFileData, "\n")
	if ok && strings.HasPrefix(first, yamlLanguageServerComment) {
		return rest
	}
	return varFileData
}
