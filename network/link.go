package network

import (
	"io"
)

type Link interface {
	io.Writer                  // 数据写入
	ReadPack() ([]byte, error) // 读取数据包
	ID() uint                  // 唯一id标识
	Start()                    // 启动
	Close(brokenType int)      // 关闭 brokenType:关闭类型
	CloseCallBack()            // 关闭回调
}

/* 消息包规则
// 包体总数([2]byte)
// 当前包体序号([2]byte)
// 包体字节总长度([2]byte)
// 包体节流([最大65535]byte)
*/

const (
	DisConnectTypeBroken = 0 // 默认网络断开
)

type linkClient interface {
}
