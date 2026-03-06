package workflow

import (
	"os"
	"testing"

	"github.com/oslokommune/ok/pkg/pkg/common"
)

func TestBuildAppInitCommand_DefaultType(t *testing.T) {
	originalEnv := os.Getenv(common.BaseUrlEnvName)
	defer restoreEnv(common.BaseUrlEnvName, originalEnv)
	os.Unsetenv(common.BaseUrlEnvName)

	opts := AppInitOptions{
		AppName: "my-app",
	}
	cmd := BuildAppInitCommand(opts)

	expectedArgs := []string{
		"--template-url", common.DefaultBaseUrl + "boilerplate/github-actions/app-cicd?ref=iac-app",
		"--output-folder", ".",
		"--non-interactive",
		"--var", "AppName=my-app",
	}
	actualArgs := cmd.Args[1:]

	assertArgs(t, expectedArgs, actualArgs)
}

func TestBuildAppInitCommand_AppWithIac(t *testing.T) {
	originalEnv := os.Getenv(common.BaseUrlEnvName)
	defer restoreEnv(common.BaseUrlEnvName, originalEnv)
	os.Unsetenv(common.BaseUrlEnvName)

	opts := AppInitOptions{
		AppName:   "my-app",
		AppType:   AppTypeAppWithIac,
		AccountID: "123",
	}
	cmd := BuildAppInitCommand(opts)

	expectedArgs := []string{
		"--template-url", common.DefaultBaseUrl + "boilerplate/github-actions/app-cicd?ref=iac-app",
		"--output-folder", ".",
		"--non-interactive",
		"--var", "AppName=my-app",
		"--var", "AppWithIac=true",
		"--var", "AccountId=123",
	}
	actualArgs := cmd.Args[1:]

	assertArgs(t, expectedArgs, actualArgs)
}

func TestBuildAppInitCommand_AllFlags(t *testing.T) {
	originalEnv := os.Getenv(common.BaseUrlEnvName)
	defer restoreEnv(common.BaseUrlEnvName, originalEnv)
	os.Unsetenv(common.BaseUrlEnvName)

	opts := AppInitOptions{
		AppName:             "my-app",
		AppType:             AppTypeAppWithIac,
		AccountID:           "456",
		Region:              "eu-north-1",
		DevEnvironmentName:  "staging",
		ProdEnvironmentName: "production",
		VarFiles:            []string{"vars.yml"},
	}
	cmd := BuildAppInitCommand(opts)

	expectedArgs := []string{
		"--template-url", common.DefaultBaseUrl + "boilerplate/github-actions/app-cicd?ref=iac-app",
		"--output-folder", ".",
		"--non-interactive",
		"--var", "AppName=my-app",
		"--var", "AppWithIac=true",
		"--var", "AccountId=456",
		"--var", "Region=eu-north-1",
		"--var", "DevEnvironmentName=staging",
		"--var", "ProdEnvironmentName=production",
		"--var-file", "vars.yml",
	}
	actualArgs := cmd.Args[1:]

	assertArgs(t, expectedArgs, actualArgs)
}

func TestValidateAppType_Valid(t *testing.T) {
	if err := ValidateAppType("app-with-iac"); err != nil {
		t.Errorf("Expected no error for valid type, got %v", err)
	}
}

func TestValidateAppType_Empty(t *testing.T) {
	if err := ValidateAppType(""); err != nil {
		t.Errorf("Expected no error for empty type, got %v", err)
	}
}

func TestValidateAppType_Invalid(t *testing.T) {
	err := ValidateAppType("invalid")
	if err == nil {
		t.Error("Expected error for invalid type, got nil")
	}
}
