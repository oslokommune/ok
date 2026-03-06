package workflow

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/oslokommune/ok/pkg/error_user_msg"
)

// AppType represents the repository type variant for app workflows.
type AppType string

const (
	AppTypeDefault    AppType = ""
	AppTypeAppWithIac AppType = "app-with-iac"
)

// ValidAppTypes contains the valid values for --type flag.
var ValidAppTypes = []AppType{AppTypeAppWithIac}

// AppInitOptions contains options for the app init command.
type AppInitOptions struct {
	AppName             string
	AppType             AppType
	DevAccountID        string
	ProdAccountID       string
	Region              string
	DevEnvironmentName  string
	ProdEnvironmentName string
	DevVarFile          string
	ProdVarFile         string
}

// ValidateAppType checks if the given type string is valid.
func ValidateAppType(t string) error {
	if t == "" {
		return nil
	}

	for _, valid := range ValidAppTypes {
		if AppType(t) == valid {
			return nil
		}
	}

	return fmt.Errorf("invalid type %q, valid values: %s", t, AppTypeAppWithIac)
}

// RunAppInit executes the boilerplate command for app-cicd workflow init.
func RunAppInit(opts AppInitOptions) error {
	if err := checkBoilerplateInstalled(); err != nil {
		return err
	}

	cmd := BuildAppInitCommand(opts)

	if err := cmd.Run(); err != nil {
		errWithMsg := error_user_msg.NewError(
			"failed to initialize app workflow",
			"Check that the template URL is correct and accessible.\nRun 'boilerplate --help' for more information.",
			err,
		)
		return &errWithMsg
	}

	return nil
}

// BuildAppInitCommand constructs the exec.Cmd for boilerplate app-cicd init.
func BuildAppInitCommand(opts AppInitOptions) *exec.Cmd {
	templateURL := buildTemplateURL(TemplateAppCicd)

	args := []string{
		"--template-url", templateURL,
		"--output-folder", ".",
		"--non-interactive",
		"--var", "AppName=" + opts.AppName,
	}

	if opts.AppType == AppTypeAppWithIac {
		args = append(args, "--var", "AppWithIac=true")
	}

	if opts.DevAccountID != "" {
		args = append(args, "--var", "DevAccountId="+opts.DevAccountID)
	}

	if opts.ProdAccountID != "" {
		args = append(args, "--var", "ProdAccountId="+opts.ProdAccountID)
	}

	if opts.Region != "" {
		args = append(args, "--var", "Region="+opts.Region)
	}

	if opts.DevEnvironmentName != "" {
		args = append(args, "--var", "DevEnvironmentName="+opts.DevEnvironmentName)
	}

	if opts.ProdEnvironmentName != "" {
		args = append(args, "--var", "ProdEnvironmentName="+opts.ProdEnvironmentName)
	}

	if opts.DevVarFile != "" {
		args = append(args, "--var-file", opts.DevVarFile)
	}

	if opts.ProdVarFile != "" {
		args = append(args, "--var-file", opts.ProdVarFile)
	}

	cmd := exec.Command("boilerplate", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}
