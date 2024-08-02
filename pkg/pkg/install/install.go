package install

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/oslokommune/ok/pkg/pkg/common"
	"os"
	"os/exec"
	"strings"
)

const DefaultBaseUrl = "git@github.com:oslokommune/golden-path-boilerplate.git//boilerplate/terraform"

func Run(pkgManifestFilename string, outputFolders []string) error {
	cmds, err := CreateBoilerplateCommands(pkgManifestFilename, outputFolders)
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

func CreateBoilerplateCommands(pkgManifestFilename string, outputFolders []string) ([]*exec.Cmd, error) {
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
	cmds := createBoilerPlateCommands(packagesToInstall)

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

func createBoilerPlateCommands(packagesToInstall []common.Package) []*exec.Cmd {
	var cmds []*exec.Cmd
	for _, pkg := range packagesToInstall {

		var EnvBaseUrl = os.Getenv("BASE_URL")
		var templateURL string
		if EnvBaseUrl == "" {
			templateURL = fmt.Sprintf("%s/%s?ref=%s", DefaultBaseUrl, pkg.Template, pkg.Ref)
		} else {
			templateURL = fmt.Sprintf("%s/%s", EnvBaseUrl, pkg.Template)
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
	return cmds
}
