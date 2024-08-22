package network

import (
	"context"
)

type TcpLink struct {
	*baseLink
}

func NewTcpLink(base *baseLink) Link {
	link := new(TcpLink)
	link.baseLink = base
	return link
}

func (l *TcpLink) Key() uint {
	return l.ID()
}

// TcpListenServer tcp服务
type TcpListenServer struct {
	*baseListenServer
}

func NewTcpListenServer(parent context.Context, address string) (*TcpListenServer, error) {
	// @param address "0.0.0.0:8080"

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
