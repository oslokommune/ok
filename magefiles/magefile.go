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
