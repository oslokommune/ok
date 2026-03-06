package workflow

import (
	"os"
	"testing"

	"github.com/oslokommune/ok/pkg/pkg/common"
)

func TestBuildTemplateURL_Default(t *testing.T) {
	// Clear env to use default
	originalEnv := os.Getenv(common.BaseUrlEnvName)
	defer restoreEnv(common.BaseUrlEnvName, originalEnv)
	os.Unsetenv(common.BaseUrlEnvName)

	url := buildTemplateURL(TemplateTerraformIac)
	expected := common.DefaultBaseUrl + "boilerplate/github-actions/terraform-iac?ref=iac-app"

	if url != expected {
		t.Errorf("Expected %s, got %s", expected, url)
	}
}

func TestBuildTemplateURL_EnvVarGitUrl(t *testing.T) {
	originalEnv := os.Getenv(common.BaseUrlEnvName)
	defer restoreEnv(common.BaseUrlEnvName, originalEnv)

	os.Setenv(common.BaseUrlEnvName, "git@github.com:myorg/myrepo.git//")

	url := buildTemplateURL(TemplateAppCicd)
	expected := "git@github.com:myorg/myrepo.git//boilerplate/github-actions/app-cicd?ref=iac-app"

	if url != expected {
		t.Errorf("Expected %s, got %s", expected, url)
	}
}

func TestBuildTemplateURL_EnvVarHttpsUrl(t *testing.T) {
	originalEnv := os.Getenv(common.BaseUrlEnvName)
	defer restoreEnv(common.BaseUrlEnvName, originalEnv)

	os.Setenv(common.BaseUrlEnvName, "https://github.com/myorg/myrepo//")

	url := buildTemplateURL(TemplateTerraformIac)
	expected := "https://github.com/myorg/myrepo//boilerplate/github-actions/terraform-iac?ref=iac-app"

	if url != expected {
		t.Errorf("Expected %s, got %s", expected, url)
	}
}

func TestBuildTemplateURL_LocalPath(t *testing.T) {
	originalEnv := os.Getenv(common.BaseUrlEnvName)
	defer restoreEnv(common.BaseUrlEnvName, originalEnv)

	os.Setenv(common.BaseUrlEnvName, "/tmp/my-boilerplate")

	url := buildTemplateURL(TemplateAppCicd)
	expected := "/tmp/my-boilerplate/boilerplate/github-actions/app-cicd"

	if url != expected {
		t.Errorf("Expected %s, got %s", expected, url)
	}
}

func restoreEnv(key, original string) {
	if original != "" {
		os.Setenv(key, original)
	} else {
		os.Unsetenv(key)
	}
}
