//go:build mage

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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

	// Step 4: Format the docs to make them more readable
	err := filepath.Walk("docs", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			text := string(content)
			text = regexp.MustCompile(`SEE ALSO`).ReplaceAllString(text, "See also")
			text = regexp.MustCompile(`(?m)^## `).ReplaceAllString(text, "# ")
			text = regexp.MustCompile(`(?m)^### `).ReplaceAllString(text, "## ")
			text = regexp.MustCompile(`(?m)^#### `).ReplaceAllString(text, "### ")
			text = regexp.MustCompile(`(?m)^##### `).ReplaceAllString(text, "#### ")

			if err := ioutil.WriteFile(path, []byte(text), info.Mode()); err != nil {
				return err
			}
		}
		return nil
	})

	return err
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
