package common

import "strings"

// IsUrl checks if a string is a URL (http://, https://, or git@).
func IsUrl(str string) bool {
	return strings.HasPrefix(str, "http://") ||
		strings.HasPrefix(str, "https://") ||
		strings.HasPrefix(str, "git@")
}
