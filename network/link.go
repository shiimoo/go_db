package network

import "github.com/shiimoo/godb/lib/base/service"

type Link interface {
	service.Service
	Key() any // 标识格式待定

	// 消息体字节长度([4]byte)
	// 包体总数([2]byte)
	// 当前包体序号([2]byte)
	// 包体节流([65535]byte)
	Read() ([]byte, error) // 读取数据(整包读取)
	Write([]byte) error    // 写入/发送数据(自动分包)
}

/* 消息包规则
消息体字节长度([4]byte)
包体总数([2]byte)
当前包体序号([2]byte)
包体节流([65535]byte)
*/
