package main

import (
	"context"

	"github.com/shiimoo/godb/lib/mlog"
	"github.com/shiimoo/godb/network"
)

func main() {
	// closeChan := make(chan any, 1)

	rootCtx := context.Background()

	server, err := network.NewTcpListenServer(rootCtx, ":8080")
	if err != nil {
		mlog.Fatal("game", "start", err.Error())
	}
	server.Start()

	for {
	}
}
