package tcp

import "github.com/shiimoo/godb/lib/base/errors"

var (
	// tcp监听地址错误
	ErrTcpListerAddrError = errors.TempErr("TcpListenAddr[%s] err : %s")
	// tcp监听创建监听错误
	ErrTcpListerError = errors.TempErr("TcpListen err : %s")
	// 链接读超时
	ErrTcpLinkReadTimeOutError = errors.TempErr("Tcp link read tyimeout")
	// 分包数据异常
	ErrTcpLinkPackNumError = errors.TempErr("Tcp link pack Err: total[%d] num[%d]")
	// 包体数据异常
	ErrTcpLinkPackSizeError = errors.TempErr("Tcp link pack size Err: total[%d] has[%d]")
)
