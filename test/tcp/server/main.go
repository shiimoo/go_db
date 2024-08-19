package main

import (
	"net"

	"github.com/shiimoo/godb/lib/mlog"
)

func main() {
	listenAddr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:8080")
	if err != nil {
		mlog.Fatal("tcp", "addr", err.Error())
	}

	tcpListener, err := net.ListenTCP("tcp", listenAddr)
	if err != nil {
		mlog.Fatal("tcp", "listen", err.Error())
	}
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			mlog.Fatal("tcp", "acceptTCP", err.Error())
		}
		tcpConn.Read()
	}

}
