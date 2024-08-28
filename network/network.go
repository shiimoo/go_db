package network

import (
	"context"
	"fmt"
)

const (
	NetTypeTcp       = "tcp"
	NetTypeWebSocket = "webSocket"
	// todo udp
	NetTypeKcp = "kcp" // todo kcp
)

func NewListen(parent context.Context, netType string, address string, parmas ...any) (ListenServer, error) {
	switch netType {
	case NetTypeTcp:
		return NewTcpListenServer(parent, address, parmas...)
	case NetTypeWebSocket:
		return NewWebSocketListenServer(parent, address, parmas...)
	case NetTypeKcp:
		return NewKcpListenServer(parent, address, parmas...)
	default:
		panic(fmt.Sprintf("unknown net type :%s", netType))
	}
}
