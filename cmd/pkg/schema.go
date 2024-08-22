package pkg

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const JsonSchemaVersionString = "http://json-schema.org/draft-07/schema#"

var TestCommand = &cobra.Command{
	Use: "test",
	RunE: func(cmd *cobra.Command, args []string) error {
		var v any = make(map[string]any)
		bin, err := os.ReadFile("app-v8.0.2.dependencies.yml")
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(bin, &v); err != nil {
			return err
		}
		fmt.Printf("Type of root: %T\n", v)
		vasmap, ok := v.(map[string]any)
		if !ok {
			fmt.Printf("NOPE\n")
			return fmt.Errorf("not a map")
		}

		for k, vv := range vasmap {
			fmt.Printf("Type of %s: %T\n", k, vv)
			if _, ok := vv.(map[string]any); ok {
				fmt.Printf("YEAH IT MATCHES map any!\n")
			} else {
				fmt.Printf("NOPE\n")
			}
		}

		return nil
	},
}

var SchemaCommand = &cobra.Command{
	Use: "schema",
	RunE: func(cmd *cobra.Command, args []string) error {

		inputFile, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer inputFile.Close()
		dec := yaml.NewDecoder(inputFile)
		var dependencies = make(map[string]*DownloadedBoilThingy)
		if err = dec.Decode(&dependencies); err != nil {
			return err
		}

		rootCfg, err := findRootConfig(dependencies)
		if err != nil {
			return err
		}
		allVariables := collectFolderVariables("", rootCfg.Path, rootCfg, dependencies)
		jsonSchema := buildJsonSchemaFromNamespaceVariables(allVariables)
		schemaFileName := "my.schema.json"
		cmd.Printf("Writing schema file to %s\n", schemaFileName)
		if err := writeJsonSchemaToFile(jsonSchema, schemaFileName); err != nil {
			return err
		}

		return nil
	},
}

func writeJsonSchemaToFile(jsonSchema JsonSchemaRoot, filename string) error {
	bts, err := json.MarshalIndent(jsonSchema, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filename, bts, 0644); err != nil {
		return err
	}

	return nil
}

func collectFolderVariables(namespace, pathPrefix string, rootCfg *DownloadedBoilThingy, stacks map[string]*DownloadedBoilThingy) map[string][]BoilerplateVariable {
	folderVariables := make(map[string][]BoilerplateVariable)
	outputPath := mustJoinUri(pathPrefix, ".")

	folderVariables[namespace] = append(folderVariables[namespace], rootCfg.Variables...)

	templatePath := rootCfg.Path
	for _, dep := range rootCfg.Dependencies {
		depTemplatePath := mustJoinUri(templatePath, dep.TemplateUrl)
		depPath := mustJoinUri(outputPath, dep.OutputFolder)
		depNs := joinNamespacePath(namespace, dep.Name)
		if depPath == outputPath {
			depNs = namespace
		}
		depCfg, ok := stacks[depTemplatePath]
		if !ok {
			log.Printf("dependency not found: %s from template %s", depTemplatePath, templatePath)
			continue
		}
		for depPath, cfg := range collectFolderVariables(depNs, depPath, depCfg, stacks) {
			variablePath := depPath //mustJoinUri(outputPath, depPath)
			folderVariables[variablePath] = append(folderVariables[variablePath], cfg...)
		}
	}

	return folderVariables
}

func findRootConfig(dependencies map[string]*DownloadedBoilThingy) (*DownloadedBoilThingy, error) {
	for _, v := range dependencies {
		if v.IsRootCfg {
			return v, nil
		}
	}
	return nil, fmt.Errorf("no root config found")
}

type (
	JsonSchemaProperty struct {
		Type        string                        `json:"type,omitempty"`
		Description string                        `json:"description,omitempty"`
		Default     any                           `json:"default,omitempty"`
		Properties  map[string]JsonSchemaProperty `json:"properties,omitempty"`
		Required    []string                      `json:"required,omitempty"`
	}
	JsonSchemaRoot struct {
		Schema     string                        `json:"$schema"`
		Title      string                        `json:"title"`
		Type       string                        `json:"type"`
		Properties map[string]JsonSchemaProperty `json:"properties"`
		Required   []string                      `json:"required,omitempty"`
	}
)

func buildJsonSchemaFromNamespaceVariables(nsVariables map[string][]BoilerplateVariable) JsonSchemaRoot {
	properties := make(map[string]JsonSchemaProperty)
	for ns, variables := range nsVariables {
		for _, v := range variables {
			namespacedVariable := joinNamespacePath(ns, v.Name)
			if _, ok := properties[namespacedVariable]; ok {
				continue
			}

			variableType, ok := mapVariableTypeToSchemaType(v.Type)
			if !ok {
				goVariableType, _ := mapGoTypeToSchemaType(v.Default)
				variableType = goVariableType
			}
			var subProperties map[string]JsonSchemaProperty = nil
			if variableType == "object" {
				subProperties = mapDefaultMapToProperties("", v)
				for k, v := range subProperties {
					v.Description = fmt.Sprintf("Override single parameter of %s", namespacedVariable)
					properties[joinNamespacePath(namespacedVariable, k)] = v
				}
			}
			properties[namespacedVariable] = JsonSchemaProperty{
				Type:        variableType,
				Description: v.Description,
				Default:     v.Default,
				Properties:  subProperties,
				Required:    getStringKeys(subProperties),
			}

		}
	}

	// Make StackName required for each stack
	var requiredFields []string
	for _, ns := range getStringKeys(nsVariables) {
		namespacedStackNameVariable := joinNamespacePath(ns, "StackName")
		if _, ok := properties[namespacedStackNameVariable]; ok {
			requiredFields = append(requiredFields, namespacedStackNameVariable)
		}
	}
	return JsonSchemaRoot{
		Schema:     JsonSchemaVersionString,
		Title:      "Boilerplate Config",
		Type:       "object",
		Properties: properties,
		Required:   requiredFields,
	}
}

func mapVariableToSchemaProperty(variable BoilerplateVariable) JsonSchemaProperty {
	return JsonSchemaProperty{
		Type:        mapTypeToJsonSchema(variable.Type),
		Description: variable.Description,
		Default:     variable.Default,
	}
}
func getStringKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func mapDefaultMapToProperties(prefix string, variable BoilerplateVariable) map[string]JsonSchemaProperty {
	defaultMap, ok := variable.Default.(map[string]any)
	if !ok {
		return nil
	}
	var properties = make(map[string]JsonSchemaProperty)

	for k, v := range defaultMap {
		schemaType, ok := mapGoTypeToSchemaType(v)
		if !ok {
			continue
		}
		variableName := joinNamespacePath(prefix, k)
		properties[variableName] = JsonSchemaProperty{
			Type: schemaType,
			//Description: fmt.Sprintf("Part of: %s", variable.Description),
			Default: v,
		}
	}
	return properties
}

func mapGoTypeToSchemaType(v any) (string, bool) {
	switch v.(type) {
	case string:
		return "string", true
	case int, int32, int64, float32, float64:
		return "number", true
	case bool:
		return "boolean", true
	case map[string]any:
		return "object", true
	default:
		return fmt.Sprintf("unknown: %T", v), false
	}
}

func mapVariableTypeToSchemaType(variableType string) (string, bool) {
	switch variableType {
	case "string", "":
		return "string", true
	case "bool":
		return "boolean", true
	case "map":
		return "object", true
	default:
		return fmt.Sprintf("unknown: %s", variableType), false
	}
}

func init() {

	ConfigCommand.AddCommand(SchemaCommand)
	ConfigCommand.AddCommand(TestCommand)
}
