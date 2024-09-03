package boilerplate

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/oslokommune/ok/pkg/pkg/common"
)

type CommandlineBoilerplate struct {
	stdout io.WriteCloser
	stderr io.WriteCloser
}

var _ TemplateRenderer = (*CommandlineBoilerplate)(nil)

func NewCommandlineRenderer(stdout, stderr io.WriteCloser) *CommandlineBoilerplate {
	return &CommandlineBoilerplate{stdout: stdout, stderr: stderr}
}

func (c *CommandlineBoilerplate) Render(ctx context.Context, pkgManifestFilename string, outputFolders []string, baseUrlOrPath string) error {
	slog.Info("Installing packages...")

	// Install packages
	cmds, err := createBoilerplateCmdsFromManifest(ctx, pkgManifestFilename, outputFolders, baseUrlOrPath)
	if err != nil {
		return fmt.Errorf("creating boilerplate commands: %w", err)
	}

	for _, cmd := range cmds {
		cmd.Stdout = c.stdout
		cmd.Stderr = c.stderr
		printPrettyCmd(cmd)
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("running boilerplate command: %w", err)
		}
	}

	return nil
}

func createBoilerplateCmdsFromManifest(ctx context.Context, pkgManifestFilename string, outputFolders []string, baseUrlOrPath string) ([]*exec.Cmd, error) {
	manifest, err := common.LoadPackageManifest(pkgManifestFilename)
	if err != nil {
		return nil, fmt.Errorf("loading package manifest: %w", err)
	}

	// Filter packages based on output folders
	var packagesToInstall []common.Package
	if len(outputFolders) == 0 {
		packagesToInstall = manifest.Packages
	} else {
		packagesToInstall = filterPackages(manifest.Packages, outputFolders)
	}

	return createBoilerPlateExecCmds(ctx, packagesToInstall, manifest.DefaultPackagePathPrefix, baseUrlOrPath)
}

func createBoilerPlateExecCmds(ctx context.Context, packagesToInstall []common.Package, packagePathPrefix string, baseUrlOrPath string) ([]*exec.Cmd, error) {
	baseUrlOrPath = GetBaseUrlOrDefault(baseUrlOrPath)
	packagePathPrefix = GetPackagePathPrefixOrDefault(packagePathPrefix)

	var cmds []*exec.Cmd
	for _, pkg := range packagesToInstall {
		templateURL := getValidTemplareUrlOrPath(baseUrlOrPath, pkg.Template, packagePathPrefix, pkg.Ref)
		cmdArgs := []string{
			"--template-url", templateURL,
			"--output-folder", pkg.OutputFolder,
			"--non-interactive",
		}

		for _, varFile := range pkg.VarFiles {
			cmdArgs = append(cmdArgs, "--var-file", varFile)
		}

		cmd := exec.CommandContext(ctx, "boilerplate", cmdArgs...)
		cmds = append(cmds, cmd)
	}

	return cmds, nil
}

func isUrl(str string) bool {
	return strings.HasPrefix(str, "http://") ||
		strings.HasPrefix(str, "https://") ||
		strings.HasPrefix(str, "git@")
}

func filterPackages(packages []common.Package, outputFolders []string) []common.Package {
	outputFoldersLookupMap := makeLookupMap(outputFolders)
	result := make([]common.Package, 0, len(outputFolders))

	for _, pkg := range packages {
		if _, ok := outputFoldersLookupMap[pkg.OutputFolder]; ok {
			result = append(result, pkg)
		}
	}

	return result
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

func makeLookupMap[T comparable](slice []T) map[T]struct{} {
	m := make(map[T]struct{}, len(slice))
	for _, item := range slice {
		m[item] = struct{}{}
	}
	return m
}
