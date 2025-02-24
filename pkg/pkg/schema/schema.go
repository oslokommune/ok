package schema

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/config"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/oslokommune/ok/pkg/jsonschema"
)

type Generator struct {
}

func NewGenerator() Generator {
	return Generator{}
}

// CreateJsonSchemaFile creates JSON schema file from a Boilerplate template configuration
func (g Generator) CreateJsonSchemaFile(
	ctx context.Context, manifestPackagePrefix string, pkg common.Package) ([]byte, error) {

	generatedSchema, err := GenerateJsonSchemaForApp(
		ctx, common.BoilerplateRepoOwner, common.BoilerplateRepoName, manifestPackagePrefix, pkg)
	if err != nil {
		return nil, fmt.Errorf("generating json schema for app: %w", err)
	}

	data, err := JsonSchemaToBytes(generatedSchema)
	if err != nil {
		return nil, fmt.Errorf("generating bytes from json: %w", err)
	}

	return data, nil
}

func GenerateJsonSchemaForApp(ctx context.Context, repoOwner string, repoName string, manifestPackagePrefix string, pkg common.Package) (*jsonschema.Document, error) {
	gh, err := githubreleases.GetGitHubClient()
	if err != nil {
		return nil, fmt.Errorf("getting GitHub client: %w", err)
	}

	downloader := githubreleases.NewFileDownloader(gh, repoOwner, repoName, pkg.Ref)
	stackPath := githubreleases.GetTemplatePath(manifestPackagePrefix, pkg.Template)

	generatedSchema, err := GenerateJsonSchemaForAppWithDownloader(ctx, downloader, stackPath, pkg.Ref)
	if err != nil {
		return nil, fmt.Errorf("generating json schema for app: %w", err)
	}

	return generatedSchema, nil
}

func GenerateJsonSchemaForAppWithDownloader(ctx context.Context, downloader config.FileDownloader, stackPath, gitRef string) (*jsonschema.Document, error) {
	stacks, err := config.DownloadBoilerplateStacksWithDependencies(ctx, downloader, stackPath)
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

func JsonSchemaToBytes(schema *jsonschema.Document) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")

	if err := enc.Encode(schema); err != nil {
		return nil, fmt.Errorf("encoding schema to bytes: %w", err)
	}

	return buf.Bytes(), nil
}

func WriteSchemaToFile(filePath string, data []byte) error {
	outputDir := filepath.Dir(filePath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	slog.Debug("writing json schema file", slog.String("path", filePath))

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("writing schema to file %s: %w", filePath, err)
	}

	return nil
}

const yamlLanguageServerComment = "# yaml-language-server:"

// GetSchemaFilePath converts a varFilePath to returns path to the schema file.
// Example
// Input:  config/app-hello.yaml
// Output: config/.schemas/app-v8.0.5.schema.json
func GetSchemaFilePath(varfilePath string, schemaName string) string {
	varFileDir := getVarFileDir(varfilePath)
	schemasDir := getSchemasDir(varFileDir)

	schemaFileName := fmt.Sprintf("%s.schema.json", schemaName)
	schemaFilePath := filepath.Join(schemasDir, schemaFileName)

	return schemaFilePath
}

func getVarFileDir(varfilePath string) string {
	return filepath.Dir(varfilePath)
}

func getSchemasDir(varFileDir string) string {
	return filepath.Join(varFileDir, ".schemas")
}

// SetVarFileSchemaDeclaration sets the first line of a var file to include a $schema reference to the schema file.
func SetVarFileSchemaDeclaration(varfilePath string, schemaName string) (string, error) {
	varFileDir := getVarFileDir(varfilePath)

	// TODO: Move this somewhhere else, it does not belong in this function.
	if err := os.MkdirAll(varFileDir, 0755); err != nil {
		return "", fmt.Errorf("creating output directory: %w", err)
	}

	schemaUri := fmt.Sprintf(
		"https://raw.githubusercontent.com/oslokommune/golden-path-boilerplate-schemas/refs/heads/main/schemas/%s.schema.json",
		schemaName)

	varFileData, err := os.ReadFile(varfilePath)
	if err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("reading file: %w", err)
	}

	varFileWithoutJsonSchemaDeclaration := stripYamlLanguageServerComment(string(varFileData))
	updatedVarFile := fmt.Sprintf("%s $schema=%s\n%s", yamlLanguageServerComment, schemaUri, varFileWithoutJsonSchemaDeclaration)

	err = os.WriteFile(varfilePath, []byte(updatedVarFile), 0644)
	if err != nil {
		return "", fmt.Errorf("overwriting config file %s: %w", varfilePath, err)
	}

	return "", nil
}

func stripYamlLanguageServerComment(varFileData string) string {
	first, rest, ok := strings.Cut(varFileData, "\n")
	if ok && strings.HasPrefix(first, yamlLanguageServerComment) {
		return rest
	}
	return varFileData
}

type Stack struct {
	Name         string
	Config       *config.BoilerplateConfig
	OutputFolder string
	Dependencies []string
}

type ModuleVariables struct {
	Namespace string
	Variables []config.BoilerplateVariable
}

type CombinedVariables struct {
	OutputFolder string
	Namespace    string
	Variables    []config.BoilerplateVariable
}

func BuildModuleVariables(configs []*config.BoilerplateStack) []*ModuleVariables {
	if len(configs) == 0 {
		return nil
	}

	return buildModuleVariables("", configs[0], configs, "some/output/folder")
}

func buildModuleVariables(namespace string, currentConfig *config.BoilerplateStack, configs []*config.BoilerplateStack, outputFolder string) []*ModuleVariables {
	// ensure input arguments follow the correct format to avoid creating invalid namespaces
	namespace = JoinNamespaces(namespace)
	outputFolder = config.JoinPath(outputFolder, currentConfig.Path)

	namespaceVariables := make(map[string][]config.BoilerplateVariable)
	namespaceVariables[namespace] = currentConfig.Config.Variables

	for _, dep := range currentConfig.Config.Dependencies {
		depPath := config.JoinPath(currentConfig.Path, dep.TemplateUrl)

		depConfig, ok := findConfigFromPath(depPath, configs)
		if !ok {
			log.Printf("dependency %s not found in configs referenced by %s", depPath, currentConfig.Path)
			continue
		}

		depNamespace := namespace
		depOutputFolder := config.JoinPath(outputFolder, dep.OutputFolder)

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

func findConfigFromPath(path string, configs []*config.BoilerplateStack) (*config.BoilerplateStack, bool) {
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
func mapVariableObjectToProperties(variable config.BoilerplateVariable) map[string]jsonschema.Property {
	defaultMap, ok := variable.Default.(map[string]any)
	if !ok {
		return make(map[string]jsonschema.Property)
	}
	properties := make(map[string]jsonschema.Property)
	for k, v := range defaultMap {
		propertyName := k

		schemaType, ok := mapGoTypeToSchemaType(v)
		if !ok && v != nil {
			slog.Debug("could not transform default map type to schema type", slog.String("variable", propertyName), slog.Any("defaultValue", v))
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
func mapVariableObjectToFlatProperties(namespace string, variable config.BoilerplateVariable) map[string]jsonschema.Property {
	defaultMap, ok := variable.Default.(map[string]any)
	if !ok {
		return make(map[string]jsonschema.Property)
	}

	properties := make(map[string]jsonschema.Property)

	for k, v := range defaultMap {
		propertyName := JoinNamespaces(namespace, k)

		schemaType, ok := mapGoTypeToSchemaType(v)
		if !ok {
			slog.Debug("could not transform default map type to schema type", slog.String("variable", propertyName))
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
