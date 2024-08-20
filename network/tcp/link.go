package tcp

import (
	"context"
	"fmt"
	"net"

	"github.com/shiimoo/godb/lib/base/service"
	"github.com/shiimoo/godb/lib/base/util"
)

type TcpLink struct {
	*service.BaseService

	id       uint         // 唯一id
	baseLink *net.TCPConn // 底层链接

	msgCount     uint64 // 接受消息数量
	msgPackBytes []byte // 消息字节缓存
}

func NewLink(parent context.Context, baseLink *net.TCPConn, id uint) *TcpLink {
	link := new(TcpLink)
	link.BaseService = service.NewService(parent, fmt.Sprintf("TcpLink_%d", id)) // todo 名称规范待定
	link.id = id
	link.baseLink = baseLink
	return link
}

// func (l *TcpLink) Start() error {
// 	panic(" Service Sub Class need to realize Service interface func Start() error")
// }

// func (l *TcpLink) Stop() error {
// 	panic(" Service Sub Class need to realize Service interface func Stop() error")
// }

// func (l *TcpLink) Close() error {
// 	panic(" Service Sub Class need to realize Service interface func Close() error")
// }

func (l *TcpLink) Key() uint {
	return l.id
}

func (l *TcpLink) Receive() {
	// go 协程获取
}

func (l *TcpLink) Read() ([]byte, error) {
	// 包体总数(uin16 [2]byte)
	packNumBuf := make([]byte, 2)
	_, err := l.baseLink.Read(packNumBuf)
	if err != nil {
		return nil, err
	}
	packNum := util.BytesToUint(packNumBuf)

	// 当前包体序号([2]byte)
	packIndexBuf := make([]byte, 2)
	_, err = l.baseLink.Read(packIndexBuf)
	if err != nil {
		return nil, err
	}
	packIndex := util.BytesToUint(packIndexBuf)

	// 包体字节总长度([2]byte)
	packSizeBuf := make([]byte, 2)
	_, err = l.baseLink.Read(packSizeBuf)
	if err != nil {
		return nil, err
	}
	packSize := util.BytesToUint(packSizeBuf)

	// 包体字节流([65535]byte)
	msgBuf := make([]byte, 65536)
	n, err := l.baseLink.Read(packSizeBuf)
	if err != nil {
		return nil, err
	}
	if uint(n) != packSize {
		return nil, nil // size不匹配
	}
	if packNum == packIndex {
		bs := l.msgPackBytes
		l.msgPackBytes = nil
		return bs, nil // 接受完毕
	}
	if l.msgPackBytes == nil {
		l.msgPackBytes = msgBuf
	} else {
		l.msgPackBytes = append(l.msgPackBytes, msgBuf...)
	}
	return l.Read()
}

// func (l *TcpLink) Write([]byte) error {
// 	return nil
// }
