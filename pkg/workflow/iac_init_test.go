package workflow

import (
	"testing"

	"github.com/oslokommune/ok/pkg/pkg/common"
)

func TestBuildIacInitCommand_NoFlags(t *testing.T) {
	t.Parallel()

	opts := IacInitOptions{}
	cmd := buildIacInitCommand(common.DefaultBaseUrl, opts)

	expectedArgs := []string{
		"--template-url", common.DefaultBaseUrl + common.BoilerplatePackageGitHubActionsPath + "/terraform-iac",
		"--output-folder", ".",
		"--non-interactive",
	}
	actualArgs := cmd.Args[1:]

	assertArgs(t, expectedArgs, actualArgs)
}

func TestBuildIacInitCommand_AllFlags(t *testing.T) {
	t.Parallel()

	opts := IacInitOptions{
		DevAccountID:        "111111111111",
		ProdAccountID:       "222222222222",
		DevRegion:           "eu-west-1",
		ProdRegion:          "eu-west-1",
		DevEnvironmentName:  "pirates-dev",
		ProdEnvironmentName: "pirates-prod",
	}
	cmd := buildIacInitCommand(common.DefaultBaseUrl, opts)

	expectedArgs := []string{
		"--template-url", common.DefaultBaseUrl + common.BoilerplatePackageGitHubActionsPath + "/terraform-iac",
		"--output-folder", ".",
		"--non-interactive",
		"--var", "DevAccountId=111111111111",
		"--var", "ProdAccountId=222222222222",
		"--var", "DevRegion=eu-west-1",
		"--var", "ProdRegion=eu-west-1",
		"--var", "DevEnvironmentName=pirates-dev",
		"--var", "ProdEnvironmentName=pirates-prod",
	}
	actualArgs := cmd.Args[1:]

	assertArgs(t, expectedArgs, actualArgs)
}

func assertArgs(t *testing.T, expected, actual []string) {
	t.Helper()

	if len(actual) != len(expected) {
		t.Errorf("Expected %d args, got %d\nExpected: %v\nActual:   %v", len(expected), len(actual), expected, actual)
		return
	}

	for i, arg := range expected {
		if actual[i] != arg {
			t.Errorf("Arg %d: expected %q, got %q", i, arg, actual[i])
		}
	}
}
