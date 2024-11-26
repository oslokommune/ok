package schema

import (
	"github.com/oslokommune/ok/pkg/jsonschema"
	"testing"

	"github.com/oslokommune/ok/pkg/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestFindConfigFromPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		configs  []*config.BoilerplateTemplate
		expected *config.BoilerplateTemplate
		found    bool
	}{
		{
			name: "config found",
			path: "path1",
			configs: []*config.BoilerplateTemplate{
				{Path: "path1"},
				{Path: "path2"},
			},
			expected: &config.BoilerplateTemplate{Path: "path1"},
			found:    true,
		},
		{
			name: "config not found",
			path: "path3",
			configs: []*config.BoilerplateTemplate{
				{Path: "path1"},
				{Path: "path2"},
			},
			expected: nil,
			found:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, found := findConfigFromPath(tt.path, tt.configs)
			require.Equal(t, tt.found, found)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestJoinNamespaces(t *testing.T) {
	tests := []struct {
		name       string
		namespaces []string
		expected   string
	}{
		{
			name:       "multiple namespaces",
			namespaces: []string{"namespace1", "namespace2", "namespace3"},
			expected:   "namespace1.namespace2.namespace3",
		},
		{
			name:       "single namespace",
			namespaces: []string{"namespace1"},
			expected:   "namespace1",
		},
		{
			name:       "empty namespaces",
			namespaces: []string{"namespace1", "", "namespace3"},
			expected:   "namespace1.namespace3",
		},
		{
			name:       "all empty namespaces",
			namespaces: []string{"", "", ""},
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := JoinNamespaces(tt.namespaces...)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestTransformModulesToJsonSchema(t *testing.T) {
	tests := []struct {
		name     string
		schemaId string
		modules  []*ModuleVariables
		expected *jsonschema.Document
	}{
		{
			name:     "single module with various types",
			schemaId: "test-schema",
			modules: []*ModuleVariables{
				{
					Namespace: "namespace1",
					Variables: []config.BoilerplateVariable{
						{Name: "var1", Type: "string", Description: "A string variable"},
						{Name: "var2", Type: "int", Description: "An integer variable"},
						{Name: "var3", Type: "bool", Description: "A boolean variable"},
					},
				},
			},
			expected: &jsonschema.Document{
				ID:     "test-schema",
				Schema: jsonschema.SchemaURI,
				Type:   "object",
				Properties: map[string]jsonschema.Property{
					"namespace1.var1": {Type: "string", Description: "A string variable"},
					"namespace1.var2": {Type: "integer", Description: "An integer variable"},
					"namespace1.var3": {Type: "boolean", Description: "A boolean variable"},
				},
				Required: []string(nil),
			},
		},
		{
			name:     "nested namespaces",
			schemaId: "test-schema",
			modules: []*ModuleVariables{
				{
					Namespace: "namespace1",
					Variables: []config.BoilerplateVariable{
						{Name: "var1", Type: "string", Description: "A string variable"},
					},
				},
				{
					Namespace: "namespace1.subnamespace",
					Variables: []config.BoilerplateVariable{
						{Name: "var2", Type: "int", Description: "An integer variable"},
					},
				},
			},
			expected: &jsonschema.Document{
				ID:     "test-schema",
				Schema: jsonschema.SchemaURI,
				Type:   "object",
				Properties: map[string]jsonschema.Property{
					"namespace1.var1":              {Type: "string", Description: "A string variable"},
					"namespace1.subnamespace.var2": {Type: "integer", Description: "An integer variable"},
				},
				Required: []string(nil),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := TransformModulesToJsonSchema(tt.schemaId, tt.modules)
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}
