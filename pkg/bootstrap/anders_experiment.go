package bootstrap

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"gopkg.in/yaml.v3"
)

type Config struct {
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

func load() Config {
	data, err := os.ReadFile("env.yml")
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error unmarshaling YAML: %v", err)
	}

	return config
}

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

func Entry() {

	//config := load()

	tfFile, err := readHCLFile("example.tf")
	if err != nil {
		log.Fatalf("error reading HCL file: %v", err)
	}

	modifiedTfFile := modifyTfFile(tfFile)

	writeHCLFile(modifiedTfFile, "modified_example.tf")

}

func modifyTfFile(tfFile *hclwrite.File) *hclwrite.File {

	localsBlock := tfFile.Body().FirstMatchingBlock("locals", nil)
	if localsBlock != nil {
		localsBlock.Body().SetAttributeValue("hei", cty.StringVal("yoooooooooooooooooooo"))
	}

	return tfFile
}

func writeHCLFile(file *hclwrite.File, filename string) error {
	modifiedContent := file.Bytes()
	err := os.WriteFile(filename, modifiedContent, 0644)
	if err != nil {
		return fmt.Errorf("error writing modified file: %w", err)
	}
	return nil
}
