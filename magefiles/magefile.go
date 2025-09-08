//go:build mage

package main

import (
	"os"

	"github.com/magefile/mage/sh"
)

// Build the project.
func Build() error {
	return sh.RunV("go", "build", "-o", "main", "main.go")
}

func Docs() error {
	// Step 1: Remove the docs directory
	if err := os.RemoveAll("docs"); err != nil {
		return err
	}

	// Step 2: Create the docs directory
	if err := os.Mkdir("docs", os.ModePerm); err != nil {
		return err
	}

	// Step 3: Generate the docs again
	if err := sh.RunV("go", "run", "main.go", "docs"); err != nil {
		return err
	}

	// Step 4: Optimize the docs
	if err := sh.RunV("node", "docs-optimizer/index.js"); err != nil {
		return err
	}

	return nil
}

// Run the project.
func Run() error {
	return sh.RunV("go", "run", "main.go")
}

// Test runs the unit tests.
func Test() error {
	return sh.RunV("go", "test", "./...")
}

