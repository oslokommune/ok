package bootstrap

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"gopkg.in/yaml.v3"
)

// OriginalConfig represents the structure of the YAML configuration file (think of `env.yml`).
type OriginalConfig struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Team        string `yaml:"team"`
		Environment string `yaml:"environment"`
	} `yaml:"metadata"`
	AWS struct {
		AccountID string `yaml:"accountID"`
		Region    string `yaml:"region"`
	} `yaml:"aws"`
}

// JSONConfig represents the structure of the JSON configuration file (think of `_config.auto.tfvars.json`).
type JSONConfig struct {
	AccountID   string `json:"account_id"`
	Region      string `json:"region"`
	TeamName    string `json:"team_name"`
	Environment string `json:"environment"`
}

// ToJSONFile writes the JSONConfig to a file in JSON format.
// filename: Name of the file to write the JSON data to.
// Returns an error if file creation or JSON encoding fails.
func (c *JSONConfig) ToJSONFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	return encoder.Encode(c)
}

// TranslateToJSONConfig converts an OriginalConfig object to a JSONConfig object.
// orig: Pointer to the OriginalConfig to be converted.
// Returns a pointer to the new JSONConfig object.
func TranslateToJSONConfig(orig *OriginalConfig) *JSONConfig {
	return &JSONConfig{
		AccountID:   orig.AWS.AccountID,
		Region:      orig.AWS.Region,
		TeamName:    orig.Metadata.Team,
		Environment: orig.Metadata.Environment,
	}
}

// load reads and unmarshals a YAML file into an OriginalConfig object.
// Returns the unmarshalled OriginalConfig.
// Exits the program if reading or unmarshalling fails.
func load() OriginalConfig {
	data, err := os.ReadFile("env.yml")
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	var config OriginalConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error unmarshaling YAML: %v", err)
	}

	return config
}

// readHCLFile reads and parses an HCL file.
// filename: Name of the HCL file to be read.
// Returns a pointer to the parsed hclwrite.File and an error if any.
// Errors are reported if file reading or parsing fails.
func readHCLFile(filename string) (*hclwrite.File, error) {
	content, err := terraformFile.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading the embedded file: %w", err)
	}

	file, diag := hclwrite.ParseConfig(content, filename, hcl.InitialPos)
	if diag.HasErrors() {
		return nil, fmt.Errorf("error parsing HCL: %v", diag)
	}

	return file, nil
}

// Entry is the main entry point of the application logic.
// It orchestrates the loading, translating, and writing of configuration files.
func Entry() {

	origConfig := load()
	jsonConfig := TranslateToJSONConfig(&origConfig)

	err := jsonConfig.ToJSONFile("_config.auto.tfvars.json")
	if err != nil {
		fmt.Println("Error writing JSON file:", err)
	}

	tfFile, err := readHCLFile("example.tf")
	if err != nil {
		log.Fatalf("error reading HCL file: %v", err)
	}

	modifiedTfFile := modifyTfFile(tfFile)

	writeHCLFile(modifiedTfFile, "modified_example.tf")

}

// modifyTfFile modifies a Terraform (HCL) file.
// tfFile: Pointer to the Terraform file to be modified.
// Returns the modified hclwrite.File.
func modifyTfFile(tfFile *hclwrite.File) *hclwrite.File {

	localsBlock := tfFile.Body().FirstMatchingBlock("locals", nil)
	if localsBlock != nil {
		localsBlock.Body().SetAttributeValue("hei", cty.StringVal("yoooooooooooooooooooo"))
	}

	return tfFile
}

// writeHCLFile writes the modified HCL content to a file.
// file: Pointer to the hclwrite.File containing the HCL content.
// filename: Name of the file to write the modified HCL content to.
// Returns an error if writing the file fails.
func writeHCLFile(file *hclwrite.File, filename string) error {
	modifiedContent := file.Bytes()
	err := os.WriteFile(filename, modifiedContent, 0644)
	if err != nil {
		return fmt.Errorf("error writing modified file: %w", err)
	}
	return nil
}
