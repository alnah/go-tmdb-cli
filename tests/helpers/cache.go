package helpers

import (
	"fmt"
	"strings"
)

// CreateKey creates a cache key from parts.
func CreateKey(parts ...any) string {
	strParts := make([]string, len(parts))
	for i, part := range parts {
		strParts[i] = fmt.Sprint(part)
	}
	return strings.Join(strParts, "-")
}

// CreateValue creates a cache value from parts.
func CreateValue(parts ...any) string {
	strParts := make([]string, len(parts))
	for i, part := range parts {
		strParts[i] = fmt.Sprint(part)
	}
	return strings.Join(strParts, "-")
}
