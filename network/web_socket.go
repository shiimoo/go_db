package network

import "context"

/* WebSocket */

type WebSocketLink struct {
	*baseLink
}

func NewWebSocketLink(base *baseLink) *WebSocketLink {
	link := new(WebSocketLink)
	link.baseLink = base
	return link
}

// WebSocketListenServer webSocket服务
type WebSocketListenServer struct {
	*baseListenServer
}

func NewWebSocketListenServer(parent context.Context, address string) (*WebSocketListenServer, error) {
	serverObj := new(WebSocketListenServer)
	base, err := newBaseListenServer(parent, "tcp", address)
	if err != nil {
		return nil, err
	}
	// CREATE
	serverObj.baseListenServer = base
	return serverObj, nil
}

// ws 比较特殊start需要特殊处理
// 需要考虑下与其他网络链接类型的通用性
