package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/oslokommune/ok/cmd/aws"
	"github.com/oslokommune/ok/cmd/pkg"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		fmt.Println(err)
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
	addCommand := pkg.NewAddCommand()
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
}

// initializeConfiguration is the function that initializes configuration using viper. It is called at the start of the application.
func initializeConfiguration() {
	setConfigFile()
	viper.SetDefault("enable_experimental", true)
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
