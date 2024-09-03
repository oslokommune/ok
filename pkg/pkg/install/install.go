package install

import (
	"context"
	"os"

	"github.com/oslokommune/ok/pkg/boilerplate"
)

const DefaultBaseUrl = "git@github.com:oslokommune/golden-path-boilerplate.git//"
const DefaultPackagePathPrefix = "boilerplate/terraform"

func Run(pkgManifestFilename string, outputFolders []string, useLegacyTemplateRenderer bool) error {
	baseUrlOrPath := os.Getenv("BASE_URL")

	var renderer boilerplate.TemplateRenderer
	if useLegacyTemplateRenderer {
		renderer = boilerplate.NewCommandlineRenderer(os.Stdout, os.Stderr)
	} else {
		renderer = boilerplate.NewBundledRenderer()
	}

	ctx := context.Background()
	return renderer.Render(ctx, pkgManifestFilename, outputFolders, baseUrlOrPath)
}
