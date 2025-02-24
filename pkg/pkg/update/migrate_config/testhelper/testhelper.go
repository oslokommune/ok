package testhelper

import "strings"

func TestNameToDir(name string) string {
	// replace spaces with dash
	return strings.ReplaceAll(name, " ", "-")
}
