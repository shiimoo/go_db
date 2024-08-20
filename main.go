package main

import (
	"fmt"

	"github.com/shiimoo/godb/lib/base/util"
)

func main() {
	// tyq(10000)
	// tyq(100000)
	// tyq(65535)
	bs := []byte{
		1, 2, 3, 4, 5, 6, 7, 8, 1,
		1, 2, 3, 4, 5, 6, 7, 8, 2,
		1, 2, 3, 4, 5, 6, 7, 8, 3,
		1, 2, 3, 4, 5, 6, 7, 8, 4,
		1, 2, 3, 4, 5, 6, 7, 8, 5,
	}

	fmt.Println(len(bs))
	for index, b := range util.SubPack(bs) {

		fmt.Println(index, b, len(b))
	}
}

func tyq(length int) {
	count := length / 65535
	other := length % 65535
	fmt.Println(count, other)
}
