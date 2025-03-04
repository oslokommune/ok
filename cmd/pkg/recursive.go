package pkg

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/spf13/cobra"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type RunRecursive func(manifestPath string, manifestDir string, style lipgloss.Style) error

func runRecursiveInSubdirs(f RunRecursive) error {
	var manifestPaths []string

	err := filepath.Walk(".", func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == "." {
			return nil
		}

		if !fileInfo.IsDir() {
			return nil
		}

		if strings.HasPrefix(fileInfo.Name(), ".") {
			return filepath.SkipDir
		}

		// Is there a package manifest in this directory?
		manifestPath := filepath.Join(path, common.PackagesManifestFilename)

		_, err = os.Stat(manifestPath)
		if err != nil {
			return nil
		}

		manifestPaths = append(manifestPaths, manifestPath)

		return nil
	})

	if err != nil {
		return fmt.Errorf("walking the path: %w", err)
	}

	for _, manifestPath := range manifestPaths {
		manifestDir := path.Dir(manifestPath)

		var style = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FF00")).
			Background(lipgloss.Color("22")).
			PaddingTop(2).
			PaddingBottom(2).
			PaddingLeft(4).
			PaddingRight(4).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFFFFF"))

		err = f(manifestPath, manifestDir, style)
		if err != nil {
			return err
		}

	}

	return nil
}

func addRecursiveFlagToCmd(cmd *cobra.Command, flag *bool, commandName string) {
	cmd.Flags().BoolVarP(flag,
		"recursive",
		"r",
		false,
		fmt.Sprintf("%s packages from manifests found in all subdirectories, but excluding the current directory.", commandName),
	)
}
