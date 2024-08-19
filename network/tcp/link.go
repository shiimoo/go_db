package tcp

import (
	"context"
	"fmt"
	"net"

	"github.com/shiimoo/godb/lib/base/service"
)

type TcpLink struct {
	*service.BaseService
	ctx    context.Context    // 上下文
	cancel context.CancelFunc // 关闭方法

	id       uint
	baseLink *net.TCPConn
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

func (l *TcpLink) Read() ([]byte, error) {

	// 该读取会读取一个完整的消息体

	// 包体总数([2]byte)
	// 当前包体序号([2]byte)
	// 包体字节总长度([2]byte)
	// 包体字节流([65535]byte)

	// bs := make([]byte, 1024)
	// n, err := tcpConn.Read(bs)
	// if err != nil {
	// 	mlog.Warn("tcp", "acceptTCP", err.Error())
	// } else {
	// 	bs = bs[:n]
	// }
	// mlog.Infof("tcp", "", "%s-%s", []mlog.Data{
	// 	{Key: "data", Value: string(bs)}, {Key: "n", Value: n},
	// }...)
}

// func (l *TcpLink) Write([]byte) error {
// 	return nil
// }
