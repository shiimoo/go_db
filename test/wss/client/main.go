package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shiimoo/godb/lib/base/util"
	"github.com/shiimoo/godb/network"
)

func WebSocketServer() {
	addr := "localhost:8002"
	http.HandleFunc("/wshandler", WebSocketUpgrade)
	log.Println("Starting websocket server at " + addr)

	go func() {
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("WebSocket 服务器正在运行。按Ctrl+C退出")
	select {}
}

func WebSocketUpgrade(resp http.ResponseWriter, req *http.Request) {
	// 初始化 Upgrader
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	} // 使用默认的选项
	// 第三个参数是响应头,默认会初始化
	conn, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// 读取客户端的发送额消息,并返回
	go ReadMessage(conn)
	select {}
}

// 读取客户端发送的消息,并返回
func ReadMessage(conn *websocket.Conn) {
	for {
		// 消息类型:文本消息和二进制消息
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("receive msg:", string(msg))

		err = conn.WriteMessage(messageType, msg)
		if err != nil {
			log.Println("write error:", err)
			return
		}
	}
}

func main() {
	bs := []byte{
		1, 2, 3, 4, 5, 6, 7, 8,
		1, 2, 3, 4, 5, 6, 7, 8,
		1, 2, 3, 4, 5, 6, 7, 8,
		1, 2, 3, 4, 5, 6, 7, 8,
		1, 2, 3, 4, 5, 6, 7, 8,
	}
	linkObj, err := net.Dial(network.NetTypeTcp, "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	subPacks := util.SubPack(bs)
	max := uint(len(subPacks))
	for index, b := range subPacks {
		msg := make([]byte, 0)
		msg = append(msg, util.UintToBytes(max, 16)...)
		msg = append(msg, util.UintToBytes(uint(index+1), 16)...)
		msg = append(msg, util.UintToBytes(uint(len(b)), 16)...)
		msg = append(msg, b...)
		linkObj.Write(msg)
	}
	time.Sleep(1000 * time.Second)
}
