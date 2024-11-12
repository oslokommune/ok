package add_apex_domain

import (
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/update/migrate_config/metadata"
	"log/slog"
	"strconv"
	"strings"

	"github.com/magefile/mage/sh"
)

func AddApexDomainSupport(varFile string, metadata metadata.VarFileMetadata) error {
	slog.Debug("adding apex support", slog.String("varFile", varFile), slog.Any("metadata", metadata))

	if metadata.Template != "app" {
		slog.Debug("not updating, template is not app", slog.String("varFile", varFile))
		return nil
	}

	isTransformed, err := isTransformed(varFile)
	if err != nil {
		return fmt.Errorf("checking if YAML already is transformed: %w", err)
	}

	if isTransformed {
		slog.Debug("not updating, is already transformed")
		return nil
	}

	return update(varFile)
}

// isTransformed returns true if AlbHostRouting.Subdomain or AlbHostRouting.Apex is set. This most likely means
// that the YAML file already has support for Apex domain routing.
func isTransformed(varFile string) (bool, error) {
	args := []string{
		"eval",
		".AlbHostRouting.Subdomain != null or .AlbHostRouting.Apex != null",
		varFile,
	}

	output, err := sh.Output("yq", args...)
	if err != nil {
		return false, fmt.Errorf("error checking AlbHostRouting fields: %w", err)
	}

	isAlreadyTransformed, err := strconv.ParseBool(strings.TrimSpace(output))
	if err != nil {
		return false, fmt.Errorf("error parsing yq output: %w", err)
	}

	return isAlreadyTransformed, nil
}

func update(varFile string) error {
	fmt.Printf("Transforming configuration to add support for Apex domain routing. File: %s\n", varFile)

	// Proceed with the transformation
	args := []string{
		"-i",
		`
    .AlbHostRouting = {
        "Enable": (.AlbHostRouting.Enable // false),
        "Internal": (.AlbHostRouting.Internal // false),
        "Subdomain": {
            "Enable": (.AlbHostRouting.Enable // false),
            "TargetGroupTargetStickiness": (.AlbHostRouting.TargetGroupTargetStickiness // false)
        },
        "Apex": {
            "Enable": false,
            "TargetGroupTargetStickiness": false
        }
    } |
	del(.AlbHostRouting.TargetGroupTargetStickiness)
    `,
		varFile,
	}

	err := sh.RunV("yq", args...)
	if err != nil {
		return fmt.Errorf("error transforming YAML: %w", err)
	}

	fmt.Println("Transformation complete")

	return nil
}
