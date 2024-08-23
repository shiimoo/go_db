package network

import (
	"context"

	"github.com/shiimoo/godb/lib/mlog"
)

type TcpLink struct {
	*baseLink
}

func NewTcpLink(base *baseLink) *TcpLink {
	link := new(TcpLink)
	link.baseLink = base
	return link
}

// TcpListenServer tcp服务
type TcpListenServer struct {
	*baseListenServer
}

func NewTcpListenServer(parent context.Context, address string, _ ...any) (*TcpListenServer, error) {
	serverObj := new(TcpListenServer)
	base, err := newBaseListenServer(parent, NetTypeTcp, address)
	if err != nil {
		return nil, err
	}
	// CREATE
	serverObj.baseListenServer = base
	return serverObj, nil
}

func (t *TcpListenServer) Start() {
	go func() {
		for {
			select {
			case <-t.Ctx().Done():
				t.CloseCallBack()
				return
			default:
				// 监听链接
				fd, err := t.GetListen().Accept()
				if err != nil {
					mlog.Warn(NetTypeTcp, "acceptTCP", err.Error())
				} else {
					linkObj := NewLink(t.Ctx(), t.NetType(), fd, t)
					t.AddLink(linkObj)
					linkObj.Start()
				}
			}
		}
	}()
}
