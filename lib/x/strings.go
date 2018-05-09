package x

import (
	"strings"
)

func StringInSlice(needle string, haystack []string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}

func isNotCommaSpace(c rune) bool {
	return c == ',' || c == ' '
}

func SplitCommaWithTrim(s string) []string {
	return strings.FieldsFunc(s, isNotCommaSpace)
}
