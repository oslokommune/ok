package install

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/charmbracelet/lipgloss"
)

var warnStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("3")) // Yellow

// warnIfInvalidBoilerplateVersion prints a warning if boilerplate has an invalid semver version.
// Some self-built versions report "latest" which causes failures when a Boilerplate template uses `required_version`
func warnIfInvalidBoilerplateVersion() {
	version, err := getBoilerplateVersion()
	if err != nil {
		// Don't block install if we can't determine the version
		return
	}

	_, err = semver.NewVersion(version)
	if err != nil {
		fmt.Println(warnStyle.Render(fmt.Sprintf(
			"️⚠️ Warning: boilerplate version %q is not a valid semantic version.", version)))
		fmt.Println("Some Golden Path templates may fail because they require a valid version (e.g. v0.12.1).")
		fmt.Println("This can happen if boilerplate was built from source without a proper version tag.")
		fmt.Println()
	}
}

// getBoilerplateVersion runs "boilerplate --version" and extracts the version string.
func getBoilerplateVersion() (string, error) {
	out, err := exec.Command("boilerplate", "--version").Output()
	if err != nil {
		return "", fmt.Errorf("running 'boilerplate --version': %w", err)
	}

	return parseBoilerplateVersion(string(out)), nil
}

// parseBoilerplateVersion extracts the version from boilerplate --version output.
// Expected format: "boilerplate version v0.12.1" or "boilerplate version latest".
func parseBoilerplateVersion(output string) string {
	output = strings.TrimSpace(output)
	parts := strings.Fields(output)

	if len(parts) >= 3 {
		return parts[len(parts)-1]
	}

	return output
}
