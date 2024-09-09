package githubreleases

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/v63/github"
)

type GithubRelease struct {
	App     string
	Version *semver.Version
	GitTag  string
}

type Client interface {
	ListGithubReleases(ctx context.Context) ([]GithubRelease, error)
}

func NewGithubReleasesClient(githubClient *github.Client) *GithubReleasesClient {
	return &GithubReleasesClient{
		githubClient: githubClient,
	}
}

type GithubReleasesClient struct {
	githubClient *github.Client
}

// ListGithubReleases implements Client.
func (g *GithubReleasesClient) ListGithubReleases(ctx context.Context) ([]GithubRelease, error) {
	options := &github.ListOptions{PerPage: 105}
	var allReleases []*github.RepositoryRelease
	for {
		releases, res, err := g.githubClient.Repositories.ListReleases(ctx, "oslokommune", "ok", options)
		if err != nil {
			return nil, fmt.Errorf("listing releases: %w", err)
		}
		allReleases = append(allReleases, releases...)
		if res.NextPage == 0 {
			break
		}
		options.Page = res.NextPage
	}

	var githubReleases []GithubRelease
	for _, release := range allReleases {
		appName, appVersion, found := strings.Cut(release.GetTagName(), "-v")
		if !found {
			continue
		}
		version, err := semver.NewVersion("v" + appVersion)
		if err != nil {
			return nil, fmt.Errorf("parsing version: %w", err)
		}
		githubReleases = append(githubReleases, GithubRelease{
			App:     appName,
			Version: version,
			GitTag:  release.GetTagName(),
		})
	}
	return githubReleases, nil
}

var _ Client = (*GithubReleasesClient)(nil)
