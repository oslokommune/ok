package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var getCommand = &cobra.Command{
	Use:   "get",
	Short: "Get a template.",
	Run: func(cmd *cobra.Command, args []string) {
		baseURL := "git@github.com:oslokommune/golden-path-boilerplate.git//boilerplate"
		ref := "main"
		if version != "" {
			ref = version
		}
		if outputFolder == "" {
			outputFolder = templateName
		}

		templateURL := fmt.Sprintf("%s/%s?ref=%s", baseURL, templateName, ref)
		varFile := fmt.Sprintf("vars-%s.yml", outputFolder)
		commonVarFile := "vars-common.yml"

		boilerplateCmd := exec.Command("boilerplate",
			"--template-url", templateURL,
			"--var-file", varFile,
			"--var-file", commonVarFile,
			"--output-folder", outputFolder,
			"--non-interactive",
		)

		boilerplateCmd.Stdout = os.Stdout
		boilerplateCmd.Stderr = os.Stderr

		err := boilerplateCmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to execute command: %v\n", err)
			os.Exit(1)
		}
	},
}

var templateName string
var outputFolder string
var version string

func init() {
	getCommand.Flags().StringVarP(&templateName, "template", "t", "", "Template name (required)")
	getCommand.MarkFlagRequired("template")
	getCommand.Flags().StringVarP(&outputFolder, "output-folder", "o", "", "Output folder (optional)")
	getCommand.Flags().StringVarP(&version, "version", "v", "", "Version (optional)")
	rootCmd.AddCommand(getCommand)
}
