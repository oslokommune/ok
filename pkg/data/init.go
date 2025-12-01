package data

import (
	"os"
	"os/exec"

	"github.com/oslokommune/ok/pkg/error_user_msg"
)

// InitOptions contains all the options for initializing a databricks bundle
type InitOptions struct {
	TemplatePath string // User-provided template (optional)
	Branch       string // databricks --branch flag
	ConfigFile   string // databricks --config-file flag
	OutputDir    string // databricks --output-dir flag
	Tag          string // databricks --tag flag
	TemplateDir  string // databricks --template-dir flag
}

// RunInit executes the databricks bundle init command with the configured template
func RunInit(opts InitOptions) error {
	// Check if databricks CLI is installed
	if err := checkDatabricksInstalled(); err != nil {
		return err
	}

	// Determine which template to use
	template := determineTemplate(opts.TemplatePath)

	// Update options with determined template
	opts.TemplatePath = template

	// Build and execute the databricks command
	cmd := buildDatabricksCommand(opts)

	if err := cmd.Run(); err != nil {
		errWithMsg := error_user_msg.NewError(
			"failed to initialize databricks bundle",
			"Check that the template URL is correct and accessible.\nRun 'databricks bundle init --help' for more information.",
			err,
		)
		return &errWithMsg
	}

	return nil
}

// checkDatabricksInstalled verifies that the databricks CLI is available
func checkDatabricksInstalled() error {
	_, err := exec.LookPath("databricks")
	if err != nil {
		errWithMsg := error_user_msg.NewError(
			"databricks CLI not found",
			"Please install the databricks CLI:\nhttps://docs.databricks.com/en/dev-tools/cli/install.html",
			err,
		)
		return &errWithMsg
	}
	return nil
}

// determineTemplate selects which template to use based on priority:
// 1. User-provided template (highest priority)
// 2. Environment variable override
// 3. Default Oslo kommune template
func determineTemplate(userTemplate string) string {
	// 1. User explicitly provided template (highest priority)
	if userTemplate != "" {
		return userTemplate
	}

	// 2. Environment variable override
	if envTemplate := os.Getenv(TemplateURLEnvName); envTemplate != "" {
		return envTemplate
	}

	// 3. Default Oslo kommune template
	return DefaultTemplateURL
}

// buildDatabricksCommand constructs the exec.Cmd for databricks bundle init
func buildDatabricksCommand(opts InitOptions) *exec.Cmd {
	args := []string{"bundle", "init"}

	// Template is positional argument (before flags)
	if opts.TemplatePath != "" {
		args = append(args, opts.TemplatePath)
	}

	// Add all optional flags
	if opts.Branch != "" {
		args = append(args, "--branch", opts.Branch)
	}
	if opts.ConfigFile != "" {
		args = append(args, "--config-file", opts.ConfigFile)
	}
	if opts.OutputDir != "" {
		args = append(args, "--output-dir", opts.OutputDir)
	}
	if opts.Tag != "" {
		args = append(args, "--tag", opts.Tag)
	}
	if opts.TemplateDir != "" {
		args = append(args, "--template-dir", opts.TemplateDir)
	}

	cmd := exec.Command("databricks", args...)
	// Pass through all I/O streams for interactive experience
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}
