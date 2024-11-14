package metadata

import (
	"fmt"
	"github.com/Masterminds/semver"
	"regexp"
	"strings"
)

type JsonSchema struct {
	Template string
	Version  *semver.Version
}

func ParseMetadata(firstLine string) (JsonSchema, error) {
	if !strings.HasPrefix(firstLine, "# yaml-language-server: $schema=") {
		return JsonSchema{}, fmt.Errorf("missing schema declaration")
	}

	template, version, err := parseSchemaLine(firstLine)
	if err != nil {
		return JsonSchema{}, fmt.Errorf("parsing schema line: %w", err)
	}

	versionSemver, err := semver.NewVersion(version)
	if err != nil {
		return JsonSchema{}, fmt.Errorf("creating semver: %w", err)
	}

	return JsonSchema{
		Template: template,
		Version:  versionSemver,
	}, nil
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
