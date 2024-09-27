package pkg

// argsContainsElement checks if an element is in an array which
// is used in the tab completion functions to filter out already selected elements
func argsContainsElement[T comparable](arr []T, element T) bool {
	for _, e := range arr {
		if e == element {
			return true
		}
	}
	return false
}
