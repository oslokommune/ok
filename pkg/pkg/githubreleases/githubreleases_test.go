package githubreleases

import (
	"testing"

	"github.com/google/go-github/v63/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		require.Len(t, components, 1, "Unexpected number of components")
		assert.Equal(t, tc.ExpectedComponent, components[0].Component, "Unexpected component name")
		assert.Equal(t, tc.ExpectedVersion, components[0].Version, "Unexpected component version")
	}
}