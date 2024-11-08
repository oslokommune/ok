package update

import (
	"context"
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/schema"
	"os"
	"strings"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
)

func Run(pkgManifestFilename string, packages []common.Package, updateSchemaConfig bool) error {
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

	updatedPackages, err := updatePackages(packages, latestReleases)
	if err != nil {
		return fmt.Errorf("updating packages: %w", err)
	}

	err = common.SavePackageManifest(pkgManifestFilename, manifest)
	if err != nil {
		return fmt.Errorf("saving package manifest: %w", err)
	}

	if updateSchemaConfig {
		err = updateSchemaConfig(context.Background(), updatedPackages, manifest, latestReleases)
		if err != nil {
			return err
		}
	}

	return nil
}

func updatePackages(packages []common.Package, latestReleases map[string]string) ([]common.Package, error) {
	updatedPackages := make([]common.Package, 0, len(packages))

	for _, pkg := range packages {
		latestRelease, ok := latestReleases[pkg.Template]
		if !ok {
			return nil, fmt.Errorf("no latest release found for package: %s", pkg.Template)
		}

		newRef := fmt.Sprintf("%s-%s", pkg.Template, latestRelease)

		if pkg.Ref != newRef {
			pkg.Ref = newRef
			updatedPackages = append(updatedPackages, pkg)
		}
	}

	return updatedPackages, nil
}

func updateSchemaConfig(ctx context.Context, updatedPackages []common.Package, manifest common.PackageManifest, latestReleases map[string]string) error {
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
		downloader := githubreleases.NewFileDownloader(gh, common.BoilerplateRepoOwner, common.BoilerplateRepoName, newRef)
		stackPath := githubreleases.GetTemplatePath(manifest.PackagePrefix(), pkg.Template)
		generatedSchema, err := schema.GenerateJsonSchemaForApp(ctx, downloader, stackPath, newRef)
		if err != nil {
			return fmt.Errorf("generating json schema for app: %w", err)
		}

		_, err = schema.CreateOrUpdateConfigurationFile(configFile, newRef, generatedSchema)
		if err != nil {
			return fmt.Errorf("creating or updating configuration file: %w", err)
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
