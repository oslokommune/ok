package doctor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeRunner answers each command with fixed output. A missing command is ErrNotFound.
type fakeRunner map[string]fakeResult

type fakeResult struct {
	out string
	err error
}

func (f fakeRunner) Run(ctx context.Context, name string, args ...string) (string, error) {
	r, ok := f[name]
	if !ok {
		return "", ErrNotFound
	}
	if r.err != nil && errors.Is(r.err, context.DeadlineExceeded) {
		<-ctx.Done()
		return "", ctx.Err()
	}
	return r.out, r.err
}

// healthy is a machine with every tool at a good version.
func healthy() fakeRunner {
	return fakeRunner{
		"git":         {out: "git version 2.53.0\n"},
		"yq":          {out: "yq (https://github.com/mikefarah/yq/) version v4.53.3\n"},
		"jq":          {out: "jq-1.8.1\n"},
		"gh":          {out: "gh version 2.83.0 (2025-11-04)\nhttps://github.com/cli/cli/releases/tag/v2.83.0\n"},
		"fzf":         {out: "0.74.2 (brew)\n"},
		"terraform":   {out: `{"terraform_version":"1.13.4","platform":"darwin_arm64"}`},
		"aws":         {out: "aws-cli/2.31.0 Python/3.13.7 Darwin/25.0.0 source/arm64\n"},
		"boilerplate": {out: "boilerplate version v0.16.0\n"},
	}
}

func run(t *testing.T, r fakeRunner, okVersion, latest string) Report {
	t.Helper()
	d := Doctor{
		Runner:          r,
		OkVersion:       okVersion,
		LatestOkVersion: func() (string, error) { return latest, nil },
		Timeout:         time.Second,
	}
	return d.Run(context.Background())
}

func find(t *testing.T, r Report, name string) Check {
	t.Helper()
	for _, c := range r.Checks {
		if c.Name == name {
			return c
		}
	}
	t.Fatalf("no check named %q in %+v", name, r.Checks)
	return Check{}
}

func TestHealthyMachinePasses(t *testing.T) {
	report := run(t, healthy(), "5.22.2", "5.22.2")

	assert.True(t, report.OK)
	assert.Empty(t, report.Failed())
	var names []string
	for _, c := range report.Checks {
		names = append(names, c.Name)
		assert.Equal(t, StatusPass, c.Status, c.Name)
		assert.Empty(t, c.Fix, c.Name)
	}
	assert.Equal(t, []string{"git", "yq", "jq", "gh", "fzf", "terraform", "aws", "boilerplate", "ok"}, names)

	assert.Equal(t, "2.53.0", find(t, report, "git").Found)
	assert.Equal(t, "4.53.3", find(t, report, "yq").Found)
	assert.Equal(t, "1.8.1", find(t, report, "jq").Found)
	assert.Equal(t, "0.16.0", find(t, report, "boilerplate").Found)
	assert.Equal(t, "1.13.4", find(t, report, "terraform").Found)
	assert.Equal(t, TerraformConstraint, find(t, report, "terraform").Required)
}

func TestTerraformBounds(t *testing.T) {
	for version, want := range map[string]Status{
		"1.9.2":  StatusFail,
		"1.10.0": StatusPass,
		"1.14.0": StatusPass,
		"2.0.0":  StatusFail,
	} {
		t.Run(version, func(t *testing.T) {
			r := healthy()
			r["terraform"] = fakeResult{out: `{"terraform_version":"` + version + `"}`}
			report := run(t, r, "dev", "")

			c := find(t, report, "terraform")
			assert.Equal(t, want, c.Status)
			assert.Equal(t, version, c.Found)
			assert.Equal(t, TerraformConstraint, c.Required)
			assert.Equal(t, want == StatusPass, report.OK)
			if want == StatusFail {
				assert.Contains(t, c.Fix, "tfswitch")
			}
		})
	}
}

func TestTerraformUnreadableOutputFails(t *testing.T) {
	r := healthy()
	r["terraform"] = fakeResult{out: "Terraform v1.13.4\n"}

	c := find(t, run(t, r, "dev", ""), "terraform")
	assert.Equal(t, StatusFail, c.Status)
	assert.Equal(t, TerraformConstraint, c.Required)
}

func TestYq(t *testing.T) {
	for name, tc := range map[string]struct {
		out    string
		status Status
	}{
		"mikefarah v4": {"yq (https://github.com/mikefarah/yq/) version v4.53.3", StatusPass},
		"mikefarah v3": {"yq version 3.4.1 (https://github.com/mikefarah/yq/)", StatusFail},
		"python yq":    {"yq 3.4.3", StatusFail},
	} {
		t.Run(name, func(t *testing.T) {
			r := healthy()
			r["yq"] = fakeResult{out: tc.out}

			c := find(t, run(t, r, "dev", ""), "yq")
			assert.Equal(t, tc.status, c.Status)
			if tc.status == StatusFail {
				assert.Equal(t, "mikefarah/yq v4", c.Required)
				assert.Equal(t, "Install yq: "+SetupDocs+"yq/", c.Fix)
			}
		})
	}
}

func TestMissingToolFailsWithLinkToItsPage(t *testing.T) {
	r := healthy()
	delete(r, "aws")
	report := run(t, r, "dev", "")

	c := find(t, report, "aws")
	assert.Equal(t, StatusFail, c.Status)
	assert.Equal(t, "not installed", c.Found)
	assert.Equal(t, "Install aws: "+SetupDocs+"aws-cli/", c.Fix)
	assert.False(t, report.OK)
	assert.Len(t, report.Failed(), 1)
}

func TestFailingCommandReportsItsOutput(t *testing.T) {
	r := healthy()
	r["gh"] = fakeResult{out: "dyld: Library not loaded\nmore\n", err: errors.New("exit status 1")}

	c := find(t, run(t, r, "dev", ""), "gh")
	assert.Equal(t, StatusFail, c.Status)
	assert.Equal(t, "`gh --version` failed: dyld: Library not loaded", c.Found)
}

func TestHangingCommandFailsAfterTimeout(t *testing.T) {
	r := healthy()
	r["boilerplate"] = fakeResult{err: context.DeadlineExceeded}
	d := Doctor{Runner: r, OkVersion: "dev", Timeout: 50 * time.Millisecond}

	c := find(t, d.Run(context.Background()), "boilerplate")
	assert.Equal(t, StatusFail, c.Status)
	assert.Equal(t, "`boilerplate --version` did not finish", c.Found)
}

func TestOkVersion(t *testing.T) {
	t.Run("current", func(t *testing.T) {
		c := find(t, run(t, healthy(), "5.22.2", "v5.22.2"), "ok")
		assert.Equal(t, StatusPass, c.Status)
	})
	t.Run("older warns but does not fail", func(t *testing.T) {
		report := run(t, healthy(), "5.21.0", "v5.22.2")
		c := find(t, report, "ok")
		assert.Equal(t, StatusWarn, c.Status)
		assert.Equal(t, "5.22.2", c.Required)
		assert.Contains(t, c.Fix, SetupDocs+"ok/")
		assert.True(t, report.OK)
	})
	t.Run("dev build is not compared", func(t *testing.T) {
		c := find(t, run(t, healthy(), "dev", "v5.22.2"), "ok")
		assert.Equal(t, StatusPass, c.Status)
		assert.Equal(t, "dev", c.Found)
	})
	t.Run("offline warns", func(t *testing.T) {
		d := Doctor{
			Runner:          healthy(),
			OkVersion:       "5.22.2",
			LatestOkVersion: func() (string, error) { return "", errors.New("no network") },
			Timeout:         time.Second,
		}
		report := d.Run(context.Background())
		c := find(t, report, "ok")
		assert.Equal(t, StatusWarn, c.Status)
		assert.Contains(t, c.Fix, "no network")
		assert.True(t, report.OK)
	})
}

func TestJSONShape(t *testing.T) {
	r := healthy()
	r["terraform"] = fakeResult{out: `{"terraform_version":"1.9.2"}`}
	b, err := json.Marshal(run(t, r, "dev", ""))
	require.NoError(t, err)

	var got struct {
		OK     bool `json:"ok"`
		Checks []map[string]any
	}
	require.NoError(t, json.Unmarshal(b, &got))
	assert.False(t, got.OK)
	for _, c := range got.Checks {
		if c["name"] == "terraform" {
			assert.Equal(t, "fail", c["status"])
			assert.Equal(t, "1.9.2", c["found"])
			assert.Equal(t, ">= 1.10.0, < 2.0.0", c["required"])
			assert.NotEmpty(t, c["fix"])
		}
		if c["name"] == "git" {
			assert.NotContains(t, c, "fix", "empty fields are left out")
		}
	}
}

func TestPrint(t *testing.T) {
	r := healthy()
	r["terraform"] = fakeResult{out: `{"terraform_version":"1.9.2"}`}
	var buf bytes.Buffer
	Print(&buf, run(t, r, "dev", ""))

	out := buf.String()
	assert.Contains(t, out, "terraform")
	assert.Contains(t, out, "1.9.2")
	assert.Contains(t, out, "needs "+TerraformConstraint)
	assert.Contains(t, out, "tfswitch")
	assert.Equal(t, 10, strings.Count(out, "\n"), "one line per check, plus the fix under terraform")
}
