package pkg

import (
	"fmt"
	"strings"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/update"
	"github.com/spf13/cobra"
)

var flagUpdateCommandUpdateSchema bool

var UpdateCommand = &cobra.Command{
	Use:   "update [package-name]",
	Short: "Update Boilerplate package manifest",
	Long: `Update Boilerplate package manifest.
If a package name is provided, only that package will be updated.
If no package name is provided, all packages will be updated.`,
	Example: `  ok pkg update
  ok pkg update my-package`,
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: updateTabCompletion,
	RunE: func(cmd *cobra.Command, args []string) error {
		var packageName string
		if len(args) == 1 {
			packageName = args[0]
		}

		err := update.Run(PackagesManifestFilename, packageName, flagUpdateCommandUpdateSchema)
		if err != nil {
			return err
		}

		if packageName != "" {
			fmt.Printf("Updated package: %s\n", packageName)
		} else {
			fmt.Println("Updated all packages")
		}

		return nil
	},
}

func init() {
	UpdateCommand.Flags().BoolVar(&flagUpdateCommandUpdateSchema, "update-schema", true, "Update the JSON schema for affected packages")
}

func updateTabCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// No completion if there are more than one argument
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	manifest, err := common.LoadPackageManifest(PackagesManifestFilename)
	if err != nil {
		cmd.PrintErrf("failed to load package manifest: %s\n", err)
		return nil, cobra.ShellCompDirectiveError
	}

	var completions []string
	for _, p := range manifest.Packages {
		if strings.HasPrefix(p.OutputFolder, toComplete) {
			completions = append(completions, p.OutputFolder)
		}
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}
