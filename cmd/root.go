package cmd

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/oslokommune/ok/cmd/aws"
	"github.com/oslokommune/ok/cmd/data"
	"github.com/oslokommune/ok/cmd/pk"
	"github.com/oslokommune/ok/cmd/pkg"
	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
	"github.com/oslokommune/ok/pkg/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
	"strconv"
	"strings"
)

var (
	// rootCmd represents the base command when called without any subcommands.
	rootCmd = &cobra.Command{
		Use:   "ok",
		Short: "Manage infrastructure and project templates.",
		Long: "\nThis toolbox helps you with:\n\n" +
			"- Infrastructure setup\n" +
			"- Terraform & CI/CD templates\n" +
			"- AWS tasks\n\n" +
			"[Get started](https://km.oslo.systems/setup/infrastructure/initialize-environment/)!",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			var nonInteractive *bool
			if cmd.Flags().Changed(FlagNonInteractiveName) {
				on, _ := cmd.Flags().GetBool(FlagNonInteractiveName)
				nonInteractive = &on
			}
			prompt.SetMode(promptMode(nonInteractive, os.Getenv, prompt.HasTerminal()))
		},
	}

	// cfgFile stores the path to the configuration file specified by the user.
	cfgFile string

	// defaultConfigPath stores the default configuration file path.
	defaultConfigPath string
)

const (
	FlagNonInteractiveName  = "non-interactive"
	FlagNonInteractiveUsage = "Never prompt; exit 2 if input is needed. On by default with OK_NON_INTERACTIVE=true, " +
		"CI=true, or without a terminal. --non-interactive=false turns prompts back on"

	exitCodeError = 1
)

// Execute is the main entry point for our program.
func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		redStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("1")) // Red text

		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, redStyle.Render("Error:"))
		prettyPrintError(err)

		var withDetails interface{ Details() string }
		if errors.As(err, &withDetails) && withDetails.Details() != "" {
			fmt.Fprintln(os.Stderr)
			fmt.Fprintln(os.Stderr, common.StyleTitle.Render("Details:"))
			fmt.Fprintln(os.Stderr, withDetails.Details())
		}

		os.Exit(exitCode(err))
	}
}

// exitCode is the exit code the error asks for, or exitCodeError.
func exitCode(err error) int {
	var withExitCode interface{ ExitCode() int }
	if errors.As(err, &withExitCode) {
		return withExitCode.ExitCode()
	}
	return exitCodeError
}

// usageTemplateWithExitCodes adds an Exit Codes section after the flags in a cobra usage template.
func usageTemplateWithExitCodes(template string) string {
	exitCodes := fmt.Sprintf("\n\nExit Codes:\n"+
		"  0  success\n"+
		"  %d  error\n"+
		"  %d  input needed, but ok may not prompt (see --%s)",
		exitCodeError, prompt.ExitCodeInputNeeded, FlagNonInteractiveName)

	const afterFlags = "{{if .HasHelpSubCommands}}"
	return strings.Replace(template, afterFlags, exitCodes+afterFlags, 1)
}

// promptMode decides whether ok may prompt. The first one set wins: --non-interactive (nil when
// not given), OK_NON_INTERACTIVE, CI=true. When none is set, ok prompts if it has a terminal.
func promptMode(nonInteractive *bool, getenv func(string) string, hasTerminal bool) prompt.Mode {
	if nonInteractive != nil {
		if *nonInteractive {
			return prompt.Mode{Reason: "--" + FlagNonInteractiveName + " is set"}
		}
		return prompt.Mode{Interactive: true}
	}
	if v := getenv("OK_NON_INTERACTIVE"); v != "" {
		if on, err := strconv.ParseBool(v); err == nil {
			if on {
				return prompt.Mode{Reason: "OK_NON_INTERACTIVE=" + v + " is set"}
			}
			return prompt.Mode{Interactive: true}
		}
	}
	if v := getenv("CI"); v != "" {
		if on, err := strconv.ParseBool(v); err == nil && on {
			return prompt.Mode{Reason: "CI=" + v + " is set"}
		}
	}
	if !hasTerminal {
		return prompt.Mode{Reason: "there is no terminal"}
	}
	return prompt.Mode{Interactive: true}
}

// init initializes the root command, setting up the configuration file path and flags.
func init() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	defaultConfigPath = path.Join(home, ".config", "ok", "config.yml")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is %s)", defaultConfigPath))
	rootCmd.PersistentFlags().Bool(FlagNonInteractiveName, false, FlagNonInteractiveUsage)
	rootCmd.SetUsageTemplate(usageTemplateWithExitCodes(rootCmd.UsageTemplate()))

	// Create dependencies
	ghReleases := githubreleases.NewGitHubReleases()
	addCommand := pkg.NewAddCommand(ghReleases)
	updateCommand := pkg.NewUpdateCommand(ghReleases)
	installCommand := pkg.NewInstallCommand()

	// Add commands
	rootCmd.AddCommand(pkgCommand)
	pkgCommand.AddCommand(installCommand)
	pkgCommand.AddCommand(addCommand)
	pkgCommand.AddCommand(updateCommand)
	pkgCommand.AddCommand(pkg.FmtCommand)

	rootCmd.AddCommand(awsCommand)
	awsCommand.AddCommand(aws.EcsExecCommand)
	awsCommand.AddCommand(aws.AdminSessionCommand)
	awsCommand.AddCommand(aws.ConfigGeneratorCommand)

	rootCmd.AddCommand(dataCommand)
	dataCommand.AddCommand(data.InitCommand)

	initializeConfiguration()

	if viper.GetBool("enable_experimental") {
		rootCmd.AddCommand(pkCommand)
		pkCommand.AddCommand(pk.NewInstallCommand())
	}
}

// initializeConfiguration is the function that initializes configuration using viper. It is called at the start of the application.
func initializeConfiguration() {
	setConfigFile()
	viper.SetDefault("enable_experimental", false)
	viper.SetEnvPrefix("ok")
	viper.AutomaticEnv()
	loadConfiguration()
}

// loadConfiguration attempts to read the configuration file using viper and prints the path of the used configuration file if successful.
func loadConfiguration() {
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// setConfigFile sets the configuration file for viper. If cfgFile is specified, it's used; otherwise, the default path is set.
func setConfigFile() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		return
	}

	setDefaultConfigPath()
}

// setDefaultConfigPath configures viper with the default configuration path.
func setDefaultConfigPath() {
	viper.AddConfigPath(path.Dir(defaultConfigPath))
	viper.SetConfigType("yaml")
	viper.SetConfigName(path.Base(defaultConfigPath))
}

func prettyPrintError(err error) {
	i := 0

	for {
		errStr := err.Error()

		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			printWithSpaces(errStr, i)
			break
		}

		unwrappedStr := unwrapped.Error()

		// As an example, errStr is:
		// middle b that wraps a: deepest error a
		//
		// But we want to print only one error at a time, like this:
		// middle b that wraps a:
		//
		// So we search for "deepest error a" (the unwrapped error) from the complete error string, and remove it.
		text := strings.Replace(errStr, unwrappedStr, "", 1)
		printWithSpaces(text, i)

		err = unwrapped
		i++

		if i > 100 {
			printWithSpaces("(Too many errors to unwrap, stopping here.)", i)
			break
		}
	}
}

func printWithSpaces(text string, depth int) {
	indent := strings.Repeat(" ", depth*2)
	out := indent + strings.ReplaceAll(text, "\n", "\n"+indent)
	fmt.Fprintln(os.Stderr, out)

}
