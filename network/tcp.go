package network

import (
	"context"
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

func NewTcpListenServer(parent context.Context, address string) (*TcpListenServer, error) {
	serverObj := new(TcpListenServer)
	base, err := newBaseListenServer(parent, "tcp", address)
	if err != nil {
		return nil, err
	}
	// CREATE
	serverObj.baseListenServer = base
	return serverObj, nil
}

func (t *TcpListenServer) Start() {
	startListen(t)
}
