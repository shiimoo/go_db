package main

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func main() {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	log.Printf("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Failed to connect to WebSocket server:", err)
	}
	defer c.Close()

	err = c.WriteMessage(websocket.TextMessage, []byte("Hello from client!"))
	if err != nil {
		log.Println("Failed to send message:", err)
		return
	}

	_, msg, err := c.ReadMessage()
	if err != nil {
		log.Println("Failed to receive message:", err)
		return
	}

	log.Printf("Received message from server: %s", msg)
}
