package metadata

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/Masterminds/semver"
	"os"
	"regexp"
	"strings"
)

var (
	ErrMissingSchemaDeclaration = errors.New("missing schema declaration")
)

type JsonSchema struct {
	// Template is a Boilerplate template, and can be for instance "app" or "networking"
	Template string

	// Version is the semantically versioned schema version, for instance 'v1.5.2'
	Version *semver.Version
}

func (j JsonSchema) Ref() string {
	return fmt.Sprintf("%s-v%s", j.Template, j.Version.String())
}

func ParseFirstLine(varFile string) (JsonSchema, error) {
	firstLine, err := ReadFirstLine(varFile)
	if err != nil {
		return JsonSchema{}, fmt.Errorf("reading first line from %s: %w", varFile, err)
	}

	varFileMetadata, err := parseMetadata(firstLine)
	if err != nil {
		return JsonSchema{}, fmt.Errorf("parsing metadata from line '%s': %w", firstLine, err)
	}

	return varFileMetadata, nil
}

func ReadFirstLine(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}
		return "", fmt.Errorf("file is empty")
	}

	return scanner.Text(), nil
}

func parseMetadata(firstLine string) (JsonSchema, error) {
	if !strings.HasPrefix(firstLine, "# yaml-language-server: $schema=") {
		return JsonSchema{}, ErrMissingSchemaDeclaration
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
	re := regexp.MustCompile(`\$schema=.*\/([\w-]+)-(\S+)\.schema\.json`)

	matches := re.FindStringSubmatch(line)
	if len(matches) != 3 {
		return "", "", fmt.Errorf("invalid schema format")
	}

	template, version := matches[1], matches[2]
	return template, version, nil
}
