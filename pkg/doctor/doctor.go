// Package doctor checks that the machine has what working with the Golden Path needs: the base tools from
// "Set up your machine", Terraform within the bounds the Golden Path templates set, and a current ok.
package doctor

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Status is the outcome of one check. Only StatusFail makes the report fail.
type Status string

const (
	StatusPass Status = "pass"
	StatusWarn Status = "warn"
	StatusFail Status = "fail"
)

// Check is one line of the report. Name is stable: scripts select checks by it.
type Check struct {
	Name     string `json:"name"`
	Status   Status `json:"status"`
	Found    string `json:"found,omitempty"`
	Required string `json:"required,omitempty"`
	Fix      string `json:"fix,omitempty"`
}

// Report is what `ok doctor --json` prints.
type Report struct {
	OK     bool    `json:"ok"`
	Checks []Check `json:"checks"`
}

// Failed returns the checks with StatusFail.
func (r Report) Failed() []Check {
	var failed []Check
	for _, c := range r.Checks {
		if c.Status == StatusFail {
			failed = append(failed, c)
		}
	}
	return failed
}

// Runner runs external commands. Tests replace it.
type Runner interface {
	// Run runs name with args and returns stdout and stderr combined. It returns ErrNotFound when
	// name is not on PATH.
	Run(ctx context.Context, name string, args ...string) (string, error)
}

// ErrNotFound is returned by a Runner when the command is not on PATH.
var ErrNotFound = errors.New("not found on PATH")

// ExecRunner runs commands with os/exec.
type ExecRunner struct{}

func (ExecRunner) Run(ctx context.Context, name string, args ...string) (string, error) {
	if _, err := exec.LookPath(name); err != nil {
		return "", ErrNotFound
	}
	out, err := exec.CommandContext(ctx, name, args...).CombinedOutput()
	return string(out), err
}

// Doctor runs the checks.
type Doctor struct {
	Runner Runner
	// OkVersion is the running ok's version; "dev" for a local build.
	OkVersion string
	// LatestOkVersion returns the newest released ok version.
	LatestOkVersion func() (string, error)
	// Timeout bounds each check.
	Timeout time.Duration
}

// New returns a Doctor that runs real commands.
func New(okVersion string, latestOkVersion func() (string, error)) Doctor {
	return Doctor{
		Runner:          ExecRunner{},
		OkVersion:       okVersion,
		LatestOkVersion: latestOkVersion,
		Timeout:         15 * time.Second,
	}
}

// Run runs every check concurrently and returns them in a fixed order.
func (d Doctor) Run(ctx context.Context) Report {
	checks := make([]func(context.Context) Check, 0, len(tools)+2)
	for _, t := range tools {
		checks = append(checks, func(ctx context.Context) Check { return d.checkTool(ctx, t) })
		if t.name == "fzf" {
			checks = append(checks, d.checkTerraform)
		}
	}
	checks = append(checks, d.checkOk)

	results := make([]Check, len(checks))
	var wg sync.WaitGroup
	for i, check := range checks {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cctx, cancel := context.WithTimeout(ctx, d.Timeout)
			defer cancel()
			results[i] = check(cctx)
		}()
	}
	wg.Wait()

	report := Report{OK: true, Checks: results}
	for _, c := range results {
		if c.Status == StatusFail {
			report.OK = false
		}
	}
	return report
}

// run runs a command and turns its failure into a failed check, so callers only handle success.
func (d Doctor) run(ctx context.Context, name string, fix string, args ...string) (string, *Check) {
	out, err := d.Runner.Run(ctx, name, args...)
	switch {
	case errors.Is(err, ErrNotFound):
		return "", &Check{Name: name, Status: StatusFail, Found: "not installed", Fix: fix}
	case ctx.Err() != nil:
		return "", &Check{Name: name, Status: StatusFail, Found: fmt.Sprintf("`%s %s` did not finish", name, strings.Join(args, " ")), Fix: fix}
	case err != nil:
		return "", &Check{Name: name, Status: StatusFail, Found: fmt.Sprintf("`%s %s` failed: %s", name, strings.Join(args, " "), firstLine(out, err)), Fix: fix}
	}
	return out, nil
}

func firstLine(out string, err error) string {
	if line, _, _ := strings.Cut(strings.TrimSpace(out), "\n"); line != "" || err == nil {
		return line
	}
	return err.Error()
}
