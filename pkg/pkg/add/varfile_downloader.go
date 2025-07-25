package add

import (
	"context"
	"errors"
	"fmt"
	"github.com/oslokommune/ok/pkg/error_user_msg"
	"github.com/oslokommune/ok/pkg/pkg/common"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func (a Adder) downloadVarFile(newPackage common.Package, varFile string, varFilePath string, outputFolder string) error {
	if _, err := os.Stat(varFilePath); err == nil {
		return fmt.Errorf("file already exists: %s", varFilePath)
	}

	// https://github.com/oslokommune/golden-path-boilerplate/tree/databases-v4.0.3/boilerplate/terraform/databases/package-config-default.yml
	varFileDownloadFilename := fmt.Sprintf("package-config-%s.yml", varFile)

	// boilerplate/terraform/{pkg.Template}/package-config-default.yml
	path := filepath.Join(common.BoilerplatePackageTerraformPath, newPackage.Template, varFileDownloadFilename)

	// Show URL to the user
	varFileUrl := fmt.Sprintf("https://github.com/%s/%s/tree/%s/%s/%s/%s",
		common.BoilerplateRepoOwner,
		common.BoilerplateRepoName,
		newPackage.Ref,
		common.BoilerplatePackageTerraformPath,
		newPackage.Template,
		varFileDownloadFilename,
	)
	fmt.Printf("Creating var file %s from %s\n", varFilePath, varFileUrl)

	fileBytes, err := a.ghReleases.DownloadGithubFile(
		context.Background(),
		common.BoilerplateRepoOwner,
		common.BoilerplateRepoName,
		path,
		newPackage.Ref,
	)
	if err != nil && strings.Contains(err.Error(), "no file named") {
		return createErrorDetails(err, newPackage)
	} else if err != nil {
		return fmt.Errorf("downloading file from GitHub: %w", err)
	}

	fileString := string(fileBytes)
	fileString = updateStackName(fileString, outputFolder)
	fileBytes = []byte(fileString)

	err = os.WriteFile(varFilePath, fileBytes, 0644)
	if err != nil {
		return fmt.Errorf("writing to file: %w", err)
	}

	return nil
}

var regexStackName = regexp.MustCompile(`StackName:\s*"[\w\-]*"`)
var regexStackNameForData = regexp.MustCompile(`([\w\-]+)\.StackName:\s*"[\w\-]*"`)

// updateStackName sets StackName in a var file.
func updateStackName(fileString, outputFolder string) string {
	/*
		Here are some examples of which cases to support.

		Example 1: Single stack, no data-stack.
		StackName: "backup"

		Example 2: Main stack and data stack.
		StackName: "app-km"
		app-data.StackName: "app-km-data"
	*/

	var replacement string

	// Replace StackName on line 1
	replacement = fmt.Sprintf(`StackName: "%s"`, outputFolder)
	fileString = regexStackName.ReplaceAllString(fileString, replacement)

	// Replace StackName on line 2, the one with "-data" in it
	fileString = regexStackNameForData.ReplaceAllStringFunc(fileString, func(m string) string {
		submatches := regexStackNameForData.FindStringSubmatch(m)

		if len(submatches) >= 2 {
			templateNameWithData := submatches[1]
			return fmt.Sprintf(`%s.StackName: "%s"`, templateNameWithData, outputFolder+"-data")
		}

		return m
	})

	return fileString
}

func createErrorDetails(sourceError error, newPackage common.Package) error {

	var errorDetails string
	errorDetails += fmt.Sprintf(
		"Template %s is missing var file %s.\n",
		error_user_msg.StyleHighlight.Render(newPackage.Template),
		error_user_msg.StyleHighlight.Render("package-config-default.yml"),
	)

	errorDetails += fmt.Sprintln()
	errorDetails += fmt.Sprintln(error_user_msg.StyleTitle.Render("Possible solutions:"))
	errorDetails += fmt.Sprintf(
		"- Use flag %s to remove this error.\n",
		error_user_msg.StyleHighlight.Render("--"+FlagNoVar),
	)
	errorDetails += fmt.Sprintf("- Ask maintainers to fix this error.\n")

	// Replace error with a new error that has the same error message and sub error, but
	// some nice error details alongside.
	errWithMsg := error_user_msg.NewError(sourceError.Error(), errorDetails, errors.Unwrap(sourceError))
	return &errWithMsg
}
