package config

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/oslokommune/ok/pkg/jsonschema"
)

func BuildJsonSchemaFromConfig(config *BoilerplateConfig, dependencies []BoilerplateConfig) (*jsonschema.Document, error) {
	return nil, fmt.Errorf("not implemented")
}

type Stack struct {
	Name         string
	Config       *BoilerplateConfig
	OutputFolder string
	Dependencies []string
}

type ModuleVariables struct {
	Namespace string
	Variables []BoilerplateVariable
}

type CombinedVariables struct {
	OutputFolder string
	Namespace    string
	Variables    []BoilerplateVariable
}

func BuildModuleVariables(namespace string, currentConfig *BoilerplateStack, configs []*BoilerplateStack, outputFolder string) []*ModuleVariables {
	// ensure input arguments follow the correct format to avoid creating invalid namespaces
	namespace = JoinNamespaces(namespace)
	outputFolder = mustJoinUri(outputFolder, currentConfig.Path)

	namespaceVariables := make(map[string][]BoilerplateVariable)
	namespaceVariables[namespace] = currentConfig.Config.Variables

	for _, dep := range currentConfig.Config.Dependencies {
		depPath := mustJoinUri(currentConfig.Path, dep.TemplateUrl)
		depConfig, ok := findConfigFromPath(depPath, configs)
		if !ok {
			log.Printf("dependency %s not found in configs referenced by %s", depPath, currentConfig.Path)
			continue
		}

		depNamespace := namespace
		depOutputFolder := mustJoinUri(outputFolder, dep.OutputFolder)
		// if we move to a different output folder, then we need to create a new namespace
		if depOutputFolder != outputFolder {
			depNamespace = JoinNamespaces(depNamespace, dep.Name)
		}
		subModuleVariables := BuildModuleVariables(depNamespace, depConfig, configs, depOutputFolder)
		for _, m := range subModuleVariables {
			namespaceVariables[m.Namespace] = append(namespaceVariables[m.Namespace], m.Variables...)
		}
	}

	moduleVariables := make([]*ModuleVariables, 0, len(namespaceVariables))
	for namespace, variables := range namespaceVariables {
		moduleVariables = append(moduleVariables, &ModuleVariables{
			Namespace: namespace,
			Variables: variables,
		})
	}

	return moduleVariables
}

func findConfigFromPath(path string, configs []*BoilerplateStack) (*BoilerplateStack, bool) {
	for _, c := range configs {
		if c.Path == path {
			return c, true
		}
	}
	return nil, false
}

func JoinNamespaces(namespaces ...string) string {
	filtered := make([]string, 0)
	for _, n := range namespaces {
		if n != "" {
			filtered = append(filtered, n)
		}
	}

	return strings.Join(filtered, ".")
}

/*
func transformConfigsToStacks(rootStack *BoilerplateStack, allStacks []*BoilerplateStack) ([]*Stack, error) {
	currentPath := rootStack.Path
	folderStacks := make(map[string][]string)
	fmt.Println()

	for _, dep := range rootStack.Dependencies {
		depPath := mustJoinUri(currentPath, dep.OutputFolder)
		folderStacks[depPath] = append(folderStacks[depPath], dep.TemplateUrl)
		dependencyConfig, ok := packageConfigs[dep.TemplateUrl]
		if !ok {
			continue
		}

		stack, err := transformConfigsToStacks(depPath, &dependencyConfig, packageConfigs)
		for _, s := range stack {
			folderStacks[depPath] = append(folderStacks[depPath], s...)
		}
		if err != nil {
			return nil, err
		}

	}
	return nil, fmt.Errorf("not implemented")
}
*/

func mustJoinUri(base, path string) string {
	uri, err := url.JoinPath(base, path)
	if err != nil {
		panic(err)
	}
	return uri
}

func TransformModulesToJsonSchema(modules []*ModuleVariables) (*jsonschema.Document, error) {
	var properties = make(map[string]jsonschema.Property)
	var requiredProperties []string

	for _, module := range modules {
		for _, variable := range module.Variables {
			if strings.Contains(variable.Description, "do NOT edit") {
				continue
			}
			typ := mapVariableTypeToJsonSchema(variable.Type)
			name := JoinNamespaces(module.Namespace, variable.Name)
			if _, ok := properties[name]; ok {
				continue
			}
			if variable.Name == "StackName" {
				requiredProperties = append(requiredProperties, name)
			}

			var requiredProperties []string
			if typ == "object" {
				requiredProperties = extractKeysFromTypeMap(variable.Default)
			}
			properties[name] = jsonschema.Property{
				Type:        typ,
				Description: variable.Description,
				Default:     variable.Default,
				Required:    requiredProperties,
			}

			// If we have a map, we need to flatten the properties
			if typ == "object" {
				prefix := JoinNamespaces(module.Namespace, variable.Name)
				flattened := mapVariableObjectToFlatProperties(prefix, variable)
				for k, v := range flattened {
					if _, ok := properties[k]; ok {
						continue
					}
					properties[k] = v
				}
				for name, p := range mapVariableObjectToFlatProperties(prefix, variable) {
					properties[name] = p
				}

			}

			/*
				if variable.Type == "map" {
					properties := mapVariableObjectToFlatProperties(module.Namespace, variable)
					for _, p := range properties {
						p.Type = mapVariableTypeToJsonSchema(p.Type)
					}
				}
			*/
		}
	}

	return &jsonschema.Document{
		Schema:     jsonschema.SchemaURI,
		Type:       "object",
		Properties: properties,
		Required:   requiredProperties,
	}, nil
}

func extractKeysFromTypeMap(d any) []string {
	m, ok := d.(map[string]any)
	if !ok {
		return nil
	}
	return getMapKeyNames(m)
}

func getMapKeyNames[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

/*
	func mapVariableToSchemaProperty(namespace string, variable BoilerplateVariable) map[string]jsonschema.Property {
		typ := mapVariableTypeToJsonSchema(variable.Type)
	}
*/

func mapVariableObjectToFlatProperties(namespace string, variable BoilerplateVariable) map[string]jsonschema.Property {
	defaultMap, ok := variable.Default.(map[string]any)
	if !ok {
		return nil
	}
	properties := make(map[string]jsonschema.Property)
	for k, v := range defaultMap {
		propertyName := JoinNamespaces(namespace, k)
		schemaType, ok := mapGoTypeToSchemaType(v)
		if !ok {
			log.Printf("could not map type %T to schema type for variable %s", v, propertyName)
			schemaType = ""
		}
		properties[propertyName] = jsonschema.Property{
			Type:        schemaType,
			Description: fmt.Sprintf("Override single parameter of %s", namespace),
			Default:     v,
		}
	}
	return properties
}
func mapVariableTypeToJsonSchema(t string) string {
	switch t {
	case "map":
		return "object"
	case "int":
		return "integer"
	case "bool":
		return "boolean"
	case "string":
		return "string"
	default:
		return "string"
	}
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
