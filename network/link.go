package network

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/shiimoo/godb/lib/base/snowflake"
	"github.com/shiimoo/godb/lib/base/util"
)

type Link interface {
	io.Reader
	io.Writer
	ID() uint       // 唯一id标识
	Type() string   // 链接类型
	Start()         // 启动
	Close()         // 关闭
	CloseCallBack() // 关闭回调
}

/* 消息包规则
// 包体总数([2]byte)
// 当前包体序号([2]byte)
// 包体字节总长度([2]byte)
// 包体节流([最大65535]byte)
*/

// 接受的链接基类
type baseLink struct {
	ctx    context.Context    // 上下文
	cancel context.CancelFunc // 关闭方法

	_fd net.Conn // 套接字

	id       uint   // 链接id
	msgCount uint64 // 接受消息数量
}

func newBaseLink(parent context.Context, fd net.Conn) *baseLink {
	obj := new(baseLink)
	obj.ctx, obj.cancel = context.WithCancel(parent)
	obj._fd = fd
	obj.id = snowflake.GenUint()
	return obj
}

// ID 唯一标识性信息
func (b *baseLink) ID() uint {
	return b.id
}

// ID 唯一标识性信息
func (b *baseLink) Type() string {
	return "" // todo 子类重写实现
}

// Read : io.Reader realize
func (b *baseLink) Read(p []byte) (int, error) {
	err := b._fd.SetDeadline(time.Now().Add(1 * time.Millisecond))
	if err != nil {
		return 0, err
	}
	return b._fd.Read(p)
}

// Write : io.Writer realize
func (b *baseLink) Write(data []byte) (int, error) {
	packs := util.SubPack(data)
	max := uint(len(packs))
	count := 0
	for index, pack := range packs {
		msg := make([]byte, 0)
		msg = append(msg, util.UintToBytes(max, 16)...)
		msg = append(msg, util.UintToBytes(uint(index+1), 16)...)
		msg = append(msg, util.UintToBytes(uint(len(pack)), 16)...)
		msg = append(msg, pack...)
		if n, err := b._fd.Write(msg); err != nil {
			return count, err
		} else {
			count += n
		}
	}
	return len(data), nil
}

// Start 启动
func (b *baseLink) Start() {
	go func() {
		for {
			select {
			case <-b.ctx.Done():
				b.CloseCallBack()
				return
			default:
				_, err := util.MergePack(b)
				if err != nil {
					if netErr, ok := err.(net.Error); !ok || !netErr.Timeout() {
						b.Close()
					}
				} else {
					b.msgCount += 1
					// todo 获取的数据派发
					// l.tcpServer.Dispatch(bs)
				}
			}
		}
	}()
}

// Close 关闭
func (b *baseLink) Close() {
	b.cancel()
}

// CloseCallBack 关闭回调
func (b *baseLink) CloseCallBack() {
	// todo 通知管理器删除
}

/* exclusive method */

func (b *baseLink) MsgCount() uint64 {
	return b.msgCount
}
