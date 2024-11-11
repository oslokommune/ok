package migrate_config

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/magefile/mage/sh"
	"strings"
)

func addApexDomainSupport(varFile string) {
	// Update version using yq. Use shell exec yq with some parameters

	combinedArgs := []string{
		"yq",
		"something-command",
	}
	green := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	fmt.Println("------------------------------------------------------------------------------------------")
	fmt.Println("Running aws command:")
	fmt.Println(green.Render("aws " + strings.Join(combinedArgs, " ")))
	fmt.Println("------------------------------------------------------------------------------------------")

	_ = sh.RunV("aws", combinedArgs...)
}

TODO GPT with tests.
