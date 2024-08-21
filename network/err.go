package network

import "github.com/shiimoo/godb/lib/base/errors"

var (
	// tcp监听地址错误
	ErrListenAddrError = errors.TempErr("ListenAddr netType[%s] addr[%s] err : %s")
	// tcp监听创建监听错误
	ErrCreateListenError = errors.TempErr("Create listen netType[%s] err : %s")
	// 链接读超时
	ErrTcpLinkReadTimeOutError = errors.TempErr("Tcp link read tyimeout")
)
