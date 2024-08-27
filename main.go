package main

import (
	"fmt"
)

func main() {
	// tyq(10000)
	// tyq(100000)
	// tyq(65535)
	// bs := []byte{
	// 	1, 2, 3, 4, 5, 6, 7, 8, 1,
	// 	1, 2, 3, 4, 5, 6, 7, 8, 2,
	// 	1, 2, 3, 4, 5, 6, 7, 8, 3,
	// 	1, 2, 3, 4, 5, 6, 7, 8, 4,
	// 	1, 2, 3, 4, 5, 6, 7, 8, 5,
	// }

	// fmt.Println(len(bs))
	// for index, b := range util.SubPack(bs) {

	// 	fmt.Println(index, b, len(b))
	// }

	original := []int{1, 2, 3, 4, 5, 6, 7}
	copied := make([]int, 5) //len(original))
	copy(copied, original)
	// if len(copied) >= len(original) {
	// 	fmt.Println(111)
	// 	copied = copied[:len(original)]
	// } else {
	// 	fmt.Println(222, original[len(copied):])
	// 	copied = append(copied, original[len(copied):]...)
	// }
	fmt.Println(copied) // 输出: [1 2 3 4 5]
}

func tyq(length int) {
	count := length / 65535
	other := length % 65535
	fmt.Println(count, other)
}
