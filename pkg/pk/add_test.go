package pk

import (
	"os"
	"path/filepath"
	"testing"
)

// mockGitHubReleases is a mock implementation of GitHubReleases.
type mockGitHubReleases struct {
	releases map[string]string
}

func (m *mockGitHubReleases) GetLatestReleases() (map[string]string, error) {
	return m.releases, nil
}

func TestAdd_AddsTemplate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pk-add-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	okDir := filepath.Join(tmpDir, ".ok")

	// Initialize first
	if err := Init(okDir); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	mock := &mockGitHubReleases{
		releases: map[string]string{
			"app":        "v10.2.2",
			"networking": "v3.0.1",
		},
	}

	opts := AddOptions{
		TemplateName: "app",
		Subfolder:    "my-app",
	}

	err = Add(okDir, opts, mock)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// Load and verify
	configs, err := LoadConfigs(okDir)
	if err != nil {
		t.Fatalf("failed to load configs: %v", err)
	}

	if len(configs) != 1 {
		t.Fatalf("expected 1 config, got %d", len(configs))
	}

	cfg := configs[0]
	if len(cfg.Templates) != 1 {
		t.Fatalf("expected 1 template, got %d", len(cfg.Templates))
	}

	tpl := cfg.Templates[0]
	if tpl.Name != "app" {
		t.Errorf("expected name %q, got %q", "app", tpl.Name)
	}
	if tpl.Path != "boilerplate/terraform/app" {
		t.Errorf("expected path %q, got %q", "boilerplate/terraform/app", tpl.Path)
	}
	if tpl.Ref != "app-v10.2.2" {
		t.Errorf("expected ref %q, got %q", "app-v10.2.2", tpl.Ref)
	}
	if tpl.Subfolder != "my-app" {
		t.Errorf("expected subfolder %q, got %q", "my-app", tpl.Subfolder)
	}
}

func TestAdd_UsesTemplateNameAsSubfolderDefault(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pk-add-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	okDir := filepath.Join(tmpDir, ".ok")

	if err := Init(okDir); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	mock := &mockGitHubReleases{
		releases: map[string]string{
			"networking": "v3.0.1",
		},
	}

	opts := AddOptions{
		TemplateName: "networking",
		// No subfolder specified
	}

	err = Add(okDir, opts, mock)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	configs, err := LoadConfigs(okDir)
	if err != nil {
		t.Fatalf("failed to load configs: %v", err)
	}

	tpl := configs[0].Templates[0]
	if tpl.Subfolder != "networking" {
		t.Errorf("expected subfolder %q, got %q", "networking", tpl.Subfolder)
	}
}

func TestAdd_WithExplicitRef(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pk-add-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	okDir := filepath.Join(tmpDir, ".ok")

	if err := Init(okDir); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	mock := &mockGitHubReleases{
		releases: map[string]string{
			"app": "v10.2.2",
		},
	}

	opts := AddOptions{
		TemplateName: "app",
		Subfolder:    "my-app",
		Ref:          "v9.0.0",
	}

	err = Add(okDir, opts, mock)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	configs, err := LoadConfigs(okDir)
	if err != nil {
		t.Fatalf("failed to load configs: %v", err)
	}

	tpl := configs[0].Templates[0]
	if tpl.Ref != "app-v9.0.0" {
		t.Errorf("expected ref %q, got %q", "app-v9.0.0", tpl.Ref)
	}
}

func TestAdd_FailsOnDuplicateSubfolder(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pk-add-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	okDir := filepath.Join(tmpDir, ".ok")

	if err := Init(okDir); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	mock := &mockGitHubReleases{
		releases: map[string]string{
			"app": "v10.2.2",
		},
	}

	// Add first template
	opts := AddOptions{
		TemplateName: "app",
		Subfolder:    "my-app",
	}

	if err := Add(okDir, opts, mock); err != nil {
		t.Fatalf("first Add failed: %v", err)
	}

	// Try to add with same subfolder
	err = Add(okDir, opts, mock)
	if err == nil {
		t.Fatal("expected Add to fail on duplicate subfolder")
	}
}

func TestAdd_MultipleTemplates(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pk-add-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	okDir := filepath.Join(tmpDir, ".ok")

	if err := Init(okDir); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	mock := &mockGitHubReleases{
		releases: map[string]string{
			"app":        "v10.2.2",
			"networking": "v3.0.1",
		},
	}

	// Add first template
	if err := Add(okDir, AddOptions{TemplateName: "app", Subfolder: "my-app"}, mock); err != nil {
		t.Fatalf("first Add failed: %v", err)
	}

	// Add second template
	if err := Add(okDir, AddOptions{TemplateName: "networking"}, mock); err != nil {
		t.Fatalf("second Add failed: %v", err)
	}

	configs, err := LoadConfigs(okDir)
	if err != nil {
		t.Fatalf("failed to load configs: %v", err)
	}

	if len(configs[0].Templates) != 2 {
		t.Fatalf("expected 2 templates, got %d", len(configs[0].Templates))
	}
}

func TestNormalizeRef(t *testing.T) {
	tests := []struct {
		template string
		ref      string
		want     string
	}{
		{"app", "v10.0.0", "app-v10.0.0"},
		{"app", "10.0.0", "app-v10.0.0"},
		{"app", "app-v10.0.0", "app-v10.0.0"},
		{"networking", "v3.0.1", "networking-v3.0.1"},
	}

	for _, tt := range tests {
		t.Run(tt.ref, func(t *testing.T) {
			got := normalizeRef(tt.template, tt.ref)
			if got != tt.want {
				t.Errorf("normalizeRef(%q, %q) = %q, want %q", tt.template, tt.ref, got, tt.want)
			}
		})
	}
}
