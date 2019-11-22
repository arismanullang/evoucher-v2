package util

import (
	"strings"
)

//StandardizeSpaces : trim redundant spaces
//All leading/trailing whitespace or new-line characters, null characters, etc.
//Any redundant spaces within a string (ex. "hello[space][space]world" would be converted to "hello[space]world")
func StandardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

//SimplifyKeyString : trim redundant spaces and conver string word by word to simple lowercase string data
func SimplifyKeyString(val string) (key string) {
	key = strings.ToLower(StandardizeSpaces(val))
	return key
}
