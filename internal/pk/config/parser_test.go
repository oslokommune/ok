package config

import (
	"os"
	"reflect"
	"testing"
)

func TestParseConfig(t *testing.T) {
	// Create a temporary YAML file for testing
	content := `
common:
  repo: "git@github.com:oslokommune/golden-path-boilerplate.git"
  path: "boilerplate/terraform/app"
  non_interactive: true
  ref: main
  var_files:
    - "../common-config.yml"
    - "config.yml"
  base_output_folder: "./stacks/dev"

templates:
  - name: "dev-networking"
    var_files:
      - "config.yml"
    ref: "networking-v2.8.3"
    subfolder: "networking"

  - name: "dev-app-km"
    non_interactive: false
    ref: "app-v2.2.3"
    subfolder: "app-km"

  - name: "dev-databases"
    ref: "databases-v9.12.0"
    subfolder: "databases"
`

	tmpfile, err := os.CreateTemp("", "config*.yml")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Error writing to temporary file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Error closing temporary file: %v", err)
	}

	// Test the ParseConfig function
	config, err := ParseConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("Error parsing config: %v", err)
	}

	// Check if the parsed config matches the expected values
	expectedConfig := &Config{
		Common: Common{
			Repo:             "git@github.com:oslokommune/golden-path-boilerplate.git",
			Path:             "boilerplate/terraform/app",
			NonInteractive:   true,
			Ref:              "main",
			VarFiles:         []string{"../common-config.yml", "config.yml"},
			BaseOutputFolder: "./stacks/dev",
		},
		Templates: []Template{
			{
				Name:      "dev-networking",
				VarFiles:  []string{"config.yml"},
				Ref:       "networking-v2.8.3",
				Subfolder: "networking",
			},
			{
				Name:           "dev-app-km",
				NonInteractive: boolPtr(false),
				Ref:            "app-v2.2.3",
				Subfolder:      "app-km",
			},
			{
				Name:      "dev-databases",
				Ref:       "databases-v9.12.0",
				Subfolder: "databases",
			},
		},
	}

	if !reflect.DeepEqual(config, expectedConfig) {
		t.Errorf("Parsed config does not match expected config.\nGot: %+v\nWant: %+v", config, expectedConfig)
	}
}

func boolPtr(b bool) *bool {
	return &b
}
