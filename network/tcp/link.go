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

func (l *TcpLink) Key() uint {
	return l.id
}
