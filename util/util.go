package util

import (
	"fmt"
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

var replacer = strings.NewReplacer("r", "0x0A", "\n", "0x0B", "\t", "0x0C")

//ToStringOneLine : replace \t \r \n to character to debug on logger
func ToStringOneLine(s interface{}) string {
	str := fmt.Sprintf("%v", s)
	return replacer.Replace(str)
}
