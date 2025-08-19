package install

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/oslokommune/ok/pkg/pkg/common"
	"os"
	"os/exec"
	"path"
	"strings"
)

// Run runs Boilerplate for the specified packages.
func Run(packagesToInstall []common.Package, manifest common.PackageManifest, workingDirectory string) error {
	cmds, err := CreateBoilerplateCommands(packagesToInstall, CreateBoilerPlateCommandsOpts{
		PackagePathPrefix: manifest.PackagePrefix(),
		BaseUrlOrPath:     os.Getenv(common.BaseUrlEnvName),
		WorkingDirectory:  workingDirectory,
	})

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

	common.PrintProcessedPackages(packagesToInstall, "installed")

	return nil
}

func printPrettyCmd(cmd *exec.Cmd) {
	cmdString := createPrettyCmdString(cmd)
	green := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))

	fmt.Println("------------------------------------------------------------------------------------------")
	fmt.Printf("Running boilerplate command in dir '%s':\n", cmd.Dir)
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

// CreateBoilerplateCommands create Boilerplate commands for the specified packages.
func CreateBoilerplateCommands(packages []common.Package, opts CreateBoilerPlateCommandsOpts) ([]*exec.Cmd, error) {
	var cmds []*exec.Cmd

	for _, pkg := range packages {
		if opts.BaseUrlOrPath == "" {
			opts.BaseUrlOrPath = common.DefaultBaseUrl
		}

		var templateURL string
		if isUrl(opts.BaseUrlOrPath) {
			pathz := strings.Join(
				[]string{opts.PackagePathPrefix, pkg.Template}, "/")

			templateURL = fmt.Sprintf("%s%s?ref=%s", opts.BaseUrlOrPath, pathz, pkg.Ref)
		} else {
			templateURL = path.Join(opts.BaseUrlOrPath, opts.PackagePathPrefix, pkg.Template)
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
		cmd.Dir = opts.WorkingDirectory
		cmds = append(cmds, cmd)
	}

	return cmds, nil
}

func isUrl(str string) bool {
	return strings.HasPrefix(str, "http://") ||
		strings.HasPrefix(str, "https://") ||
		strings.HasPrefix(str, "git@")
}

type CreateBoilerPlateCommandsOpts struct {
	PackagePathPrefix string
	BaseUrlOrPath     string
	WorkingDirectory  string
}
