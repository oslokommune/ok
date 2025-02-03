package update

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/oslokommune/ok/pkg/pkg/schema"
	"github.com/oslokommune/ok/pkg/pkg/update/migrate_config"
	"github.com/oslokommune/ok/pkg/pkg/update/migrate_config/metadata"
	"os"
	"strings"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
)

type Updater struct {
	ghReleases GitHubReleases
}

type GitHubReleases interface {
	GetLatestReleases() (map[string]string, error)
}

func NewUpdater(ghReleases GitHubReleases) Updater {
	return Updater{
		ghReleases: ghReleases,
	}
}

// Run updates the package manifest with the latest releases.
func (u Updater) Run(pkgManifestFilename string, selectedPackagesInput []common.Package, opts Options) error {
	var selectedPackages []common.Package
	var manifest common.PackageManifest
	{
		currentManifest, err := common.LoadPackageManifest(pkgManifestFilename)
		if err != nil {
			return fmt.Errorf("loading package manifest: %w", err)
		}

		if opts.DisableManifestUpdate {
			manifest = currentManifest
			selectedPackages = selectedPackagesInput
		} else {
			updateManifest, updatedSelectedPackages, err := u.updateManifest(currentManifest, selectedPackagesInput)
			if err != nil {
				return fmt.Errorf("updating package manifest: %w", err)
			}

			manifest = updateManifest
			selectedPackages = updatedSelectedPackages

			err = common.SavePackageManifest(pkgManifestFilename, manifest)
			if err != nil {
				return fmt.Errorf("saving package manifest: %w", err)
			}
		}
	}

	if opts.UpdateSchemaConfig {
		err := updateSchemaConfiguration(context.Background(), selectedPackages, manifest.PackagePrefix())
		if err != nil {
			return err
		}
	}

	if opts.MigrateConfig {
		err := migrate_config.MigratePackageConfig(selectedPackages)
		if err != nil {
			return fmt.Errorf("migrating package config: %w", err)
		}
	} else {
		fmt.Println("Not migrating package configuration files.")
	}

	common.PrintProcessedPackages(selectedPackages, "updated")

	return nil
}

// updateManifest updates package versions in the manifest with the latest releases from GitHub
func (u Updater) updateManifest(manifest common.PackageManifest, selectedPackages []common.Package) (common.PackageManifest, []common.Package, error) {
	latestReleases, err := u.ghReleases.GetLatestReleases()
	if err != nil {
		if strings.Contains(err.Error(), "secret not found in keyring") {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", githubreleases.AuthErrorHelpMessage)
		}

		return common.PackageManifest{}, nil, fmt.Errorf("failed getting latest github releases: %w", err)
	}

	updatedManifest, updatedPackages, err := updatePackages(manifest, selectedPackages, latestReleases)
	if err != nil {
		return common.PackageManifest{}, nil, fmt.Errorf("updating packages: %w", err)
	}

	return updatedManifest, updatedPackages, nil
}

// updatePackages updates the package manifest with the latest releases. It only updates the packages found in selectedPackages.
func updatePackages(manifest common.PackageManifest, selectedPackages []common.Package, latestReleases map[string]string) (common.PackageManifest, []common.Package, error) {
	updatedManifest := manifest.Clone()
	updatedPackages := make([]common.Package, 0)

	for i, _ := range manifest.Packages {
		pkg := &updatedManifest.Packages[i]

		_, err := pkg.PackageVersion()
		if errors.Is(err, semver.ErrInvalidSemVer) {
			// pkg.Ref is not a semver version, e.g. "main". We don't update these.
			continue
		}

		if !common.ContainsPackage(selectedPackages, *pkg) {
			continue
		}

		// pkg is a package that is in selectedPackages, i.e. it should be updated
		latestRelease, ok := latestReleases[pkg.Template] // e.g. v2.1.3
		if !ok {
			return common.PackageManifest{}, nil, fmt.Errorf("no latest release found for package: %s", pkg.Template)
		}

		newRef := fmt.Sprintf("%s-%s", pkg.Template, latestRelease) // e.g. app-v2.1.3
		pkg.Ref = newRef
		updatedPackages = append(updatedPackages, *pkg)
	}

	return updatedManifest, updatedPackages, nil
}

// updateSchemaConfiguration does two things. For each package in the package manifest, that is also in selectedPackages:
// 1) Download the JSON schema file for each template. The version download is the one found in the package manifest.
// 2) Update the stack configuration file header with the downloaded JSON schema. For instance: "# yaml-language-server: $schema=.schemas/app-v8.0.5.schema.json"
func updateSchemaConfiguration(ctx context.Context, selectedPackages []common.Package, manifestPackagePrefix string) error {
	gh, err := githubreleases.GetGitHubClient()
	if err != nil {
		return fmt.Errorf("getting GitHub client: %w", err)
	}

	fmt.Println("Updating json schemas:")

	for _, pkg := range selectedPackages {
		_, err := pkg.PackageVersion()
		if err != nil {
			// pkg.Ref might be "main". To keep code simple, we avoid dealing with non-semver versions.
			continue
		}

		varFile, ok := getLastVarFile(pkg)
		if !ok {
			continue
		}

		fmt.Printf("- %s\n", pkg.OutputFolder)

		// Get current JSON schema for the package
		jsonSchema, err := metadata.ParseFirstLine(varFile)
		if err != nil && errors.Is(err, metadata.ErrMissingSchemaDeclaration) {
			// Proceeed with downloading JSON schema and updating the varFile, so that the JSON schema declaration is
			// added to the varFile. The next time this code is run, the schema declaration will then be found.
		} else if err != nil {
			return fmt.Errorf("parsing first line of file '%s': %w", varFile, err)
		}

		existingRef := fmt.Sprintf("%s-%s", jsonSchema.Template, jsonSchema.Version)
		if existingRef == pkg.Ref {
			// No need to update the varFile with a new JSON schema, as the existing one is as declared in the pacckage
			// manifest.
			continue
		}

		// Update the JSON schema, i.e. download it and update the varFile's schema declaration to point to it.
		downloader := githubreleases.NewFileDownloader(gh, common.BoilerplateRepoOwner, common.BoilerplateRepoName, pkg.Ref)
		stackPath := githubreleases.GetTemplatePath(manifestPackagePrefix, pkg.Template)

		generatedSchema, err := schema.GenerateJsonSchemaForApp(ctx, downloader, stackPath, pkg.Ref)
		if err != nil {
			return fmt.Errorf("generating json schema for app: %w", err)
		}

		_, err = schema.CreateOrUpdateVarFile(varFile, pkg.Ref, generatedSchema)
		if err != nil {
			return fmt.Errorf("creating or updating configuration file: %w", err)
		}
	}

	return nil
}

func getLastVarFile(pkg common.Package) (string, bool) {
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
