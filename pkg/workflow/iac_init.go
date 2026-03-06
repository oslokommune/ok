package workflow

import (
	"os"
	"os/exec"

	"github.com/oslokommune/ok/pkg/error_user_msg"
)

// IacInitOptions contains options for the iac init command.
type IacInitOptions struct {
	AccountID           string
	Region              string
	DevEnvironmentName  string
	ProdEnvironmentName string
	VarFiles            []string
}

// RunIacInit executes the boilerplate command for terraform-iac workflow init.
func RunIacInit(opts IacInitOptions) error {
	if err := checkBoilerplateInstalled(); err != nil {
		return err
	}

	cmd := BuildIacInitCommand(opts)

	if err := cmd.Run(); err != nil {
		errWithMsg := error_user_msg.NewError(
			"failed to initialize IAC workflow",
			"Check that the template URL is correct and accessible.\nRun 'boilerplate --help' for more information.",
			err,
		)
		return &errWithMsg
	}

	return nil
}

// BuildIacInitCommand constructs the exec.Cmd for boilerplate terraform-iac init.
func BuildIacInitCommand(opts IacInitOptions) *exec.Cmd {
	templateURL := buildTemplateURL(TemplateTerraformIac)

	args := []string{
		"--template-url", templateURL,
		"--output-folder", ".",
		"--non-interactive",
	}

	if opts.AccountID != "" {
		args = append(args, "--var", "AccountId="+opts.AccountID)
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

	for _, varFile := range opts.VarFiles {
		args = append(args, "--var-file", varFile)
	}

	cmd := exec.Command("boilerplate", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}
