package network

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/shiimoo/godb/lib/base/errors"
	"github.com/shiimoo/godb/lib/base/util"
)

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
	routing  string              // 路由地址
	upgrader *websocket.Upgrader // ws链接升级(http-->ws)
}

func NewWebSocketListenServer(parent context.Context, address string, parmas ...any) (*WebSocketListenServer, error) {
	serverObj := new(WebSocketListenServer)
	base, err := newBaseListenServer(parent, NetTypeWebSocket, address)
	if err != nil {
		return nil, err
	}
	if len(parmas) > 0 {
		routing, ok := parmas[0].(string)
		if !ok {
			return nil, errors.NewErr(
				ErrWsRouting,
				parmas[0],
				fmt.Sprintf("route Type must be string, but it's %s", reflect.TypeOf(parmas[0])),
			)
		}
		routing = strings.TrimSpace(routing)
		if routing != "" {
			serverObj.routing = "/" + routing
		}
	}

	// CREATE
	serverObj.baseListenServer = base
	serverObj.upgrader = &websocket.Upgrader{
		ReadBufferSize:  util.PackBytesLimit(),
		WriteBufferSize: util.PackBytesLimit(),
		CheckOrigin: func(r *http.Request) bool {
			// todo 请求合法性检查
			return true // 允许所有源，生产环境中应根据需要设置更严格的CORS策略
		},
	}
	return serverObj, nil
}

func (w *WebSocketListenServer) Start() {
	http.HandleFunc(w.routing, w.serveWs)
	log.Printf("Starting WebSocket server on %s...\n", w.address)
	err := http.ListenAndServe(w.address, nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}

func (w *WebSocketListenServer) @(resp http.ResponseWriter, req *http.Request) {
	conn, err := w.upgrader.Upgrade(resp, req, nil)
	if err != nil {
		log.Println("Failed to set up WebSocket connection:", err)
		return
	}
	// linkObj := NewLink(w.Ctx(), w.NetType(), fd, t)
	// t.AddLink(linkObj)
	// linkObj.Start()

	for {
		// 读取客户端发送的消息
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read WebSocket message:", err)
			break
		}

		log.Printf("Received message from client: %s", msg)

		// 假设我们只是简单地将接收到的消息回传给客户端
		err = conn.WriteMessage(msgType, msg)
		if err != nil {
			log.Println("Failed to send WebSocket message:", err)
			break
		}
	}
}
