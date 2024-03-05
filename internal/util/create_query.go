package util

import (
	"strings"
)

// CreateQueryStr creates a query string by extracting non-empty string fields from a struct.
// It takes an interface{} parameter 's' which should be a struct.
// It returns a string representing the query, where non-empty string fields are joined by '$'.
func CreateQueryStr(s []string) string {
	var result []string
	for _, value := range s {
		if value != "" {
			result = append(result, value)
		}
	}

	return strings.Join(result, "$")
}
