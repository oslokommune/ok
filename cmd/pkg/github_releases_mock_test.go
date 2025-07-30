package pkg_test

import (
	"context"
	"os"
	"path/filepath"
)

type GitHubReleasesMock struct {
	LatestReleases            map[string]string
	TestWorkingDirectory      string
	BoilerplateRepositoryPath string
}

func (g *GitHubReleasesMock) GetLatestReleases() (map[string]string, error) {
	return g.LatestReleases, nil
}

func (g *GitHubReleasesMock) DownloadGithubFile(ctx context.Context, owner, repo, path, gitRef string) ([]byte, error) {
	testFilePath := filepath.Join(g.TestWorkingDirectory, g.BoilerplateRepositoryPath, inputDir, gitRef, path)

	data, err := os.ReadFile(testFilePath)
	if err != nil {
		return nil, err
	}

	return data, nil
}
