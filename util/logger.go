package util

import (
	"fmt"
)

const (
	tagDebug   = "[DBG]"
	tagInfo    = "[INF]"
	tagWarning = "[WRN]"
	tagError   = "[ERR]"
)

//DEBUG debug
func DEBUG(i ...interface{}) {
	//fmt.Println(tagDebug, ToStringOneLine(i))
	fmt.Println(tagDebug, i)
}
