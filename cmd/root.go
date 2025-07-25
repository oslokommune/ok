package cmd

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/oslokommune/ok/cmd/aws"
	"github.com/oslokommune/ok/cmd/pk"
	"github.com/oslokommune/ok/cmd/pkg"
	"github.com/oslokommune/ok/pkg/error_user_msg"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
	"strings"
)

var (
	// rootCmd represents the base command when called without any subcommands.
	okString = "`ok`"
	rootCmd  = &cobra.Command{
		Use:   "ok",
		Short: "The `ok` infrastructure toolbox.",
		Long: fmt.Sprintf(`The %s tool is a comprehensive infrastructure management toolbox designed to streamline the setup and maintenance of Terraform environments. It provides a variety of commands to bootstrap infrastructure, manage environment configurations, handle AWS operations, and more.

Key functionalities include:

- Executing AWS-specific commands.
- Managing and updating Boilerplate templates.

Whether you're setting up a new environment or maintaining an existing one, %s simplifies and automates many of the repetitive tasks involved in infrastructure management.`, okString, okString),
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// cfgFile stores the path to the configuration file specified by the user.
	cfgFile string

	// defaultConfigPath stores the default configuration file path.
	defaultConfigPath string
)

// Execute is the main entry point for our program.
func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		redStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("1")) // Red text

		fmt.Println()
		fmt.Println(redStyle.Render("Error:"))
		prettyPrintError(err)

		var userError *error_user_msg.ErrorUserMessage
		if errors.As(err, &userError) {
			fmt.Println()
			fmt.Println(error_user_msg.StyleTitle.Render("Details:"))
			fmt.Println(userError.Details())
		}

		os.Exit(1)
	}
}

// init initializes the root command, setting up the configuration file path and flags.
func init() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	defaultConfigPath = path.Join(home, ".config", "ok", "config.yml")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is %s)", defaultConfigPath))

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
		fmt.Println("Using config file:", viper.ConfigFileUsed())
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
	out := strings.Repeat(" ", depth*2) + text
	fmt.Println(out)

}
