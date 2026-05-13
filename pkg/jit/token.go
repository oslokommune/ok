package jit

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

// ExtractGroupsFromJWT parses a JWT and extracts the "groups" claim.
func ExtractGroupsFromJWT(token string) []string {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil
	}

	raw, ok := claims["groups"]
	if !ok {
		return nil
	}

	arr, ok := raw.([]interface{})
	if !ok {
		return nil
	}

	groups := make([]string, 0, len(arr))
	for _, v := range arr {
		if s, ok := v.(string); ok {
			groups = append(groups, s)
		}
	}
	return groups
}
