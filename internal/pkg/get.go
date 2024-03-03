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

func Get() {

	manifest, err := readPackageManifest("ok.json")
	if err != nil {
		panic(err)
	}

	fs := createFileStore()
	defer fs.Close()

	authClient := createAuthClient()

	ctx := context.Background()
	for _, pkg := range manifest.Packages {
		repo, err := remote.NewRepository(pkg.Repo)
		if err != nil {
			panic(err)
		}
		repo.Client = authClient

		manifestDescriptor, err := oras.Copy(ctx, repo, pkg.Tag, fs, pkg.Name, oras.DefaultCopyOptions)
		if err != nil {
			panic(err)
		}
		fmt.Println("manifest descriptor:", manifestDescriptor)
	}

}

func createFileStore() *file.Store {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fs, err := file.New(dir)
	if err != nil {
		panic(err)
	}

	return fs
}

func createAuthClient() *auth.Client {
	storeOpts := credentials.StoreOptions{}
	credStore, err := credentials.NewStoreFromDocker(storeOpts)
	if err != nil {
		panic(err)
	}

	authClient := &auth.Client{
		Client:     retry.DefaultClient,
		Cache:      auth.NewCache(),
		Credential: credentials.Credential(credStore),
	}

	return authClient
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
