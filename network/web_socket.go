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
	"github.com/shiimoo/godb/lib/base/snowflake"
	"github.com/shiimoo/godb/lib/base/util"
)

/* WebSocket */

type WebSocketLink struct {
	ctx    context.Context    // 上下文
	cancel context.CancelFunc // 关闭方法

	_fd           *websocket.Conn // websocket 升级conn
	_listenServer ListenServer    // 归属的监听服务(todo 专门建立管理服务，不依赖于监听服务?)
	byteCache     []byte          // 已接受的字节缓存

	id       uint   // 链接id
	msgCount uint64 // 接受消息数量

	brokenType int // 链接断开类型(关闭时写入)
}

func NewWebSocketLink(parent context.Context, fd *websocket.Conn, listenServer ListenServer) *WebSocketLink {
	obj := new(WebSocketLink)
	obj.ctx, obj.cancel = context.WithCancel(parent)
	obj._fd = fd
	obj._listenServer = listenServer
	obj.id = snowflake.GenUint()
	return obj
}

// ID 唯一标识性信息
func (wl *WebSocketLink) ID() uint {
	return wl.id
}

// ReadPack 读取数据包
func (wl *WebSocketLink) ReadPack() ([]byte, error) {
	msgType, bs, err := wl._fd.ReadMessage()
	if err != nil {
		return nil, err
	}
	if msgType != websocket.BinaryMessage {
		return nil, nil
	}

	// 包体总数(uin16 [2]byte)
	packNum := util.BytesToUint(bs[:2])
	// 当前包体序号([2]byte)
	packIndex := util.BytesToUint(bs[2:4])
	if packIndex > packNum {
		return nil, errors.NewErr(util.ErrPackNumError, packNum, packIndex)
	}

	// 包体字节总长度([2]byte)
	packSize := util.BytesToUint(bs[4:6])

	// 包体字节流(最大[65535]byte)
	msgBuf := bs[6:]
	if uint(len(msgBuf)) != packSize {
		return nil, errors.NewErr(util.ErrPackSizeError, packSize, len(msgBuf))
	}

	if packNum != packIndex {
		buf, err := wl.ReadPack()
		if err != nil {
			return nil, err
		}
		msgBuf = append(msgBuf, buf...)
	}
	return msgBuf, nil // 接受完毕
}

// Write : io.Writer realize
func (wl *WebSocketLink) Write(data []byte) (int, error) {
	packs := util.SubPack(data)
	max := uint(len(packs))
	for index, pack := range packs {
		msg := make([]byte, 0)
		msg = append(msg, util.UintToBytes(max, 16)...)
		msg = append(msg, util.UintToBytes(uint(index+1), 16)...)
		msg = append(msg, util.UintToBytes(uint(len(pack)), 16)...)
		msg = append(msg, pack...)
		if err := wl._fd.WriteMessage(websocket.BinaryMessage, msg); err != nil {
			return 0, err
		}
	}
	return len(data), nil
}

func (wl *WebSocketLink) Start() {
	go func() {
		for {
			select {
			case <-wl.ctx.Done():
				wl.CloseCallBack()
				return
			default:
				data, err := wl.ReadPack()
				if err != nil {
					wl.Close(DisConnectTypeBroken)
				} else {
					wl.msgCount += 1
					wl._listenServer.Dispatch(wl.id, data)
				}
			}
		}
	}()
}

// Close 关闭
func (wl *WebSocketLink) Close(brokenType int) {
	wl.brokenType = brokenType
	wl.cancel()
}

// CloseCallBack 关闭回调
func (wl *WebSocketLink) CloseCallBack() {
	wl._listenServer.DelLink(wl, wl.brokenType)
	wl._fd.Close()
}

func (wl *WebSocketLink) MsgCount() uint64 {
	return wl.msgCount
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
	go func() {
		http.HandleFunc(w.routing, w.serveWs)
		log.Printf("Starting WebSocket server on %s...\n", w.address)
		err := http.ListenAndServe(w.address, nil)
		if err != nil {
			log.Fatal("Error starting server:", err)
		}
	}()
}

func (w *WebSocketListenServer) serveWs(resp http.ResponseWriter, req *http.Request) {
	conn, err := w.upgrader.Upgrade(resp, req, nil)
	if err != nil {
		log.Println("Failed to set up WebSocket connection:", err)
		return
	}
	linkObj := NewWebSocketLink(w.Ctx(), conn, w)
	w.AddLink(linkObj)
	linkObj.Start()
}
