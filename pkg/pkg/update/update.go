package update

import (
	"context"
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/schema"
	"github.com/oslokommune/ok/pkg/pkg/update/migrate_config"
	"os"
	"strings"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
)

// Run updates the package manifest with the latest releases.
func Run(pkgManifestFilename string, packages []common.Package, opts Options) error {
	manifest, err := common.LoadPackageManifest(pkgManifestFilename)
	if err != nil {
		return fmt.Errorf("loading package manifest: %w", err)
	}

	latestReleases, err := githubreleases.GetLatestReleases()
	if err != nil {
		if strings.Contains(err.Error(), "secret not found in keyring") {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", githubreleases.AuthErrorHelpMessage)
		}

		return fmt.Errorf("failed getting latest github releases: %w", err)
	}

	if !opts.DisableManifestUpdate {
		updatedManifest, err := updatePackages(packages, latestReleases, manifest)
		if err != nil {
			return fmt.Errorf("updating packages: %w", err)
		}

		err = common.SavePackageManifest(pkgManifestFilename, updatedManifest)
		if err != nil {
			return fmt.Errorf("saving package manifest: %w", err)
		}
	}

	if opts.UpdateSchemaConfig {
		err = updateSchemaConfiguration(context.Background(), packages, manifest, latestReleases)
		if err != nil {
			return err
		}
	}

	if opts.MigrateConfig {
		err = migrate_config.MigratePackageConfig(packages)
		if err != nil {
			return fmt.Errorf("migrating package config: %w", err)
		}
	} else {
		fmt.Println("Not migrating package configuration files.")
	}

	common.PrintProcessedPackages(packages, "updated")

	return nil
}

func updatePackages(packagestoUpdate []common.Package, latestReleases map[string]string, manifest common.PackageManifest) (common.PackageManifest, error) {
	updatedManifest := manifest // This works as long as the manifest do not contain any pointers.
	updatedPackages := make([]common.Package, 0, len(packagestoUpdate))

	for _, pkg := range packagestoUpdate {
		latestRelease, ok := latestReleases[pkg.Template] // e.g. v2.1.3
		if !ok {
			return common.PackageManifest{}, fmt.Errorf("no latest release found for package: %s", pkg.Template)
		}

		newRef := fmt.Sprintf("%s-%s", pkg.Template, latestRelease) // e.g. app-v2.1.3

		if pkg.Ref != newRef {
			pkg.Ref = newRef
			updatedPackages = append(updatedPackages, pkg)
		}
	}

	updatedManifest.Packages = updatedPackages

	return updatedManifest, nil
}

func updateSchemaConfiguration(ctx context.Context, packages []common.Package, manifest common.PackageManifest, latestReleases map[string]string) error {
	gh, err := githubreleases.GetGitHubClient()
	if err != nil {
		return fmt.Errorf("getting GitHub client: %w", err)
	}

	fmt.Printf("Updating schema configuration files: ")

	for i, pkg := range packages {
		newRef := fmt.Sprintf("%s-%s", pkg.Template, latestReleases[pkg.Template])
		if manifest.Packages[i].Ref == newRef {
			continue
		}

		manifest.Packages[i].Ref = newRef

		if i > 0 {
			fmt.Printf(", %s", pkg.OutputFolder)
		} else {
			fmt.Printf("%s", pkg.OutputFolder)
		}

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

	fmt.Println()

	return nil
}

func getLastConfigFile(pkg common.Package) (string, bool) {
	if len(pkg.VarFiles) > 0 {
		return pkg.VarFiles[len(pkg.VarFiles)-1], true
	}
	return "", false
}

type Options struct {
	DisableManifestUpdate bool
	MigrateConfig         bool
	UpdateSchemaConfig    bool
}
