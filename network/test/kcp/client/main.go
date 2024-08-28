package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shiimoo/godb/network"
)

func main() {

	bs := []byte{
		1, 2, 3, 4, 5, 6, 7, 8,
	}
	client, err := network.NewKcpClient(context.Background(), "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	client.Start()

	fmt.Println(client.Write(append([]byte{1, 0, 0}, bs...)))
	// fmt.Println(client.Write(append([]byte{2, 0, 0}, bs...)))
	time.Sleep(time.Second)
	fmt.Println(client.Write(append([]byte{3, 0, 0}, bs...)))

	time.Sleep(1000 * time.Second)
}
