package main

import (
	"context"

	"github.com/shiimoo/godb/lib/mlog"
	"github.com/shiimoo/godb/network/tcp"
)

func main() {
	closeChan := make(chan any, 1)

	rootCtx := context.Background()

	server, err := tcp.NewServer(rootCtx, "test", ":8080")
	if err != nil {
		mlog.Fatal("game", "start", err.Error())
	}
	server.Start()

	<-closeChan
}
