package network

import (
	"io"
	"net"
	"time"

	"github.com/shiimoo/godb/lib/base/service"
)

type Link interface {
	service.Service
	Key() any              // 标识格式待定
	Read() ([]byte, error) // 读取数据(整包读取)
	Write([]byte) error    // 写入/发送数据(自动分包)
}

/* 消息包规则
// 包体总数([2]byte)
// 当前包体序号([2]byte)
// 包体字节总长度([2]byte)
// 包体节流([最大65535]byte)
*/

// 接受的链接基类
type baseLink struct {
	_fd net.Conn // 套接字

	id                uint   // 链接id
	msgCount          uint64 // 接受消息数量
	receiveBytesCache []byte // 消息字节缓存(整体接受完毕后清空)

	t io.Reader
}

func (b *baseLink) Read(p []byte) (int, error) {
	err := b._fd.SetDeadline(time.Now().Add(1 * time.Millisecond))
	if err != nil {
		return 0, err
	}
	return b._fd.Read(p)
}
