package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/oslokommune/ok/pkg/jsonschema"
)

func GenerateJsonSchemaForApp(ctx context.Context, downloader FileDownloader, stackPath, gitRef string) (*jsonschema.Document, error) {
	stacks, err := DownloadBoilerplateStacksWithDependencies(ctx, downloader, stackPath)
	if err != nil {
		return nil, fmt.Errorf("downloading boilerplate stacks: %w", err)
	}
	mobules := BuildModuleVariables(stacks)
	schema, err := TransformModulesToJsonSchema(fmt.Sprintf("%s-%s", stackPath, gitRef), mobules)
	if err != nil {
		return nil, fmt.Errorf("transforming modules to json schema: %w", err)
	}
	return schema, nil
}

// WriteJsonSchemaFile writes a json schema to a file in the output directory with the given template and version.
// The file will be named <template>-<version>.schema.json.
// The return value is the path to the file.
func WriteJsonSchemaFile(filePath string, schema *jsonschema.Document) (string, error) {
	outputDir := filepath.Dir(filePath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("creating output directory: %w", err)
	}

	slog.Debug("writing json schema file", slog.String("path", filePath), slog.String("schemaId", schema.ID))
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(schema); err != nil {
		return "", fmt.Errorf("encoding schema to file %s: %w", filePath, err)
	}

	return filePath, nil
}

const yamlLanguageServerComment = "# yaml-language-server:"

func CreateOrUpdateConfigurationFile(configFilePath string, schemaName string, schema *jsonschema.Document) (string, error) {
	configDir := filepath.Dir(configFilePath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("creating output directory: %w", err)
	}
	// write schema to the config dir with the name <schemaName>.schema.json
	schemaFileName := fmt.Sprintf("%s.schema.json", schemaName)
	schemaFilePath := filepath.Join(configDir, schemaFileName)
	slog.Info("writing schema file", slog.String("path", schemaFilePath), slog.String("configDir", configDir), slog.String("schemaId", schema.ID))
	_, err := WriteJsonSchemaFile(schemaFilePath, schema)
	if err != nil {
		return "", fmt.Errorf("writing schema file to %s: %w", schemaFilePath, err)
	}

	relativeSchemaPath, err := filepath.Rel(configDir, schemaFilePath)
	if err != nil {
		return "", fmt.Errorf("getting relative schema path: %w", err)
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("reading file: %w", err)
	}
	// Find if the first line starts with # yaml-language-server
	cleanedConfig := stripYamlLanguageServerComment(string(data))
	newConfig := appendYamlLanguageServerComment(cleanedConfig, relativeSchemaPath)
	err = os.WriteFile(configFilePath, []byte(newConfig), 0644)
	if err != nil {
		return "", fmt.Errorf("overwriting config file %s: %w", configFilePath, err)
	}

	return "", nil
}

func stripYamlLanguageServerComment(data string) string {
	first, rest, ok := strings.Cut(data, "\n")
	if ok && strings.HasPrefix(first, yamlLanguageServerComment) {
		return rest
	}
	return data
}

func appendYamlLanguageServerComment(data, schemaPath string) string {
	return fmt.Sprintf("%s $schema=%s\n%s", yamlLanguageServerComment, schemaPath, data)
}

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

func BuildModuleVariables(configs []*BoilerplateStack) []*ModuleVariables {
	if len(configs) == 0 {
		return nil
	}

	return buildModuleVariables("", configs[0], configs, "some/output/folder")
}

func buildModuleVariables(namespace string, currentConfig *BoilerplateStack, configs []*BoilerplateStack, outputFolder string) []*ModuleVariables {
	// ensure input arguments follow the correct format to avoid creating invalid namespaces
	namespace = JoinNamespaces(namespace)
	outputFolder = JoinPath(outputFolder, currentConfig.Path)

	namespaceVariables := make(map[string][]BoilerplateVariable)
	namespaceVariables[namespace] = currentConfig.Config.Variables

	for _, dep := range currentConfig.Config.Dependencies {
		depPath := JoinPath(currentConfig.Path, dep.TemplateUrl)
		depConfig, ok := findConfigFromPath(depPath, configs)
		if !ok {
			log.Printf("dependency %s not found in configs referenced by %s", depPath, currentConfig.Path)
			continue
		}

		depNamespace := namespace
		depOutputFolder := JoinPath(outputFolder, dep.OutputFolder)
		// if we move to a different output folder, then we need to create a new namespace
		if depOutputFolder != outputFolder {
			depNamespace = JoinNamespaces(depNamespace, dep.Name)
		}
		subModuleVariables := buildModuleVariables(depNamespace, depConfig, configs, depOutputFolder)
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

func JoinPath(base, path string) string {
	uri, err := url.JoinPath(base, path)
	if err != nil {
		slog.Error("could not join paths", slog.String("base", base), slog.String("path", path), slog.String("error", err.Error()))
		panic(err)
	}
	return uri
}

func TransformModulesToJsonSchema(schemaId string, modules []*ModuleVariables) (*jsonschema.Document, error) {
	var properties = make(map[string]jsonschema.Property)
	var requiredProperties []string

	for _, module := range modules {
		for _, variable := range module.Variables {
			if strings.Contains(variable.Description, "do NOT edit") {
				continue
			}
			typ := mapBoilerplateVariableTypeToSchemaType(variable.Type)
			name := JoinNamespaces(module.Namespace, variable.Name)
			if _, ok := properties[name]; ok {
				continue
			}
			// Special case for StackName, since it is a required property in all stacks.
			if variable.Name == "StackName" {
				requiredProperties = append(requiredProperties, name)
			}

			var requiredProperties []string
			var objectProperties map[string]jsonschema.Property
			// If we have an incoming map, we need to extract the keys and create a list of required properties
			// since all properties have to be overridden in an object to avoid unknown null values.
			// We also need to create a properties map for the default values in order to
			// give autocomplete suggestions in the editor.
			if typ == "object" {
				requiredProperties = extractKeysFromTypeMap(variable.Default)
				objectProperties = mapVariableObjectToProperties(variable)
				// We also need to add the properties as flat properties to the root object to allow for easy overriding
				// of single properties within a namespace.
				// Example: Override single property of "a.b.c" by setting "a.b.c: somevalue" to a new value.
				prefix := JoinNamespaces(module.Namespace, variable.Name)
				flattened := mapVariableObjectToFlatProperties(prefix, variable)
				addPropertiesIfNotExists(properties, flattened)

			}
			properties[name] = jsonschema.Property{
				Type:        typ,
				Description: variable.Description,
				Default:     variable.Default,
				Required:    requiredProperties,
				Properties:  objectProperties,
			}

		}
	}

	return &jsonschema.Document{
		ID:         schemaId,
		Schema:     jsonschema.SchemaURI,
		Type:       "object",
		Properties: properties,
		Required:   requiredProperties,
	}, nil
}

func addPropertiesIfNotExists(properties map[string]jsonschema.Property, newProperties map[string]jsonschema.Property) {
	for k, v := range newProperties {
		addPropertyIfNotExists(properties, k, v)
	}
}

func addPropertyIfNotExists(properties map[string]jsonschema.Property, name string, property jsonschema.Property) {
	if _, ok := properties[name]; ok {
		return
	}
	properties[name] = property
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

// mapVariableObjectToProperties turns a map into a list of properties with type information
func mapVariableObjectToProperties(variable BoilerplateVariable) map[string]jsonschema.Property {
	defaultMap, ok := variable.Default.(map[string]any)
	if !ok {
		return make(map[string]jsonschema.Property)
	}
	properties := make(map[string]jsonschema.Property)
	for k, v := range defaultMap {
		propertyName := k
		schemaType, ok := mapGoTypeToSchemaType(v)
		if !ok && v != nil {
			slog.Warn("could not transform default map type to schema type", slog.String("variable", propertyName), slog.Any("defaultValue", v))
		}
		properties[propertyName] = jsonschema.Property{
			Type:    schemaType,
			Default: v,
		}
	}
	return properties
}

// mapVariableObjectToFlatProperties turns a map into a flat list of properties prefixed with the namespace
// For example if the namespace is "a.b" and the map is {"c": 1, "d": 2} the result will be {"a.b.c": 1, "a.b.d": 2}
func mapVariableObjectToFlatProperties(namespace string, variable BoilerplateVariable) map[string]jsonschema.Property {
	defaultMap, ok := variable.Default.(map[string]any)
	if !ok {
		return make(map[string]jsonschema.Property)
	}
	properties := make(map[string]jsonschema.Property)
	for k, v := range defaultMap {
		propertyName := JoinNamespaces(namespace, k)
		schemaType, ok := mapGoTypeToSchemaType(v)
		if !ok {
			slog.Warn("could not transform default map type to schema type", slog.String("variable", propertyName))
		}
		properties[propertyName] = jsonschema.Property{
			Type:        schemaType,
			Description: fmt.Sprintf("Override single parameter of %s", namespace),
			Default:     v,
		}
	}
	return properties
}

func mapBoilerplateVariableTypeToSchemaType(t string) string {
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
		return "", false
	}
}
