package add_path_based_routing

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/oslokommune/ok/pkg/pkg/update/migrate_config/metadata"

	"github.com/magefile/mage/sh"
)

const RequiredVersion = "11.0.0"

func AddPathBasedRouting(varFile string, varFileJsonSchema metadata.JsonSchema) error {
	slog.Debug("adding support for path based routing",
		slog.String("varFile", varFile),
		slog.Any("varFileJsonSchema", varFileJsonSchema),
	)

	/*
	 * Check that var file supports the new feature.
	 * https://github.com/oslokommune/golden-path-boilerplate/releases/tag/app-v11.0.0
	 */

	requiredVersion, err := semver.NewVersion(RequiredVersion)
	if err != nil {
		return fmt.Errorf("creating semver: %w", err)
	}

	if varFileJsonSchema.Version.LessThan(requiredVersion) {
		slog.Debug("not updating, package less than required version",
			slog.String("requiredVersion", requiredVersion.String()),
			slog.String("varFile", varFile))
		return nil
	}

	if varFileJsonSchema.Template != "app" {
		slog.Debug("not updating, var file template is not 'app'")
		return nil
	}

	// Check if already transformed
	migrated, err := isMigrated(varFile)
	if err != nil {
		return fmt.Errorf("checking if YAML already is transformed: %w", err)
	}

	if migrated {
		slog.Debug("not updating, is already transformed")
		return nil
	}

	return migrate(varFile)
}

// isMigrated returns true if var file already has been transformed.
func isMigrated(varFile string) (bool, error) {
	// To test this logic in your terminal, run:
	// yq ".ApplicationLoadBalancer != null" app-too-tikki.yml

	args := []string{
		"eval",
		".ApplicationLoadBalancer != null",
		varFile,
	}

	output, err := sh.Output("yq", args...)
	if err != nil {
		return false, fmt.Errorf("error running yq: %w", err)
	}

	isAlreadyTransformed, err := strconv.ParseBool(strings.TrimSpace(output))
	if err != nil {
		return false, fmt.Errorf("error parsing yq output: %w", err)
	}

	return isAlreadyTransformed, nil
}

func migrate(varFile string) error {
	fmt.Printf("Changing var file from using 'AlbHostRouting' to 'ApplicationLoadBalancer'. File: %s\n", varFile)

	hasMetadata := false
	firstLine, err := metadata.ReadFirstLine(varFile)
	if err != nil {
		return fmt.Errorf("reading file %s: %w", varFile, err)
	}

	_, err = metadata.ParseMetadataLine(firstLine)
	if err == nil {
		hasMetadata = true
	}

	args := getArgs(varFile)

	err = sh.RunV("yq", args...)
	if err != nil {
		return fmt.Errorf("error transforming YAML: %w", err)
	}

	if hasMetadata {
		// The yq command will remove the json schema declaration, so let's write it back.
		content, err := os.ReadFile(varFile)
		if err != nil {
			return fmt.Errorf("reading file to restore metadata: %w", err)
		}

		newContent := firstLine + "\n" + string(content)
		err = os.WriteFile(varFile, []byte(newContent), 0644)
		if err != nil {
			return fmt.Errorf("writing file with restored metadata: %w", err)
		}
	}

	return nil
}

func getArgs(varFile string) []string {
	/*
		    # yq command to test the transformation logic:
			yq '
			# Convert to entries to maintain order
			to_entries |
			# Process each entry
			map(
			  # For AlbHostRouting entries
			  select(.key == "AlbHostRouting") |= (
				# Store original value first
				.orig = .value |
				# Change key to ApplicationLoadBalancer
				.key = "ApplicationLoadBalancer" |
				# Create new structure with basic fields
				.value = {
				  "Enable": (.orig.Enable // false),
				  "Internal": (.orig.Internal // false),
				  "HostRouting": {}
				} |
				# Add Subdomain if it exists
				select(.orig.Subdomain != null) |= (
				  .value.HostRouting.Subdomain = {
					"Enable": (.orig.Subdomain.Enable // false),
					"TargetGroupTargetStickiness": (.orig.Subdomain.TargetGroupTargetStickiness // false)
				  }
				) |
				# Add ApexDomain if it exists
				select(.orig.ApexDomain != null) |= (
				  .value.HostRouting.ApexDomain = {
					"Enable": (.orig.ApexDomain.Enable // false),
					"TargetGroupTargetStickiness": (.orig.ApexDomain.TargetGroupTargetStickiness // false)
				  }
				) |
				# Remove temporary field
				del(.orig)
			  )
			) |
			# Convert back to object
			from_entries
			' app-too-tikki.yml
	*/

	// Proceed with the transformation
	args := []string{
		"-i",
		`
		   # Convert to entries to maintain order
			to_entries |
			# Process each entry
			map(
			  # For AlbHostRouting entries
			  select(.key == "AlbHostRouting") |= (
				# Store original value first
				.orig = .value |
				# Change key to ApplicationLoadBalancer
				.key = "ApplicationLoadBalancer" |
				# Create new structure with basic fields
				.value = {
				  "Enable": (.orig.Enable // false),
				  "Internal": (.orig.Internal // false),
				  "HostRouting": {}
				} |
				# Add Subdomain if it exists
				select(.orig.Subdomain != null) |= (
				  .value.HostRouting.Subdomain = {
					"Enable": (.orig.Subdomain.Enable // false),
					"TargetGroupTargetStickiness": (.orig.Subdomain.TargetGroupTargetStickiness // false)
				  }
				) |
				# Add ApexDomain if it exists
				select(.orig.ApexDomain != null) |= (
				  .value.HostRouting.ApexDomain = {
					"Enable": (.orig.ApexDomain.Enable // false),
					"TargetGroupTargetStickiness": (.orig.ApexDomain.TargetGroupTargetStickiness // false)
				  }
				) |
				# Remove temporary field
				del(.orig)
			  )
			) |
			# Convert back to object
			from_entries
    `,
		varFile,
	}
	return args
}
