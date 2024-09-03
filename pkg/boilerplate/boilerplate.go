package boilerplate

import (
	"context"
	"fmt"
	"path"
	"strings"
)

const (
	DefaultBaseUrl           = "git@github.com:oslokommune/golden-path-boilerplate.git//"
	DefaultPackagePathPrefix = "boilerplate/terraform"
)

type (
	// An interface for a template renderer
	TemplateRenderer interface {
		Render(ctx context.Context, manifestPath string, outputFolders []string, baseUrlOrPath string) error
	}
)

func GetBaseUrlOrDefault(baseUrl string) string {
	if baseUrl == "" {
		return DefaultBaseUrl
	}
	return baseUrl
}

func GetPackagePathPrefixOrDefault(packagePathPrefix string) string {
	if packagePathPrefix == "" {
		return DefaultPackagePathPrefix
	}
	return packagePathPrefix
}

func getValidTemplareUrlOrPath(baseUrlOrPath, templateName, packagePathPrefix, gitRef string) string {
	if !isUrl(baseUrlOrPath) {
		return path.Join(baseUrlOrPath, packagePathPrefix, templateName)
	}

	pathz := strings.Join(
		[]string{packagePathPrefix, templateName},
		"/",
	)
	return fmt.Sprintf("%s%s?ref=%s", baseUrlOrPath, pathz, gitRef)

}
