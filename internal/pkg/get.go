package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/credentials"
	"oras.land/oras-go/v2/registry/remote/retry"
)

type PackageManifest struct {
	Packages []Package `json:"packages"`
}

type Package struct {
	Name string `json:"name"`
	Repo string `json:"repo"`
	Tag  string `json:"tag"`
}

func Get() error {
	manifest, err := readPackageManifest("ok.json")
	if err != nil {
		return fmt.Errorf("failed to read package manifest: %w", err)
	}

	fs, err := createFileStore()
	if err != nil {
		return fmt.Errorf("failed to create file store: %w", err)
	}
	defer fs.Close()

	authClient, err := createAuthClient()
	if err != nil {
		return fmt.Errorf("failed to create auth client: %w", err)
	}

	ctx := context.Background()
	for _, pkg := range manifest.Packages {
		err := processPackage(ctx, pkg, authClient, fs)
		if err != nil {
			return fmt.Errorf("failed to process package: %w", err)
		}
	}

	return nil
}

func createFileStore() (*file.Store, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	fs, err := file.New(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to create new file: %w", err)
	}

	return fs, nil
}

func createAuthClient() (*auth.Client, error) {
	storeOpts := credentials.StoreOptions{}
	credStore, err := credentials.NewStoreFromDocker(storeOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create new store from docker: %w", err)
	}

	authClient := &auth.Client{
		Client:     retry.DefaultClient,
		Cache:      auth.NewCache(),
		Credential: credentials.Credential(credStore),
	}

	return authClient, nil
}

func readPackageManifest(filePath string) (*PackageManifest, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	var manifest PackageManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return &manifest, nil
}

func processPackage(ctx context.Context, pkg Package, authClient *auth.Client, fs *file.Store) error {
	repo, err := remote.NewRepository(pkg.Repo)
	if err != nil {
		return fmt.Errorf("failed to create new repository: %w", err)
	}
	repo.Client = authClient

	_, err = oras.Copy(ctx, repo, pkg.Tag, fs, pkg.Name, oras.DefaultCopyOptions)
	if err != nil {
		return fmt.Errorf("failed to copy package: %w", err)
	}

	return nil
}
