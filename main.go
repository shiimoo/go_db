package main

import (
	"log"

	"github.com/shiimoo/godb/base"
)

func main() {
	// listTest()
	// err := savectrl.SaveBox(func() {
	// 	asveText()
	// })
	// log.Println("111", err)

	aa := []int{
		11, 22,
	}
	log.Println(aa[0], aa[1], aa[2:])
}

func asveText() {
	a := 0
	a = 10 / a
	// panic("savectrl.SaveBox test panic")
}

func listTest() {
	list := base.NewList()
	log.Println(list)
	list.Add(11)
	log.Println(list)
	list.Add(2)
	log.Println(list)

	log.Println("---------------")
	log.Println(list.Pop())
	log.Println(list)
}
