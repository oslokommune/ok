package install

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/oslokommune/ok/pkg/pkg/common"
	"os"
	"os/exec"
	"strings"
)

const DefaultBaseUrl = "git@github.com:oslokommune/golden-path-boilerplate.git//"
const DefaultPackagePathPrefix = "boilerplate/terraform"

func Run(pkgManifestFilename string, outputFolders []string) error {
	baseUrlOrPath := os.Getenv("BASE_URL")

	cmds, err := CreateBoilerplateCommands(pkgManifestFilename, outputFolders, baseUrlOrPath)
	if err != nil {
		return fmt.Errorf("creating boilerplate command: %w", err)
	}

	for _, cmd := range cmds {
		printPrettyCmd(cmd)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("running boilerplate command: %w", err)
		}
	}

	return nil
}

func printPrettyCmd(cmd *exec.Cmd) {
	cmdString := createPrettyCmdString(cmd)
	green := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))

	fmt.Println("------------------------------------------------------------------------------------------")
	fmt.Println("Running boilerplate command:")
	fmt.Println(green.Render(cmdString))
	fmt.Println("------------------------------------------------------------------------------------------")
}

func createPrettyCmdString(cmd *exec.Cmd) string {
	var argsStr string

	for _, arg := range cmd.Args[1:] {
		if strings.HasPrefix(arg, "--") {
			argsStr += "\n  " + arg
		} else {
			argsStr += " " + arg
		}
	}

	cmdString := fmt.Sprintf("%s%s", cmd.Path, argsStr)

	return cmdString
}

func CreateBoilerplateCommands(pkgManifestFilename string, outputFolders []string, baseUrlOrPath string) ([]*exec.Cmd, error) {
	fmt.Println("Installing packages...")

	manifest, err := common.LoadPackageManifest(pkgManifestFilename)
	if err != nil {
		return nil, fmt.Errorf("loading package manifest: %w", err)
	}

	// Filter packages based on output folders
	packagesToInstall := make([]common.Package, 0)
	if len(outputFolders) == 0 {
		packagesToInstall = manifest.Packages
	} else {
		packagesToInstall = filterPackages(manifest.Packages, outputFolders)
	}

	// Install packages
	cmds, err := createBoilerPlateCommands(packagesToInstall, manifest.DefaultPackagePathPrefix, baseUrlOrPath)
	if err != nil {
		return nil, fmt.Errorf("creating boilerplate commands: %w", err)
	}

	return cmds, nil
}

func filterPackages(packages []common.Package, outputFolders []string) []common.Package {
	result := make([]common.Package, 0)

	for _, pkg := range packages {
		for _, outputFolder := range outputFolders {

			if pkg.OutputFolder == outputFolder {
				result = append(result, pkg)
				break
			}

		}
	}

	return result
}

func createBoilerPlateCommands(packagesToInstall []common.Package, packagePathPrefix string, baseUrlOrPath string) ([]*exec.Cmd, error) {
	var cmds []*exec.Cmd
	for _, pkg := range packagesToInstall {
		if baseUrlOrPath == "" {
			baseUrlOrPath = DefaultBaseUrl
		}

		if packagePathPrefix == "" {
			packagePathPrefix = DefaultPackagePathPrefix
		}

		path := strings.Join(
			[]string{packagePathPrefix, pkg.Template},
			"/")

		// envBaseUrl can be a URL or a path
		var templateURL string
		if isUrl(baseUrlOrPath) {
			templateURL = fmt.Sprintf("%s%s?ref=%s", baseUrlOrPath, path, pkg.Ref)
		} else {
			templateURL = fmt.Sprintf("%s%s", baseUrlOrPath, path)

		}

		cmdArgs := []string{
			"--template-url", templateURL,
			"--output-folder", pkg.OutputFolder,
			"--non-interactive",
		}

		for _, varFile := range pkg.VarFiles {
			cmdArgs = append(cmdArgs, "--var-file", varFile)
		}

		cmd := exec.Command("boilerplate", cmdArgs...)
		cmds = append(cmds, cmd)
	}

	return cmds, nil
}

func isUrl(str string) bool {
	return strings.HasPrefix(str, "http://") ||
		strings.HasPrefix(str, "https://") ||
		strings.HasPrefix(str, "git@")
}
