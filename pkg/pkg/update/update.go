package update

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/config"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
)

func Run(pkgManifestFilename string, packageName string, updateConfigSchema bool) error {
	manifest, err := common.LoadPackageManifest(pkgManifestFilename)
	if err != nil {
		return fmt.Errorf("loading package manifest: %w", err)
	}
	latestReleases, err := githubreleases.GetLatestReleases()
	if err != nil {
		if strings.Contains(err.Error(), "secret not found in keyring") {
			fmt.Fprintf(os.Stderr, "%s\n\n", githubreleases.AuthErrorHelpMessage)
		}
		return fmt.Errorf("failed getting latest github releases: %w", err)
	}

	// Set each package to the latest release
	updatedPackages := make([]common.Package, 0, len(manifest.Packages))
	if packageName != "" {
		// Update only the specified package
		updated := false
		for i, pkg := range manifest.Packages {
			if pkg.Template == packageName {
				latestRelease, ok := latestReleases[pkg.Template]
				if !ok {
					return fmt.Errorf("no latest release found for package: %s", packageName)
				}
				newRef := fmt.Sprintf("%s-%s", pkg.Template, latestRelease)
				if manifest.Packages[i].Ref != newRef {
					manifest.Packages[i].Ref = fmt.Sprintf("%s-%s", pkg.Template, latestRelease)
					updatedPackages = append(updatedPackages, manifest.Packages[i])
				}
				updated = true
				break
			}
		}
		if !updated {
			return fmt.Errorf("package not found: %s", packageName)
		}
	} else {
		// Update all packages
		for i, pkg := range manifest.Packages {
			latestRelease, ok := latestReleases[pkg.Template]
			if !ok {
				return fmt.Errorf("no latest release found for package: %s", pkg.Template)
			}
			newRef := fmt.Sprintf("%s-%s", pkg.Template, latestRelease)
			if manifest.Packages[i].Ref != newRef {
				manifest.Packages[i].Ref = fmt.Sprintf("%s-%s", pkg.Template, latestRelease)
				updatedPackages = append(updatedPackages, manifest.Packages[i])
			}
		}
	}

	err = common.SavePackageManifest(pkgManifestFilename, manifest)
	if err != nil {
		return fmt.Errorf("saving package manifest: %w", err)
	}

	if updateConfigSchema {
		ctx := context.Background()
		gh, err := githubreleases.GetGitHubClient()
		if err != nil {
			return fmt.Errorf("getting GitHub client: %w", err)
		}
		for i, pkg := range updatedPackages {
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

	}

	return nil
}

func getLastConfigFile(pkg common.Package) (string, bool) {
	if len(pkg.VarFiles) > 0 {
		return pkg.VarFiles[len(pkg.VarFiles)-1], true
	}
	return "", false
}
