package main

import (
	"context"

	"github.com/shiimoo/godb/lib/mlog"
	"github.com/shiimoo/godb/network"
)

func main() {

	rootCtx := context.Background()

	server, err := network.NewListen(rootCtx, network.NetTypeKcp, ":8080")
	if err != nil {
		mlog.Fatal("game", "start", err.Error())
	}
	server.Start()

	for {
	}
	// fmt.Println("kcp listens on 10000")
	// lis, err := kcp.ListenWithOptions(":10000", nil, 10, 3)
	// if err != nil {
	// 	panic(err)
	// }
	// for {
	// 	conn, e := lis.AcceptKCP()
	// 	if e != nil {
	// 		panic(e)
	// 	}
	// 	go func(conn net.Conn) {
	// 		var buffer = make([]byte, 1024, 1024)
	// 		for {
	// 			n, e := conn.Read(buffer)
	// 			if e != nil {
	// 				if e == io.EOF {
	// 					break
	// 				}
	// 				fmt.Println(e)
	// 				break
	// 			}

	// 			fmt.Println("receive from client:", string(buffer[:n]))
	// 		}
	// 	}(conn)
	// }
}
