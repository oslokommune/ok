package add_apex_domain

import (
	"fmt"
	"strings"

	"github.com/magefile/mage/sh"
)

func AddApexDomainSupport(varFile string) error {
	// Check if AlbHostRouting.Subdomain and AlbHostRouting.Apex exist
	checkCmd := []string{
		"eval",
		".AlbHostRouting.Subdomain, .AlbHostRouting.Apex",
		varFile,
	}

	output, err := sh.Output("yq", checkCmd...)
	if err != nil {
		return fmt.Errorf("error checking AlbHostRouting fields: %w", err)
	}

	// TODO this should be true/false..                                                      <------------------------
	// If both fields are null (don't exist), proceed with the transformation
	if strings.TrimSpace(output) != "null\nnull" {
		// AlbHostRouting.Subdomain or AlbHostRouting.Apex already exist. Skipping transformation.
		return nil
	}

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

	err = sh.RunV("yq", args...)
	if err != nil {
		return fmt.Errorf("error transforming YAML: %w", err)
	}

	fmt.Println("Transformation complete")

	return nil
}

/*
green := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	fmt.Println("------------------------------------------------------------------------------------------")
	fmt.Println("Running aws command:")
	fmt.Println(green.Render("aws " + strings.Join(combinedArgs, " ")))
	fmt.Println("------------------------------------------------------------------------------------------")

*/
