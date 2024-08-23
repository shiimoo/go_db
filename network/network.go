package network

import (
	"context"
	"fmt"
	"net"
)

const (
	NetTypeTcp       = "tcp"
	NetTypeWebSocket = "webSocket"
)

func NewLink(parent context.Context, netType string, baseLink net.Conn, listenServer ListenServer) Link {
	base := newBaseLink(parent, baseLink, listenServer)
	switch netType {
	case NetTypeTcp:
		return NewTcpLink(base)
	case NetTypeWebSocket:
		return NewWebSocketLink(base)
	default:
		panic(fmt.Sprintf("unknown net type :%s", netType))
	}
}

func NewListen(parent context.Context, netType string, address string, parmas ...any) (ListenServer, error) {
	switch netType {
	case NetTypeTcp:
		return NewTcpListenServer(parent, address, parmas...)
	case NetTypeWebSocket:
		return NewWebSocketListenServer(parent, address, parmas...)
	default:
		panic(fmt.Sprintf("unknown net type :%s", netType))
	}
}
