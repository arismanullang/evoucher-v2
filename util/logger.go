package util

import (
	"fmt"
)

//DEBUG debug
func DEBUG(i ...interface{}) {
	var tag interface{}
	tag = "[DEBUG]"
	fmt.Println(tag, i)
}
