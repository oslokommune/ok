package metadata

import (
	"fmt"
	"regexp"
	"strings"
)

func ParseMetadata(firstLine string) (VarFileMetadata, error) {
	varFileMetadata := VarFileMetadata{
		HasVersion: false,
		Template:   "",
		Version:    "",
	}

	if strings.HasPrefix(firstLine, "# yaml-language-server: $schema=") {
		template, version, err := parseSchemaLine(firstLine)
		if err != nil {
			return VarFileMetadata{}, fmt.Errorf("parsing schema line: %w", err)
		}

		varFileMetadata.HasVersion = true
		varFileMetadata.Template = template
		varFileMetadata.Version = version
	}

	return varFileMetadata, nil
}

func parseSchemaLine(line string) (string, string, error) {
	re := regexp.MustCompile(`\.schemas/(\w+)-(\S+)\.schema\.json`)

	matches := re.FindStringSubmatch(line)
	if len(matches) != 3 {
		return "", "", fmt.Errorf("invalid schema format")
	}

	template, version := matches[1], matches[2]
	return template, version, nil
}
