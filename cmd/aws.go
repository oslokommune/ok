package cmd

import (
	"github.com/spf13/cobra"
)

var awsCommand = &cobra.Command{
	Use:   "aws",
	Short: "Do AWS stuff",
}
