package main

import (
	"log"
	"net"
	"time"

	"github.com/shiimoo/godb/lib/base/util"
)

func main() {
	bs := []byte{
		1, 2, 3, 4, 5, 6, 7, 8,
		1, 2, 3, 4, 5, 6, 7, 8,
		1, 2, 3, 4, 5, 6, 7, 8,
		1, 2, 3, 4, 5, 6, 7, 8,
		1, 2, 3, 4, 5, 6, 7, 8,
	}
	linkObj, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	subPacks := util.SubPack(bs)
	max := uint(len(subPacks))
	for index, b := range subPacks {
		msg := make([]byte, 0)
		msg = append(msg, util.UintToBytes(max, 16)...)
		msg = append(msg, util.UintToBytes(uint(index+1), 16)...)
		msg = append(msg, util.UintToBytes(uint(len(b)), 16)...)
		msg = append(msg, b...)
		linkObj.Write(msg)
	}
	time.Sleep(1000 * time.Second)
}
