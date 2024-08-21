package network

import "github.com/shiimoo/godb/lib/base/service"

type Listen interface {
	service.Service
	Dispatch(bs []byte)            // 数据派发
	LinkCount() int                // 当前拥有的链接数量
	RemoveLink(id uint)            // 移除链接
	SendData(id uint, data []byte) // 发送数据
}
