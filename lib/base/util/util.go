package util

import (
	"encoding/binary"
	"math"
)

/* 字节<-->整型数字(无符号) */

// 无符号整形 --> 成字节
func UintToBytes(num, bit uint) []byte {
	var buf []byte
	if bit == 16 && num <= math.MaxUint16 {
		buf = make([]byte, 2)
		binary.BigEndian.PutUint16(buf, uint16(num))
	} else if bit == 32 && num <= math.MaxUint32 {
		buf = make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(num))
	} else if bit == 64 {
		buf = make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(num))
	}
	return buf
}

// 字节 --> 无符号整形
func BytesToUint(b []byte) uint {
	bsLen := len(b)
	var num uint = 0 // uint16 |uint32 | uint64
	if bsLen == 2 {
		num = uint(binary.BigEndian.Uint16(b))
	} else if bsLen == 4 {
		num = uint(binary.BigEndian.Uint32(b))
	} else if bsLen == 8 {
		num = uint(binary.BigEndian.Uint64(b))
	}
	return num
}
