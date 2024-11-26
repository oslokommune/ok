package schema

import (
	"fmt"
	"github.com/oslokommune/ok/pkg/jsonschema"
	"github.com/oslokommune/ok/pkg/pkg/config"
	"log/slog"
	"strings"
)

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
