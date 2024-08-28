package main

import (
	"context"
	"log"
	"time"

	"github.com/shiimoo/godb/network"
)

func main() {
	// kcpconn, err := kcp.DialWithOptions("localhost:10000", nil, 10, 3)
	// if err != nil {
	// 	panic(err)
	// }

	// kcpconn.Write([]byte("hello kcp.emmmmmmmmmmmmmmm"))
	// select {}

	bs := []byte{
		1, 2, 3, 4, 5, 6, 7, 8,
		1, 2, 3, 4, 5, 6, 7, 8,
		1, 2, 3, 4, 5, 6, 7, 8,
		1, 2, 3, 4, 5, 6, 7, 8,
		1, 2, 3, 4, 5, 6, 7, 8,
	}
	client, err := network.NewKcpClient(context.Background(), "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	client.Start()

	client.Write(bs)

	time.Sleep(1000 * time.Second)
}
