package add_apex_domain

import (
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/update/migrate_config"
	"github.com/rogpeppe/go-internal/semver"
	"strconv"
	"strings"

	"github.com/magefile/mage/sh"
)

func AddApexDomainSupport(varFile string, metadata migrate_config.VarFileMetadata) error {
	runUpdate :=
		!metadata.HasVersion ||
			(metadata.HasVersion &&
				(metadata.Template == "app" && semver.Compare(metadata.Version, "9.0.0") >= 0)) // metadata.Version >= "9.0.0"
	if !runUpdate {
		return nil
	}

	isAlreadyTransformed, err := isAlreadyTransformed(varFile)
	if err != nil {
		return fmt.Errorf("checking if YAML already is transformed: %w", err)
	}

	if isAlreadyTransformed {
		return nil
	}

	return update(varFile)
}

// isAlreadyTransformed returns true if AlbHostRouting.Subdomain or AlbHostRouting.Apex is set. This most likely means
// that the YAML file already has support for Apex domain routing.
func isAlreadyTransformed(varFile string) (bool, error) {
	checkCmd := []string{
		"eval",
		".AlbHostRouting.Subdomain != null or .AlbHostRouting.Apex != null",
		varFile,
	}

	output, err := sh.Output("yq", checkCmd...)
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
	fmt.Printf("Transforming %s to add support for Apex domain routing\n", varFile)

	// Proceed with the transformation
	args := []string{
		"-i",
		`
		.AlbHostRouting = {
			"Enable": .AlbHostRouting.Enable,
			"Internal": .AlbHostRouting.Internal,
			"Subdomain": {
				"Enable": .AlbHostRouting.Enable,
				"TargetGroupTargetStickiness": .AlbHostRouting.TargetGroupTargetStickiness
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
