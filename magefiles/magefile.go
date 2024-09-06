//go:build mage

package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/magefile/mage/sh"
	"github.com/staticaland/brandish/pipelines/markdownlint"
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

// Lint
func Lint(ctx context.Context) error {

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	fmt.Println(markdownlint.Lint(client))

	return nil
}

// Lint and fix
func LintFix(ctx context.Context) error {

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	fmt.Println(
		markdownlint.Lint(
			client,
			markdownlint.WithGlobs("README.md"),
			markdownlint.WithFix(),
		),
	)

	return nil
}
