package util

import (
	"fmt"

	"github.com/teris-io/shortid"
)

func init() {

}

func main() {
	sid := shortid.GetDefault()
	for index := 0; index < 100; index++ {
		fmt.Println(sid.Generate())
	}
}
