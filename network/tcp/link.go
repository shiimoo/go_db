package tcp

import (
	"context"
	"fmt"
	"net"
	"time"

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
				bs, err := util.MergePack(l)
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

func (l *TcpLink) Read(b []byte) (int, error) {
	err := l.baseLink.SetDeadline(time.Now().Add(1 * time.Millisecond))
	if err != nil {
		return 0, err
	}
	return l.baseLink.Read(b)
}

func (l *TcpLink) clear() {
	l.msgPackBytes = nil
}

func (l *TcpLink) Write(data []byte) error {
	packs := util.SubPack(data)
	max := uint(len(packs))
	for index, b := range packs {
		msg := make([]byte, 0)
		msg = append(msg, util.UintToBytes(max, 16)...)
		msg = append(msg, util.UintToBytes(uint(index+1), 16)...)
		msg = append(msg, util.UintToBytes(uint(len(b)), 16)...)
		msg = append(msg, b...)
		if _, err := l.baseLink.Write(msg); err != nil {
			return err
		}
	}
	return nil
}
