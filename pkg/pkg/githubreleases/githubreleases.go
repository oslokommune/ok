package githubreleases

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/oslokommune/ok/pkg/pkg/common"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/v63/github"
	"github.com/zalando/go-keyring"
	"gopkg.in/yaml.v3"
)

const AuthErrorHelpMessage = `
GitHub token not found in keyring or environment variables.

Steps to resolve:

1. Ensure you have the latest GitHub CLI version, which in most cases should store the token in the OS keyring:
   https://cli.github.com/

2. Try logging in again with GitHub CLI:
   gh auth login

3. If you're still encountering issues, you can bypass the keyring by setting the token as an environment variable:
   export GH_TOKEN=$(gh auth token)`

type GitHubReleasesImpl struct {
}

func NewGitHubReleases() GitHubReleasesImpl {
	return GitHubReleasesImpl{}
}

func (g GitHubReleasesImpl) GetLatestReleases() (map[string]string, error) {
	return GetLatestReleases()
}

type Release struct {
	Component string
	Version   string
}

// GhConfig represents the structure of the hosts.yml file
type GhConfig struct {
	Github struct {
		OAuthToken string `yaml:"oauth_token,omitempty"`
	} `yaml:"github.com,omitempty"`
}

// getGitHubTokenFromHostsFile retrieves the GitHub OAuth token from the gh CLI hosts.yml file
func getGitHubTokenFromHostsFile() (string, error) {
	configPath, err := getGhConfigPath()
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", err
	}

	var config GhConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return "", err
	}

	return config.Github.OAuthToken, nil
}

func GetLatestOkVersion() (*semver.Version, error) {
	authToken, err := getGitHubToken()
	if err != nil {
		return nil, fmt.Errorf("getting GitHub token: %w", err)
	}
	client := github.NewClient(nil).WithAuthToken(authToken)

	release, _, err := client.Repositories.GetLatestRelease(context.Background(), "oslokommune", "ok")
	if err != nil {
		return nil, fmt.Errorf("error getting latest release: %v", err)
	}

	versionTag := release.GetTagName()
	versionSemver, err := semver.NewVersion(versionTag)
	if err != nil {
		return nil, fmt.Errorf("parsing version string '%s': %w", versionTag, err)
	}

	return versionSemver, nil
}

func GetLatestReleases() (map[string]string, error) {

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

func GetGitHubToken() (string, error) {
	return getGitHubToken()
}

func GetGitHubClient() (*github.Client, error) {
	token, err := getGitHubToken()
	if err != nil {
		return nil, fmt.Errorf("getting GitHub token: %w", err)
	}

	return github.NewClient(nil).WithAuthToken(token), nil
}

// getGHConfigPath returns the path to the gh CLI hosts.yml file
func getGhConfigPath() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return "", fmt.Errorf("$HOME environment variable not set")
	}

	return filepath.Join(home, ".config", "gh", "hosts.yml"), nil
}

func getGitHubToken() (string, error) {
	// Try environment variables first
	if token := os.Getenv("GH_TOKEN"); token != "" {
		return token, nil
	}

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token, nil
	}

	// Try system keyring
	token, err := keyring.Get("gh:github.com", "")
	if err == nil {
		return token, nil
	}

	// Try hosts.yml file
	token, err = getGitHubTokenFromHostsFile()
	if err == nil {
		return token, nil
	}

	return "", err
}

func listReleases(client *github.Client) ([]*github.RepositoryRelease, error) {
	options := &github.ListOptions{PerPage: 105}
	var allReleases []*github.RepositoryRelease

	for {
		releases, response, err := client.Repositories.ListReleases(
			context.Background(), common.BoilerplateRepoOwner, common.BoilerplateRepoName, options)
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

var reNamespaceSemver = regexp.MustCompile(`(?m)^(?:(?P<namespace>[a-zA-Z][a-zA-Z0-9\-_]*?)(?:-?v?))?(?P<version>(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?)$`)

func splitComponentAndVersion(releases []*github.RepositoryRelease) []Release {
	var reGroups = reNamespaceSemver.SubexpNames()
	var components []Release

	for _, release := range releases {
		if release.TagName == nil {
			continue
		}

		tagName := *release.TagName
		matches := reNamespaceSemver.FindStringSubmatch(tagName)
		var namespace, version string
		for groupID, group := range matches {
			switch reGroups[groupID] {
			case "namespace":
				namespace = group
			case "version":
				version = group
			}
		}

		if version == "" {
			continue
		} else if v, err := semver.NewVersion(version); err != nil {
			continue
		} else {
			version = "v" + v.String()
		}

		components = append(components, Release{Component: namespace, Version: version})
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

		if componentVersionSemver.Prerelease() != "" {
			// Skip prerelease versions
			continue
		}

		if componentVersionSemver.GreaterThan(latestVersionSemver) {
			latestComponents[component.Component] = component.Version
		}
	}
	return latestComponents, nil
}

func DownloadGithubFile(ctx context.Context, client *github.Client, owner, repo, path, ref string) ([]byte, error) {
	rc, response, err := client.Repositories.DownloadContents(ctx, owner, repo, path, &github.RepositoryContentGetOptions{Ref: ref})
	if err != nil {
		return nil, fmt.Errorf("downloading file: %w", err)
	}

	// check if response headers contains "application/vnd.github.raw"
	if !strings.Contains(response.Header.Get("Content-Type"), "application/vnd.github.raw") {
		return nil, fmt.Errorf("response content type is not application/vnd.github.raw")
	}

	defer rc.Close()

	bodyText, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	return bodyText, nil
}

func GetTemplatePath(stackPath string, app string) string {
	return fmt.Sprintf("%s/%s", stackPath, app)
}
