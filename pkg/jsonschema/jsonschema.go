package jsonschema

const SchemaURI = "https://json-schema.org/draft-07/schema"

type (
	Document struct {
		ID         string              `json:"$id,omitempty"`
		Schema     string              `json:"$schema"`
		Title      string              `json:"title,omitempty"`
		Type       string              `json:"type,omitempty"`
		Properties map[string]Property `json:"properties,omitempty"`
		Required   []string            `json:"required,omitempty"`
	}

	Property struct {
		Type        string              `json:"type,omitempty"`
		Description string              `json:"description,omitempty"`
		Default     any                 `json:"default,omitempty"`
		Properties  map[string]Property `json:"properties,omitempty"`
		Required    []string            `json:"required,omitempty"`
	}
)

func MapGoTypeToSchemaType(v any) (string, bool) {
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
