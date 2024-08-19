package tcp

import "github.com/shiimoo/godb/lib/base/errors"

var (
	// tcp监听地址错误
	ErrTcpListerAddrError = errors.TempErr("TcpListenAddr[%s] err : %s")
	// tcp监听创建监听错误
	ErrTcpListerError = errors.TempErr("TcpListen err : %s")
)
