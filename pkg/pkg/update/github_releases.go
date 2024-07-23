package update

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/v63/github"
	"github.com/zalando/go-keyring"
)

func getLatestReleases() (map[string]string, error) {

	authToken, err := getGitHubToken()
	if err != nil {
		return nil, fmt.Errorf("getting GitHub token: %w", err)
	}

	client := github.NewClient(nil).WithAuthToken(authToken)

	githubReleases, err := listReleases(client)
	if err != nil {
		return nil, fmt.Errorf("listing releases: %w", err)
	}

	releases := splitComponentAndVersion(githubReleases)

	latestReleases, err := parseLatestReleases(releases)
	if err != nil {
		return nil, fmt.Errorf("getting latest releases: %w", err)
	}
	return latestReleases, nil
}

func getGitHubToken() (string, error) {
	if token := os.Getenv("GH_TOKEN"); token != "" {
		return token, nil
	}

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token, nil
	}

	token, err := keyring.Get("gh:github.com", "")
	if err != nil {
		return "", fmt.Errorf("getting GitHub token from keyring: %w", err)
	}

	return token, nil
}

func listReleases(client *github.Client) ([]*github.RepositoryRelease, error) {
	options := &github.ListOptions{PerPage: 105}

	var allReleases []*github.RepositoryRelease

	for {

		releases, response, err := client.Repositories.ListReleases(
			context.Background(), "oslokommune", "golden-path-boilerplate", options)

		if err != nil {
			return nil, fmt.Errorf("listing releases: %w", err)
		}

		allReleases = append(allReleases, releases...)

		if response.NextPage == 0 {
			break
		}

		options.Page = response.NextPage

	}

	return allReleases, nil
}

func splitComponentAndVersion(releases []*github.RepositoryRelease) []Release {
	var components []Release

	for _, release := range releases {
		tagName := *release.TagName
		lastDashIndex := strings.LastIndex(tagName, "-")

		if lastDashIndex > -1 {
			component := tagName[:lastDashIndex]
			version := tagName[lastDashIndex+1:]
			components = append(components, Release{Component: component, Version: version})
		}
	}

	return components
}

func parseLatestReleases(components []Release) (map[string]string, error) {
	latestComponents := make(map[string]string)

	for _, component := range components {
		latestVersion, keyFound := latestComponents[component.Component]
		if !keyFound {
			// Set version if not found
			latestComponents[component.Component] = component.Version
			continue
		}

		// Convert to semver
		latestVersionSemver, err := semver.NewVersion(latestVersion)
		if err != nil {
			return nil, fmt.Errorf("parsing version string '%s': %w", latestVersion, err)
		}

		componentVersionSemver, err := semver.NewVersion(component.Version)
		if err != nil {
			return nil, fmt.Errorf("parsing version string '%s': %w", component.Version, err)
		}

		if componentVersionSemver.GreaterThan(latestVersionSemver) {
			latestComponents[component.Component] = component.Version
		}
	}
	return latestComponents, nil
}
