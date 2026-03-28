package pk

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInit_CreatesDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pk-init-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	okDir := filepath.Join(tmpDir, ".ok")

	err = Init(okDir, InitOptions{})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// Check directory exists
	info, err := os.Stat(okDir)
	if err != nil {
		t.Fatalf("failed to stat .ok dir: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected .ok to be a directory")
	}

	// Check config file exists
	configPath := filepath.Join(okDir, DefaultConfigFileName)
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("config file not created: %v", err)
	}

	// Load and verify config
	configs, err := LoadConfigs(okDir)
	if err != nil {
		t.Fatalf("failed to load configs: %v", err)
	}

	if len(configs) != 1 {
		t.Fatalf("expected 1 config, got %d", len(configs))
	}

	cfg := configs[0]
	if cfg.Common.Repo != DefaultRepo {
		t.Errorf("expected repo %q, got %q", DefaultRepo, cfg.Common.Repo)
	}
	if !cfg.Common.NonInteractive {
		t.Error("expected non_interactive to be true")
	}
	if cfg.Common.BaseOutputFolder != "." {
		t.Errorf("expected base_output_folder %q, got %q", ".", cfg.Common.BaseOutputFolder)
	}
	if len(cfg.Templates) != 0 {
		t.Errorf("expected 0 templates, got %d", len(cfg.Templates))
	}
}

func TestInit_FailsIfConfigExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pk-init-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	okDir := filepath.Join(tmpDir, ".ok")

	// Create the directory and config file first
	if err := os.Mkdir(okDir, 0755); err != nil {
		t.Fatalf("failed to create .ok dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(okDir, "config.yaml"), []byte(""), 0644); err != nil {
		t.Fatalf("failed to create config file: %v", err)
	}

	err = Init(okDir, InitOptions{})
	if err == nil {
		t.Fatal("expected Init to fail when config exists")
	}
}

func TestInit_WithCustomOptions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pk-init-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	okDir := filepath.Join(tmpDir, ".ok")

	err = Init(okDir, InitOptions{
		ConfigFileName:   "dev.yaml",
		BaseOutputFolder: "stacks/dev",
	})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// Check dev.yaml was created
	configPath := filepath.Join(okDir, "dev.yaml")
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("dev.yaml not created: %v", err)
	}

	// Load and verify
	configs, err := LoadConfigs(okDir)
	if err != nil {
		t.Fatalf("failed to load configs: %v", err)
	}

	cfg := configs[0]
	if cfg.Common.BaseOutputFolder != "stacks/dev" {
		t.Errorf("expected base_output_folder %q, got %q", "stacks/dev", cfg.Common.BaseOutputFolder)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Common.Repo != DefaultRepo {
		t.Errorf("expected repo %q, got %q", DefaultRepo, cfg.Common.Repo)
	}
	if !cfg.Common.NonInteractive {
		t.Error("expected non_interactive to be true")
	}
	if cfg.Common.BaseOutputFolder != "." {
		t.Errorf("expected base_output_folder %q, got %q", ".", cfg.Common.BaseOutputFolder)
	}
	if cfg.Templates == nil {
		t.Error("expected Templates to be initialized")
	}
}
