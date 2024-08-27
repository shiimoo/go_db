package main

import (
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shiimoo/godb/lib/base/util"
)

func main() {
	bs := []byte{
		1, 2, 3, 4, 5, 6, 7, 8,
		1, 2, 3, 4, 5, 6, 7, 8,
		1, 2, 3, 4, 5, 6, 7, 8,
		1, 2, 3, 4, 5, 6, 7, 8,
		1, 2, 3, 4, 5, 6, 7, 8,
	}

	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8080", Path: "/ws"}
	log.Printf("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Failed to connect to WebSocket server:", err)
	}
	defer c.Close()
	subPacks := util.SubPack(bs)
	max := uint(len(subPacks))
	for index, b := range subPacks {
		msg := make([]byte, 0)
		msg = append(msg, util.UintToBytes(max, 16)...)
		msg = append(msg, util.UintToBytes(uint(index+1), 16)...)
		msg = append(msg, util.UintToBytes(uint(len(b)), 16)...)
		msg = append(msg, b...)
		err = c.WriteMessage(websocket.BinaryMessage, msg)
		if err != nil {
			log.Println("Failed to send message:", err)
			return
		}
	}

	time.Sleep(1000 * time.Second)
}
