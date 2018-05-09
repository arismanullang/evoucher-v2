package valid

import (
	"regexp"
)

var (
	rxSlug = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9-_]+[a-zA-Z0-9]$")
)

func IsSlug(s string) bool {
	return rxSlug.MatchString(s)
}
