package main

import (
	"fmt"

	"github.com/shiimoo/godb/lib/base/util"
)

func main() {
	var n uint = 39125678349
	fmt.Println(util.UintToBytes(n, 64))
	fmt.Println(util.BytesToUint(util.UintToBytes(n, 64)))
}
