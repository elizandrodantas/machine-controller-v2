package util

import "regexp"

func CleanSpecialCharacters(name string) string {
	reg := regexp.MustCompile(`[^a-zA-Z\\s]+`)
	cleanedName := reg.ReplaceAllString(name, "")
	return cleanedName
}
