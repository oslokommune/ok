package schema

import (
	"encoding/json"
	"fmt"
	"github.com/oslokommune/ok/pkg/jsonschema"
	"log/slog"
	"os"
	"path/filepath"
)

// writeJsonSchemaFile writes a json schema to a file in the output directory with the given template and version.
// The file will be named <template>-<version>.schema.json.
// The return value is the path to the file.
func writeJsonSchemaFile(filePath string, schema *jsonschema.Document) (string, error) {
	outputDir := filepath.Dir(filePath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("creating output directory: %w", err)
	}

	slog.Debug("writing json schema file", slog.String("path", filePath), slog.String("schemaId", schema.ID))

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("opening file: %w", err)
	}

	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(schema); err != nil {
		return "", fmt.Errorf("encoding schema to file %s: %w", filePath, err)
	}

	return filePath, nil
}
