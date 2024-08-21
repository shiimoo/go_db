package util

import (
	"io"

	"github.com/shiimoo/godb/lib/base/errors"
)

const packBytesLimit = 65535

/* 消息包规则
// 包体总数([2]byte)
// 当前包体序号([2]byte)
// 包体字节总长度([2]byte)
// 包体节流([最大65535]byte)
*/

// SubPack 分包
func SubPack(bs []byte) [][]byte {
	subs := make([][]byte, 0)
	if bs == nil {
		return subs
	}
	length := len(bs)
	startIndex := 0
	endIndex := 0
	for startIndex < length-1 {
		endIndex = startIndex + packBytesLimit
		if endIndex >= length {
			endIndex = length
		}
		subs = append(subs, bs[startIndex:endIndex])
		startIndex = endIndex
	}
	return subs
}

// MergePack 合包
func MergePack(r io.Reader) ([]byte, error) {
	// 包体总数(uin16 [2]byte)
	packNumBuf := make([]byte, 2)
	_, err := r.Read(packNumBuf)
	if err != nil {
		return nil, err
	}
	packNum := BytesToUint(packNumBuf)
	// 当前包体序号([2]byte)
	packIndexBuf := make([]byte, 2)
	_, err = r.Read(packIndexBuf)
	if err != nil {
		return nil, err
	}
	packIndex := BytesToUint(packIndexBuf)
	if packIndex > packNum {
		return nil, errors.NewErr(ErrPackNumError, packNum, packIndex)
	}

	// 包体字节总长度([2]byte)
	packSizeBuf := make([]byte, 2)
	_, err = r.Read(packSizeBuf)
	if err != nil {
		return nil, err
	}
	packSize := BytesToUint(packSizeBuf)

	// 包体字节流(最大[65535]byte)
	msgBuf := make([]byte, packSize)
	n, err := r.Read(msgBuf)
	if err != nil {
		return nil, err
	}
	if uint(n) != packSize {
		return nil, errors.NewErr(ErrPackSizeError, packSize, n)
	}

	if packNum != packIndex {
		buf, err := MergePack(r)
		if err != nil {
			return nil, err
		}
		msgBuf = append(msgBuf, buf...)
	}
	return msgBuf, nil // 接受完毕
}
