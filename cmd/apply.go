package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"embed"
	"strings"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"gopkg.in/yaml.v3"
	"github.com/zclconf/go-cty/cty"
)

func init() {
	rootCmd.AddCommand(applyCommand)
}

//go:embed my-app-config.yaml ecs_container_definition_main.tf
var configFiles embed.FS 

// Define a struct that matches the structure of the YAML content
type AppConfig struct {
	AppName     string `yaml:"appName"`
	DockerImage string `yaml:"dockerImage"`
}

// Define a function to convert AppConfig to a Terraform variable file
func OutputTerraformVariables(config AppConfig) string {
	var sb strings.Builder

	// Header for Terraform variable file
	sb.WriteString("# Generated Terraform variable file\n\n")

	// Define variables for appName
	sb.WriteString(fmt.Sprintf("variable \"appName\" {\n"))
	sb.WriteString(fmt.Sprintf("  description = \"The name of the application\"\n"))
	sb.WriteString(fmt.Sprintf("  default     = \"%s\"\n", config.AppName))
	sb.WriteString("}\n\n")

	// Define variables for dockerImage
	sb.WriteString(fmt.Sprintf("variable \"dockerImage\" {\n"))
	sb.WriteString(fmt.Sprintf("  description = \"The Docker image for the application\"\n"))
	sb.WriteString(fmt.Sprintf("  default     = \"%s\"\n", config.DockerImage))
	sb.WriteString("}\n")

	return sb.String()
}

var applyCommand = &cobra.Command{
	Use:   "apply",
	Short: "Creates a new env.yml file with placeholder values.",


	// # cd my-app
	// # Lag en my-app-config.yaml
	// # Kjøre ok apply my-app-config.yaml
	// # Kjøre tf apply

	Run: func(cmd *cobra.Command, args []string) {
		// Read the file content
		var configContent, _ = configFiles.ReadFile("my-app-config.yaml")

		// Parse the YAML content into AppConfig struct
		var config AppConfig
		err := yaml.Unmarshal(configContent, &config)
		if err != nil {
			fmt.Printf("Error parsing YAML file: %s\n", err)
			return
		}

		// Now you can use the parsed data
		fmt.Printf("appName: %s\n", config.AppName)
		fmt.Printf("dockerImage: %s\n", config.DockerImage)

		

		// Write the parsed data back to a new file
		hclContent, _ := readHCLFile("ecs_container_definition_main.tf")

		// Modify the HCL content
		hclContent = modifyTfFile(hclContent)
		

		// Write the modified HCL content to a file
		err = writeHCLFile(hclContent, "new.tf")
		if err != nil {
			fmt.Printf("Error writing HCL file: %s\n", err)
			return
		}

		fmt.Println(hclContent)
	},
}





// modifyTfFile modifies a Terraform (HCL) file.
// tfFile: Pointer to the Terraform file to be modified.
// Returns the modified hclwrite.File.
// func modifyTfFile(tfFile *hclwrite.File) *hclwrite.File {

// 	localsBlock := tfFile.Body().FirstMatchingBlock("locals", nil)
// 	if localsBlock != nil {
// 		mainContainer := localsBlock.Body().FirstMatchingBlock("main_container", nil)
// 		body := mainContainer.Body()
// 		if body == nil {
// 			fmt.Println("mainContainer.Body() is nil")
// 		} else {
// 			body.SetAttributeValue("name", cty.StringVal("yoooooooooooooooooooo"))
// 		}
// 	} else {
// 		fmt.Println("No locals block found")
// 	}



// 	return tfFile
// }
// modifyHCL modifies the given HCL file to change the image attribute in the main_container block.
func modifyTfFile(file *hclwrite.File) *hclwrite.File {
	// Find the locals block
	for _, block := range file.Body().Blocks() {
		if block.Type() == "locals" {
			// Find the main_container block within the locals block
			for _, nestedBlock := range block.Body().Blocks() {
				if nestedBlock.Type() == "main_container" {
					// Set the image attribute to "hello"
					nestedBlock.Body().SetAttributeValue("image", cty.StringVal("hello"))
					return file // Return the modified file
				}
			}
		}
	}

	fmt.Println("No main_container block found")
	return file // Return the file unchanged if the specific block was not found
}



// readHCLFile reads and parses an HCL file.
// filename: Name of the HCL file to be read.
// Returns a pointer to the parsed hclwrite.File and an error if any.
// Errors are reported if file reading or parsing fails.
func readHCLFile(filename string) (*hclwrite.File, error) {
	content, err := configFiles.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading the embedded file: %w", err)
	}

	file, diag := hclwrite.ParseConfig(content, filename, hcl.InitialPos)
	if diag.HasErrors() {
		return nil, fmt.Errorf("error parsing HCL: %v", diag)
	}

	return file, nil
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