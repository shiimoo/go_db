package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
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
	WebSocketServer()
}
