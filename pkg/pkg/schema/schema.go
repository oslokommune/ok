package schema

import (
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/config"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/oslokommune/ok/pkg/jsonschema"
)

func CreateOrUpdateConfigurationFile(configFilePath string, schemaName string, schema *jsonschema.Document) (string, error) {
	configDir := filepath.Dir(configFilePath)
	schemasDir := filepath.Join(configDir, ".schemas")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("creating output directory: %w", err)
	}
	if err := os.MkdirAll(schemasDir, 0755); err != nil {
		return "", fmt.Errorf("creating output directory: %w", err)
	}
	// write schema to the config dir with the name <schemaName>.schema.json
	schemaFileName := fmt.Sprintf("%s.schema.json", schemaName)
	schemaFilePath := filepath.Join(schemasDir, schemaFileName)
	slog.Debug("writing schema file", slog.String("path", schemaFilePath), slog.String("configDir", configDir), slog.String("schemaDir", schemasDir), slog.String("schemaId", schema.ID))
	_, err := writeJsonSchemaFile(schemaFilePath, schema)
	if err != nil {
		return "", fmt.Errorf("writing schema file to %s: %w", schemaFilePath, err)
	}

	relativeSchemaPath, err := filepath.Rel(configDir, schemaFilePath)
	if err != nil {
		return "", fmt.Errorf("getting relative schema path: %w", err)
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("reading file: %w", err)
	}
	// Find if the first line starts with # yaml-language-server
	cleanedConfig := stripYamlLanguageServerComment(string(data))
	newConfig := appendYamlLanguageServerComment(cleanedConfig, relativeSchemaPath)
	err = os.WriteFile(configFilePath, []byte(newConfig), 0644)
	if err != nil {
		return "", fmt.Errorf("overwriting config file %s: %w", configFilePath, err)
	}

	return "", nil
}

const yamlLanguageServerComment = "# yaml-language-server:"

func stripYamlLanguageServerComment(data string) string {
	first, rest, ok := strings.Cut(data, "\n")
	if ok && strings.HasPrefix(first, yamlLanguageServerComment) {
		return rest
	}
	return data
}

func appendYamlLanguageServerComment(data, schemaPath string) string {
	return fmt.Sprintf("%s $schema=%s\n%s", yamlLanguageServerComment, schemaPath, data)
}

func BuildJsonSchemaFromConfig(config *config.BoilerplateConfig, dependencies []config.BoilerplateConfig) (*jsonschema.Document, error) {
	return nil, fmt.Errorf("not implemented")
}
