package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	// ALPHABET string
	ALPHABET = "alphabet"
	// NUMERALS string
	NUMERALS = "numeric"
	// ALPHANUMERIC string
	ALPHANUMERIC = "alphanumeric"

	AlphabetString     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	NumeralsString     = "1234567890"
	AlphaNumericString = AlphabetString + NumeralsString

	// LENGTH = default random length number
	LENGTH = 8
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

var replacer = strings.NewReplacer("\r", "0x0A", "\n", "0x0B", "\t", "0x0C")

//ToStringOneLine : replace \t \r \n to character to debug on logger
func ToStringOneLine(s interface{}) string {
	str := fmt.Sprintf("%v", s)
	return replacer.Replace(str)
}

//RandomizeString : randomize string with custom length and random type
func RandomizeString(ln int, fm string) string {
	CharsType := map[string]string{
		ALPHABET:     AlphabetString,
		NUMERALS:     NumeralsString,
		ALPHANUMERIC: AlphaNumericString,
	}

	rand.Seed(time.Now().UTC().UnixNano())
	chars := CharsType[fm]
	result := make([]byte, ln)
	for i := 0; i < ln; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}

	return string(result)
}
