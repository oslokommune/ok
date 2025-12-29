package pk

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/oslokommune/ok/pkg/pk"
)

// setupTestRepo copies testdata to a temp dir and initializes it as a git repo.
func setupTestRepo(t *testing.T, testdataName string) string {
	t.Helper()

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "pk-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Copy testdata to temp dir
	testdataDir := filepath.Join("testdata", testdataName)
	if err := copyDir(testdataDir, tmpDir); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("failed to copy testdata: %v", err)
	}

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("failed to init git repo: %v", err)
	}

	// Configure git user for the test repo
	cmd = exec.Command("git", "config", "user.email", "test@test.com")
	cmd.Dir = tmpDir
	cmd.Run()
	cmd = exec.Command("git", "config", "user.name", "Test")
	cmd.Dir = tmpDir
	cmd.Run()

	// Add and commit so git rev-parse works
	cmd = exec.Command("git", "add", ".")
	cmd.Dir = tmpDir
	cmd.Run()
	cmd = exec.Command("git", "commit", "-m", "initial")
	cmd.Dir = tmpDir
	cmd.Run()

	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	return tmpDir
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(dstPath, data, info.Mode())
	})
}

func TestInstallContextAware_FiltersBySubfolder(t *testing.T) {
	repoDir := setupTestRepo(t, "context-aware")

	tests := []struct {
		name            string
		cwd             string
		expectedMatches []string
	}{
		{
			name:            "exact match on app-hello",
			cwd:             "app-hello",
			expectedMatches: []string{"app-hello"},
		},
		{
			name:            "within app-hello subdir",
			cwd:             "app-hello/src",
			expectedMatches: []string{"app-hello"},
		},
		{
			name:            "exact match on networking",
			cwd:             "networking",
			expectedMatches: []string{"networking"},
		},
		{
			name:            "at repo root returns no matches",
			cwd:             ".",
			expectedMatches: nil,
		},
		{
			name:            "unrelated dir returns no matches",
			cwd:             "other",
			expectedMatches: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the cwd if it doesn't exist
			cwdPath := filepath.Join(repoDir, tt.cwd)
			os.MkdirAll(cwdPath, 0755)

			// Load configs
			okDir := filepath.Join(repoDir, ".ok")
			configs, err := pk.LoadConfigs(okDir)
			if err != nil {
				t.Fatalf("failed to load configs: %v", err)
			}

			templates, err := pk.ApplyCommon(configs)
			if err != nil {
				t.Fatalf("failed to apply common: %v", err)
			}

			// Filter by working directory
			matched := pk.FilterTemplatesByWorkingDir(templates, cwdPath, repoDir)

			if tt.expectedMatches == nil {
				if matched != nil && len(matched) > 0 {
					t.Errorf("expected no matches, got %d", len(matched))
				}
				return
			}

			if len(matched) != len(tt.expectedMatches) {
				t.Errorf("expected %d matches, got %d", len(tt.expectedMatches), len(matched))
				return
			}

			for i, tpl := range matched {
				expectedSubfolder := tt.expectedMatches[i]
				actualOutput := filepath.Join(tpl.BaseOutputFolder, tpl.Subfolder)
				if filepath.Clean(actualOutput) != expectedSubfolder {
					t.Errorf("expected subfolder %q, got %q", expectedSubfolder, actualOutput)
				}
			}
		})
	}
}

func TestInstallCommand_DryRunWithAll(t *testing.T) {
	repoDir := setupTestRepo(t, "context-aware")

	// Change to repo dir
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(repoDir)

	cmd := NewInstallCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--dry-run", "--all"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	output := out.String()

	// Should have 2 dry-run lines (one per template)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 dry-run lines, got %d: %s", len(lines), output)
	}

	// Verify all templates are present
	if !strings.Contains(output, "app-hello") {
		t.Error("expected app-hello in output")
	}
	if !strings.Contains(output, "networking") {
		t.Error("expected networking in output")
	}
}

func TestInstallCommand_DryRunFromSubfolder(t *testing.T) {
	repoDir := setupTestRepo(t, "context-aware")

	// Change to app-hello subfolder
	appHelloDir := filepath.Join(repoDir, "app-hello")
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(appHelloDir)

	cmd := NewInstallCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--dry-run"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	output := out.String()

	// Should only have 1 dry-run line (app-hello template)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 1 {
		t.Errorf("expected 1 dry-run line, got %d: %s", len(lines), output)
	}

	if !strings.Contains(output, "app-hello") {
		t.Error("expected app-hello in output")
	}
	if strings.Contains(output, "networking") {
		t.Error("should not contain networking in output")
	}
}
