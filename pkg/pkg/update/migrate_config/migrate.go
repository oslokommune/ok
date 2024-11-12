package migrate_config

import (
	"bufio"
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/update/migrate_config/add_apex_domain"
	"github.com/oslokommune/ok/pkg/pkg/update/migrate_config/metadata"
	"os"
	"regexp"
	"strings"
)

func UpdatePackageConfig(packagesToUpdate []common.Package) error {
	for _, pkg := range packagesToUpdate {
		for _, varFile := range pkg.VarFiles {
			if err := updateVarFile(varFile); err != nil {
				return err
			}
		}
	}
	return nil
}

func updateVarFile(varFile string) error {
	firstLine, err := readFirstLine(varFile)
	if err != nil {
		return fmt.Errorf("reading first line from %s: %w", varFile, err)
	}

	// If the varFile has a schema line, get the template and version from it. It makes the update process more robust.
	// If not, attempt to update the varFile anyway.

	metadata := metadata.VarFileMetadata{
		HasVersion: false,
		Template:   "",
		Version:    "",
	}

	if strings.HasPrefix(firstLine, "# yaml-language-server: $schema=") {
		template, version, err := parseSchemaLine(firstLine, varFile)
		if err != nil {
			return fmt.Errorf("parsing schema line: %w", err)
		}

		metadata.HasVersion = true
		metadata.Template = template
		metadata.Version = version
	}

	err = update(varFile, metadata)
	if err != nil {
		return fmt.Errorf("updating varFile %s: %w", varFile, err)
	}
}

func readFirstLine(filename string) (string, error) {
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

func parseSchemaLine(line, varFile string) (string, string, error) {
	re := regexp.MustCompile(`\.schemas/(\w+)-(\S+)\.schema\.json`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 3 {
		return "", "", fmt.Errorf("invalid schema format in file %s", varFile)
	}

	template, version := matches[1], matches[2]
	return template, version, nil
}

func update(varFile string, metadata VarFileMetadata) error {
	err := add_apex_domain.AddApexDomainSupport(varFile, metadata)
	if err != nil {
		return err
	}

	return nil
}
