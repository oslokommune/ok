package main

import (
	"fmt"
	"os"

	"github.com/oslokommune/ok/cmd"
)

func main() {
	rootCmd := cmd.NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
