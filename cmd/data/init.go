package data

import (
	"github.com/oslokommune/ok/pkg/data"
	"github.com/spf13/cobra"
)

var (
	flagBranch     string
	flagConfigFile string
	flagOutputDir  string
	flagTag        string
	flagTemplateDir string
)

// InitCommand initializes a new Databricks bundle using Oslo Kommune template
var InitCommand = &cobra.Command{
	Use:   "init [TEMPLATE_PATH]",
	Short: "Initialize a new Databricks bundle using Oslo Kommune template",
	Long: `Initialize a new Databricks bundle with the Oslo Kommune template as default.

The TEMPLATE_PATH is optional. If not provided, uses the Oslo Kommune custom template.
You can override the default template URL with the DATA_TEMPLATE_URL environment variable.

This command is a wrapper around 'databricks bundle init' that uses Oslo Kommune's
standard Databricks project template by default, while still allowing you to use
any other template (built-in or custom) by specifying it explicitly.`,
	Example: `  # Use Oslo Kommune template (default)
  ok data init

  # Use built-in Python template
  ok data init default-python

  # Use custom Git template
  ok data init github.com/my/template

  # Specify output directory
  ok data init --output-dir ./my-project

  # Override default template via environment variable
  DATA_TEMPLATE_URL=custom/url ok data init

  # Use specific branch of template
  ok data init --branch main

  # Use specific tag of template
  ok data init --tag v1.0.0`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		templatePath := ""
		if len(args) > 0 {
			templatePath = args[0]
		}

		return data.RunInit(data.InitOptions{
			TemplatePath: templatePath,
			Branch:       flagBranch,
			ConfigFile:   flagConfigFile,
			OutputDir:    flagOutputDir,
			Tag:          flagTag,
			TemplateDir:  flagTemplateDir,
		})
	},
}

func init() {
	InitCommand.Flags().StringVar(&flagBranch, "branch", "", "Git branch to use for template initialization")
	InitCommand.Flags().StringVar(&flagConfigFile, "config-file", "", "JSON file containing key value pairs of input parameters")
	InitCommand.Flags().StringVar(&flagOutputDir, "output-dir", "", "Directory to write the initialized template to")
	InitCommand.Flags().StringVar(&flagTag, "tag", "", "Git tag to use for template initialization")
	InitCommand.Flags().StringVar(&flagTemplateDir, "template-dir", "", "Directory path within a Git repository containing the template")
}
