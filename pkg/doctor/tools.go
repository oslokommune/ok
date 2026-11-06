package doctor

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"
)

// SetupDocs is the "Set up your machine" guide. Each tool's fix links to its page there.
const SetupDocs = "https://km.oslo.systems/set-up-your-machine/"

// TerraformConstraint must match the default of TerraformVersion in golden-path-boilerplate's
// boilerplate/terraform/versions/boilerplate.yml, which every Golden Path stack is generated with.
const TerraformConstraint = ">= 1.10.0, < 2.0.0"

// tool is a base tool: the required tools on "Set up your machine", plus git, which that guide assumes.
type tool struct {
	name    string
	args    []string
	install string
	// check validates the version output and returns what was found, and, when the tool is not usable,
	// what is required.
	check func(out string) (found, required string, ok bool)
}

var tools = []tool{
	{name: "git", args: []string{"--version"}, install: "https://git-scm.com/downloads", check: anyVersion},
	{name: "yq", args: []string{"--version"}, install: SetupDocs + "yq/", check: checkYq},
	{name: "jq", args: []string{"--version"}, install: SetupDocs + "jq/", check: anyVersion},
	{name: "gh", args: []string{"--version"}, install: SetupDocs + "gh/", check: anyVersion},
	{name: "fzf", args: []string{"--version"}, install: SetupDocs + "fzf/", check: anyVersion},
	{name: "aws", args: []string{"--version"}, install: SetupDocs + "aws-cli/", check: anyVersion},
	{name: "boilerplate", args: []string{"--version"}, install: SetupDocs + "boilerplate/", check: anyVersion},
}

var versionPattern = regexp.MustCompile(`\d+\.\d+(\.\d+)?`)

// anyVersion accepts any output and reports the first version number in it.
func anyVersion(out string) (string, string, bool) {
	if v := versionPattern.FindString(out); v != "" {
		return v, "", true
	}
	return strings.TrimSpace(firstLine(out, nil)), "", true
}

// checkYq requires mikefarah's yq v4. The Python yq, a jq wrapper with the same name, fails on ok's yq
// expressions.
func checkYq(out string) (string, string, bool) {
	const required = "mikefarah/yq v4"
	v := versionPattern.FindString(out)
	if !strings.Contains(out, "mikefarah") {
		return strings.TrimSpace("yq " + v + " (not mikefarah/yq)"), required, false
	}
	if !strings.HasPrefix(v, "4.") {
		return v, required, false
	}
	return v, "", true
}

func (d Doctor) checkTool(ctx context.Context, t tool) Check {
	fix := "Install " + t.name + ": " + t.install
	out, failed := d.run(ctx, t.name, fix, t.args...)
	if failed != nil {
		return *failed
	}
	found, required, ok := t.check(out)
	if !ok {
		return Check{Name: t.name, Status: StatusFail, Found: found, Required: required, Fix: fix}
	}
	return Check{Name: t.name, Status: StatusPass, Found: found}
}

func (d Doctor) checkTerraform(ctx context.Context) Check {
	fix := fmt.Sprintf("Install a Terraform version that satisfies %s, for example with tfswitch: %sterraform/",
		TerraformConstraint, SetupDocs)
	out, failed := d.run(ctx, "terraform", fix, "version", "-json")
	if failed != nil {
		failed.Required = TerraformConstraint
		return *failed
	}
	var v struct {
		TerraformVersion string `json:"terraform_version"`
	}
	if err := json.Unmarshal([]byte(out), &v); err != nil || v.TerraformVersion == "" {
		return Check{Name: "terraform", Status: StatusFail, Found: "unreadable `terraform version -json`",
			Required: TerraformConstraint, Fix: fix}
	}
	ok, err := satisfies(v.TerraformVersion, TerraformConstraint)
	if err != nil || !ok {
		return Check{Name: "terraform", Status: StatusFail, Found: v.TerraformVersion, Required: TerraformConstraint, Fix: fix}
	}
	return Check{Name: "terraform", Status: StatusPass, Found: v.TerraformVersion, Required: TerraformConstraint}
}

func satisfies(version, constraint string) (bool, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return false, err
	}
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		return false, err
	}
	return c.Check(v), nil
}

// checkOk warns when a newer ok is released. It never fails: an old ok still works, and the check needs
// the network.
func (d Doctor) checkOk(ctx context.Context) Check {
	c := Check{Name: "ok", Status: StatusPass, Found: d.OkVersion}
	if d.OkVersion == "dev" || d.LatestOkVersion == nil {
		return c
	}
	current, err := semver.NewVersion(d.OkVersion)
	if err != nil {
		return c
	}

	type result struct {
		version string
		err     error
	}
	done := make(chan result, 1)
	go func() {
		v, err := d.LatestOkVersion()
		done <- result{v, err}
	}()
	var r result
	select {
	case r = <-done:
	case <-ctx.Done():
		r.err = ctx.Err()
	}
	if r.err != nil {
		c.Status = StatusWarn
		c.Fix = "Could not get the latest ok release (" + r.err.Error() + "); run `ok version` later"
		return c
	}

	latest, err := semver.NewVersion(r.version)
	if err != nil {
		return c
	}
	if current.LessThan(latest) {
		c.Status = StatusWarn
		c.Required = latest.String()
		c.Fix = "Update ok to " + latest.String() + ": " + SetupDocs + "ok/"
	}
	return c
}
