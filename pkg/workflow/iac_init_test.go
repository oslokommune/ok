package workflow

import (
	"os"
	"testing"

	"github.com/oslokommune/ok/pkg/pkg/common"
)

func TestBuildIacInitCommand_NoFlags(t *testing.T) {
	originalEnv := os.Getenv(common.BaseUrlEnvName)
	defer restoreEnv(common.BaseUrlEnvName, originalEnv)
	os.Unsetenv(common.BaseUrlEnvName)

	opts := IacInitOptions{}
	cmd := BuildIacInitCommand(opts)

	expectedArgs := []string{
		"--template-url", common.DefaultBaseUrl + "boilerplate/github-actions/terraform-iac?ref=iac-app",
		"--output-folder", ".",
		"--non-interactive",
	}
	actualArgs := cmd.Args[1:]

	assertArgs(t, expectedArgs, actualArgs)
}

func TestBuildIacInitCommand_AllFlags(t *testing.T) {
	originalEnv := os.Getenv(common.BaseUrlEnvName)
	defer restoreEnv(common.BaseUrlEnvName, originalEnv)
	os.Unsetenv(common.BaseUrlEnvName)

	opts := IacInitOptions{
		AccountID:           "123",
		Region:              "eu-west-1",
		DevEnvironmentName:  "staging",
		ProdEnvironmentName: "production",
		VarFiles:            []string{"common-config.yml"},
	}
	cmd := BuildIacInitCommand(opts)

	expectedArgs := []string{
		"--template-url", common.DefaultBaseUrl + "boilerplate/github-actions/terraform-iac?ref=iac-app",
		"--output-folder", ".",
		"--non-interactive",
		"--var", "AccountId=123",
		"--var", "Region=eu-west-1",
		"--var", "DevEnvironmentName=staging",
		"--var", "ProdEnvironmentName=production",
		"--var-file", "common-config.yml",
	}
	actualArgs := cmd.Args[1:]

	assertArgs(t, expectedArgs, actualArgs)
}

func TestBuildIacInitCommand_MultipleVarFiles(t *testing.T) {
	originalEnv := os.Getenv(common.BaseUrlEnvName)
	defer restoreEnv(common.BaseUrlEnvName, originalEnv)
	os.Unsetenv(common.BaseUrlEnvName)

	opts := IacInitOptions{
		VarFiles: []string{"a.yml", "b.yml"},
	}
	cmd := BuildIacInitCommand(opts)

	expectedArgs := []string{
		"--template-url", common.DefaultBaseUrl + "boilerplate/github-actions/terraform-iac?ref=iac-app",
		"--output-folder", ".",
		"--non-interactive",
		"--var-file", "a.yml",
		"--var-file", "b.yml",
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
