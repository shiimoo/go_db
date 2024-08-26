package network

import (
	"context"
	"fmt"
)

const (
	NetTypeTcp       = "tcp"
	NetTypeWebSocket = "webSocket"
)

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
