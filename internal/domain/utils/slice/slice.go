package slice

// IsStringInSlice returns true/false based of if the element "a" is present within "list"
func IsStringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
