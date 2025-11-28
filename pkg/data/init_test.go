package data

import (
	"os"
	"testing"
)

func TestDetermineTemplate_UserProvided(t *testing.T) {
	// User template should take highest priority
	result := determineTemplate("user/custom-template")
	expected := "user/custom-template"

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestDetermineTemplate_EnvVar(t *testing.T) {
	// Save original env var and restore after test
	originalEnv := os.Getenv(TemplateURLEnvName)
	defer func() {
		if originalEnv != "" {
			os.Setenv(TemplateURLEnvName, originalEnv)
		} else {
			os.Unsetenv(TemplateURLEnvName)
		}
	}()

	// Set environment variable
	envTemplate := "env/custom-template"
	os.Setenv(TemplateURLEnvName, envTemplate)

	// When no user template, should use env var
	result := determineTemplate("")
	if result != envTemplate {
		t.Errorf("Expected %s, got %s", envTemplate, result)
	}
}

func TestDetermineTemplate_Default(t *testing.T) {
	// Save original env var and restore after test
	originalEnv := os.Getenv(TemplateURLEnvName)
	defer func() {
		if originalEnv != "" {
			os.Setenv(TemplateURLEnvName, originalEnv)
		} else {
			os.Unsetenv(TemplateURLEnvName)
		}
	}()

	// Clear environment variable
	os.Unsetenv(TemplateURLEnvName)

	// When no user template and no env var, should use default
	result := determineTemplate("")
	if result != DefaultTemplateURL {
		t.Errorf("Expected %s, got %s", DefaultTemplateURL, result)
	}
}

func TestBuildDatabricksCommand_MinimalArgs(t *testing.T) {
	opts := InitOptions{
		TemplatePath: "test-template",
	}

	cmd := buildDatabricksCommand(opts)

	// Check args (Args[0] is the full path to the command, so we check the rest)
	expectedArgs := []string{"bundle", "init", "test-template"}
	actualArgs := cmd.Args[1:] // Skip first arg which is the command path

	if len(actualArgs) != len(expectedArgs) {
		t.Errorf("Expected %d args, got %d", len(expectedArgs), len(actualArgs))
	}

	for i, arg := range expectedArgs {
		if actualArgs[i] != arg {
			t.Errorf("Arg %d: expected '%s', got '%s'", i, arg, actualArgs[i])
		}
	}
}

func TestBuildDatabricksCommand_AllFlags(t *testing.T) {
	opts := InitOptions{
		TemplatePath: "test-template",
		Branch:       "main",
		ConfigFile:   "config.json",
		OutputDir:    "./output",
		Tag:          "v1.0.0",
		TemplateDir:  "templates",
	}

	cmd := buildDatabricksCommand(opts)

	// Check that all flags are present in args (skip first arg which is command path)
	expectedArgs := []string{
		"bundle",
		"init",
		"test-template",
		"--branch", "main",
		"--config-file", "config.json",
		"--output-dir", "./output",
		"--tag", "v1.0.0",
		"--template-dir", "templates",
	}
	actualArgs := cmd.Args[1:]

	if len(actualArgs) != len(expectedArgs) {
		t.Errorf("Expected %d args, got %d", len(expectedArgs), len(actualArgs))
	}

	for i, arg := range expectedArgs {
		if actualArgs[i] != arg {
			t.Errorf("Arg %d: expected '%s', got '%s'", i, arg, actualArgs[i])
		}
	}
}

func TestBuildDatabricksCommand_NoTemplate(t *testing.T) {
	opts := InitOptions{}

	cmd := buildDatabricksCommand(opts)

	// Check args when no template is provided (skip first arg which is command path)
	expectedArgs := []string{"bundle", "init"}
	actualArgs := cmd.Args[1:]

	if len(actualArgs) != len(expectedArgs) {
		t.Errorf("Expected %d args, got %d", len(expectedArgs), len(actualArgs))
	}

	for i, arg := range expectedArgs {
		if actualArgs[i] != arg {
			t.Errorf("Arg %d: expected '%s', got '%s'", i, arg, actualArgs[i])
		}
	}
}
