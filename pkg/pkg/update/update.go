package update

import (
	"context"
	"fmt"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/config"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
)

func Run(pkgManifestFilename string) error {
	ctx := context.Background()
	gh, err := githubreleases.GetGitHubClient()
	if err != nil {
		return fmt.Errorf("getting GitHub client: %w", err)
	}

	manifest, err := common.LoadPackageManifest(pkgManifestFilename)
	if err != nil {
		return fmt.Errorf("loading package manifest: %w", err)
	}

	latestReleases, err := githubreleases.GetLatestReleases()
	if err != nil {
		return fmt.Errorf("getting latest releases: %w", err)
	}

	// Set each package to the latest release
	for i, pkg := range manifest.Packages {
		newRef := fmt.Sprintf("%s-%s", pkg.Template, latestReleases[pkg.Template])
		manifest.Packages[i].Ref = newRef

		configFile, ok := getLastConfigFile(pkg)
		if !ok {
			continue
		}
		downloader := githubreleases.NewFileDownloader(gh, githubreleases.GithubOwner, githubreleases.GithubRepo, newRef)
		stackPath := githubreleases.GetTemplatePath(pkg.Template)
		schema, err := config.GenerateJsonSchemaForApp(ctx, downloader, stackPath, newRef)
		if err != nil {
			return fmt.Errorf("generating json schema for app: %w", err)
		}

		_, err = config.CreateOrUpdateConfigurationFile(configFile, newRef, schema)
		if err != nil {
			return fmt.Errorf("creating or updating configuration file: %w", err)
		}

	}

	err = common.SavePackageManifest(pkgManifestFilename, manifest)
	if err != nil {
		return err
	}

	return nil
}

func getLastConfigFile(pkg common.Package) (string, bool) {
	if len(pkg.VarFiles) > 0 {
		return pkg.VarFiles[len(pkg.VarFiles)-1], true
	}
	return "", false
}
