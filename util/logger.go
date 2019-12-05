package util

import (
	"fmt"
)

const tag = "[DEBUG]"

//DEBUG debug
func DEBUG(i ...interface{}) {
	var tag interface{}
	fmt.Println(tag, ToStringOneLine(i))
}
