package workflow

import (
	"testing"

	"github.com/oslokommune/ok/pkg/pkg/common"
)

func TestBuildTemplateURL_Default(t *testing.T) {
	t.Parallel()

	url := buildTemplateURL(common.DefaultBaseUrl, TemplateTerraformIac)
	expected := common.DefaultBaseUrl + "boilerplate/github-actions/terraform-iac?ref=iac-app"

	if url != expected {
		t.Errorf("Expected %s, got %s", expected, url)
	}
}

func TestBuildTemplateURL_GitUrl(t *testing.T) {
	t.Parallel()

	url := buildTemplateURL("git@github.com:myorg/myrepo.git//", TemplateAppCicd)
	expected := "git@github.com:myorg/myrepo.git//boilerplate/github-actions/app-cicd?ref=iac-app"

	if url != expected {
		t.Errorf("Expected %s, got %s", expected, url)
	}
}

func TestBuildTemplateURL_HttpsUrl(t *testing.T) {
	t.Parallel()

	url := buildTemplateURL("https://github.com/myorg/myrepo//", TemplateTerraformIac)
	expected := "https://github.com/myorg/myrepo//boilerplate/github-actions/terraform-iac?ref=iac-app"

	if url != expected {
		t.Errorf("Expected %s, got %s", expected, url)
	}
}

func TestBuildTemplateURL_LocalPath(t *testing.T) {
	t.Parallel()

	url := buildTemplateURL("/tmp/my-boilerplate", TemplateAppCicd)
	expected := "/tmp/my-boilerplate/boilerplate/github-actions/app-cicd"

	if url != expected {
		t.Errorf("Expected %s, got %s", expected, url)
	}
}
