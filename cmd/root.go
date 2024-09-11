package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/oslokommune/ok/cmd/aws"
	"github.com/oslokommune/ok/cmd/pkg"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// rootCmd represents the base command when called without any subcommands.
	rootCmd = &cobra.Command{
		Use:   "ok",
		Short: "The `ok` infrastructure toolbox.",
		Long: "The " + "`ok` " + `tool helps you to create a fresh Terraform environment (like prod or development) and configure it to use remote state storage.

Your environment is configured using a ` + "`packages.yml`" + ` file. This file is a package manifest listing the components from Golden Path that you wish to use. An example can be found in the ` + "`pirates-iac` " + `repository.`,
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

	rootCmd.AddCommand(pkgCommand)
	pkgCommand.AddCommand(pkg.AddCommand)
	pkgCommand.AddCommand(pkg.InstallCommand)
	pkgCommand.AddCommand(pkg.UpdateCommand)
	pkgCommand.AddCommand(pkg.FmtCommand)
	pkgCommand.AddCommand(pkg.SchemaCommand)

	rootCmd.AddCommand(awsCommand)
	awsCommand.AddCommand(aws.EcsExecCommand)

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
