package tcp

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/shiimoo/godb/lib/base/errors"
	"github.com/shiimoo/godb/lib/base/service"
	"github.com/shiimoo/godb/lib/base/util"
)

type TcpLink struct {
	*service.BaseService
	tcpServer *TcpServer // 归属的tcp服务

	id       uint         // 唯一id
	baseLink *net.TCPConn // 底层链接

	msgCount     uint64 // 接受消息数量
	msgPackBytes []byte // 消息字节缓存
}

func NewLink(parent context.Context, tcpServer *TcpServer, baseLink *net.TCPConn, id uint) *TcpLink {
	link := new(TcpLink)
	link.BaseService = service.NewService(parent, fmt.Sprintf("TcpLink_%d", id))
	link.tcpServer = tcpServer
	link.id = id
	link.baseLink = baseLink
	return link
}

func (l *TcpLink) Start() {
	go func() {
		for {
			select {
			case <-l.Context().Done():
				l.Close()
				return
			default:
				// 读取
				bs, err := l.Read()
				if err != nil {
					if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
						// 超时
						l.clear()
					} else {
						l.tcpServer.removeLinkObj(l, err)
					}
				} else {
					l.tcpServer.Dispatch(bs)
				}
			}
		}
	}()
}

func (l *TcpLink) Close() {
	l.baseLink.Close()
}

func (l *TcpLink) Key() uint {
	return l.id
}

func (l *TcpLink) Read() ([]byte, error) {
	// 包体总数(uin16 [2]byte)
	packNumBuf := make([]byte, 2)
	_, err := l.read(packNumBuf)
	if err != nil {
		return nil, err
	}
	packNum := util.BytesToUint(packNumBuf)
	// 当前包体序号([2]byte)
	packIndexBuf := make([]byte, 2)
	_, err = l.read(packIndexBuf)
	if err != nil {
		return nil, err
	}
	packIndex := util.BytesToUint(packIndexBuf)
	if packIndex > packNum {
		return nil, errors.NewErr(ErrTcpLinkPackNumError, packNum, packIndex)
	}

	// 包体字节总长度([2]byte)
	packSizeBuf := make([]byte, 2)
	_, err = l.read(packSizeBuf)
	if err != nil {
		return nil, err
	}
	packSize := util.BytesToUint(packSizeBuf)

	// 包体字节流(最大[65535]byte)
	msgBuf := make([]byte, packSize)
	n, err := l.read(msgBuf)
	if err != nil {
		return nil, err
	}
	if uint(n) != packSize {
		return nil, errors.NewErr(ErrTcpLinkPackSizeError, packSize, n)
	}
	if l.msgPackBytes == nil {
		l.msgPackBytes = msgBuf
	} else {
		l.msgPackBytes = append(l.msgPackBytes, msgBuf...)
	}

	if packNum == packIndex {
		bs := l.msgPackBytes
		l.clear()
		l.msgCount += 1
		return bs, nil // 接受完毕
	}
	return l.Read()
}

func (l *TcpLink) read(b []byte) (int, error) {
	err := l.baseLink.SetDeadline(time.Now().Add(1 * time.Millisecond))
	if err != nil {
		return 0, err
	}
	return l.baseLink.Read(b)
}

func (l *TcpLink) clear() {
	l.msgPackBytes = nil
}

// func (l *TcpLink) Write([]byte) error {
// 	return nil
// }
