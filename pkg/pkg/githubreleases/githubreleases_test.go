package githubreleases

import (
	"testing"

	"github.com/google/go-github/v63/github"
)

func TestSplitComponentAndVersion(t *testing.T) {
	type testcase struct {
		TagName           string
		ExpectedComponent string
		ExpectedVersion   string
	}
	testcases := []testcase{
		{
			TagName:           "app-v6.1.1",
			ExpectedComponent: "app",
			ExpectedVersion:   "v6.1.1",
		},
		{
			TagName:           "app-versions-v5.0.0",
			ExpectedComponent: "app-versions",
			ExpectedVersion:   "v5.0.0",
		},
		{
			TagName:           "load-balancing-stack-v0.0.1",
			ExpectedComponent: "load-balancing-stack",
			ExpectedVersion:   "v0.0.1",
		},
	}

	for _, tc := range testcases {
		releases := []*github.RepositoryRelease{
			{TagName: &tc.TagName},
		}

		components := splitComponentAndVersion(releases)
		if len(components) != 1 {
			t.Errorf("Expected 1 component, got %d", len(components))
		}

		if components[0].Component != tc.ExpectedComponent {
			t.Errorf("Expected component to be %s, got %s", tc.ExpectedComponent, components[0].Component)
		}

		if components[0].Version != tc.ExpectedVersion {
			t.Errorf("Expected version to be %s, got %s", tc.ExpectedVersion, components[0].Version)
		}
	}
}
