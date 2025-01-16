package utils

import "regexp"

// ValidateItemName checks if the item name contains only alphabetic characters
func ValidateItemName(name string) bool {
	re := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	return re.MatchString(name)
}
