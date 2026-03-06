package workflow

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	okcommon "github.com/oslokommune/ok/pkg/common"
	"github.com/oslokommune/ok/pkg/error_user_msg"
	"github.com/oslokommune/ok/pkg/pkg/common"
)

const (
	BoilerplateGitHubActionsPath = common.BoilerplatePackageGitHubActionsPath
	TemplateRefDefault           = "iac-app"
	TemplateTerraformIac         = "terraform-iac"
	TemplateAppCicd              = "app-cicd"
)

// resolveBaseURL returns the base URL from the environment, or the default.
func resolveBaseURL() string {
	baseURL := os.Getenv(common.BaseUrlEnvName)
	if baseURL == "" {
		return common.DefaultBaseUrl
	}
	return baseURL
}

// buildTemplateURL constructs the full template URL for boilerplate.
// It supports git URLs (git@, http://, https://) and local filesystem paths.
func buildTemplateURL(baseURL, templateName string) string {
	templatePath := strings.Join([]string{BoilerplateGitHubActionsPath, templateName}, "/")

	if okcommon.IsUrl(baseURL) {
		return fmt.Sprintf("%s%s?ref=%s", baseURL, templatePath, TemplateRefDefault)
	}

	return path.Join(baseURL, templatePath)
}

// checkBoilerplateInstalled verifies that the boilerplate CLI is available.
func checkBoilerplateInstalled() error {
	_, err := exec.LookPath("boilerplate")
	if err != nil {
		errWithMsg := error_user_msg.NewError(
			"boilerplate CLI not found",
			"Please install boilerplate:\nhttps://github.com/gruntwork-io/boilerplate",
			err,
		)
		return &errWithMsg
	}
	return nil
}
