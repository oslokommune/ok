package install

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"strings"
)

const DefaultBaseUrl = "git@github.com:oslokommune/golden-path-boilerplate.git//boilerplate/terraform"

func Run(pkgManifestFilename string, stacks []string) error {
	cmds, err := CreateBoilerplateCommands(pkgManifestFilename, stacks)
	if err != nil {
		return fmt.Errorf("creating boilerplate command: %w", err)
	}

	curDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current directory: %w", err)
	}

	log.Debug().Msgf("Current working directory: %s", curDir)

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

func CreateBoilerplateCommands(filePath string, stacks []string) ([]*exec.Cmd, error) {
	var EnvBaseUrl = os.Getenv("BASE_URL")

	fmt.Println("Installing packages...")

	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	var manifest common.PackageManifest

	err = yaml.Unmarshal(fileContents, &manifest)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling YAML: %w", err)
	}

	var cmds []*exec.Cmd
	for _, pkg := range manifest.Packages {
		if len(stacks) > 0 {
			var found bool
			for _, stack := range stacks {
				if stack == pkg.OutputFolder {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

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

	return cmds, nil
}
