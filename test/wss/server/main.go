package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有源，生产环境中应根据需要设置更严格的CORS策略
	},
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to set up WebSocket connection:", err)
		return
	}
	defer conn.Close()

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

func main() {
	http.HandleFunc("/ws", serveWs)
	log.Println("Starting WebSocket server on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
