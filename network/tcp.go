package network

import (
	"context"
	"net"
	"time"

	"github.com/shiimoo/godb/lib/base/errors"
	"github.com/shiimoo/godb/lib/base/snowflake"
	"github.com/shiimoo/godb/lib/base/util"
	"github.com/shiimoo/godb/lib/mlog"
)

type TcpLink struct {
	ctx    context.Context    // 上下文
	cancel context.CancelFunc // 关闭方法
	// *baseLink
	_fd           net.Conn     // 套接字
	_listenServer ListenServer // 归属的监听服务(todo 专门建立管理服务，不依赖于监听服务?)

	id         uint   // 链接id
	msgCount   uint64 // 接受消息数量
	brokenType int    // 链接断开类型(关闭时写入)
}

func NewTcpLink(parent context.Context, netType string, fd net.Conn, listenServer ListenServer) *TcpLink {
	link := new(TcpLink)
	link.ctx, link.cancel = context.WithCancel(parent)
	link._fd = fd
	link._listenServer = listenServer
	link.id = snowflake.GenUint()
	return link
}

// ID 唯一标识性信息
func (tl *TcpLink) ID() uint {
	return tl.id
}

// Read : io.Reader realize
func (tl *TcpLink) Read(p []byte) (int, error) {
	err := tl._fd.SetDeadline(time.Now().Add(1 * time.Millisecond))
	if err != nil {
		return 0, err
	}
	return tl._fd.Read(p)
}

// ReadPack 读取数据包
func (tl *TcpLink) ReadPack() ([]byte, error) {
	// 包体总数(uin16 [2]byte)
	packNumBuf := make([]byte, 2)
	_, err := tl.Read(packNumBuf)
	if err != nil {
		return nil, err
	}
	packNum := util.BytesToUint(packNumBuf)
	// 当前包体序号([2]byte)
	packIndexBuf := make([]byte, 2)
	_, err = tl.Read(packIndexBuf)
	if err != nil {
		return nil, err
	}
	packIndex := util.BytesToUint(packIndexBuf)
	if packIndex > packNum {
		return nil, errors.NewErr(util.ErrPackNumError, packNum, packIndex)
	}

	// 包体字节总长度([2]byte)
	packSizeBuf := make([]byte, 2)
	_, err = tl.Read(packSizeBuf)
	if err != nil {
		return nil, err
	}
	packSize := util.BytesToUint(packSizeBuf)

	// 包体字节流(最大[65535]byte)
	msgBuf := make([]byte, packSize)
	n, err := tl.Read(msgBuf)
	if err != nil {
		return nil, err
	}
	if uint(n) != packSize {
		return nil, errors.NewErr(util.ErrPackSizeError, packSize, n)
	}

	if packNum != packIndex {
		buf, err := tl.ReadPack()
		if err != nil {
			return nil, err
		}
		msgBuf = append(msgBuf, buf...)
	}
	return msgBuf, nil // 接受完毕
}

// Write : io.Writer realize
func (tl *TcpLink) Write(data []byte) (int, error) {
	packs := util.SubPack(data)
	max := uint(len(packs))
	count := 0
	for index, pack := range packs {
		msg := make([]byte, 0)
		msg = append(msg, util.UintToBytes(max, 16)...)
		msg = append(msg, util.UintToBytes(uint(index+1), 16)...)
		msg = append(msg, util.UintToBytes(uint(len(pack)), 16)...)
		msg = append(msg, pack...)
		if n, err := tl._fd.Write(msg); err != nil {
			return count, err
		} else {
			count += n
		}
	}
	return len(data), nil
}

// Start 启动
func (tl *TcpLink) Start() {
	go func() {
		for {
			select {
			case <-tl.ctx.Done():
				tl.CloseCallBack()
				return
			default:
				data, err := tl.ReadPack()
				if err != nil {
					if netErr, ok := err.(net.Error); !ok || !netErr.Timeout() {
						tl.Close(DisConnectTypeBroken)
					}
				} else {
					tl.msgCount += 1
					tl._listenServer.Dispatch(tl.id, data)
				}
			}
		}
	}()
}

// Close 关闭
func (tl *TcpLink) Close(brokenType int) {
	tl.brokenType = brokenType
	tl.cancel()

}

// CloseCallBack 关闭回调
func (tl *TcpLink) CloseCallBack() {
	tl._listenServer.DelLink(tl, tl.brokenType)
	tl._fd.Close()
}

/* exclusive method */

func (tl *TcpLink) MsgCount() uint64 {
	return tl.msgCount
}

// TcpListenServer tcp服务
type TcpListenServer struct {
	*baseListenServer
}

func NewTcpListenServer(parent context.Context, address string, _ ...any) (*TcpListenServer, error) {
	serverObj := new(TcpListenServer)
	base, err := newBaseListenServer(parent, NetTypeTcp, address)
	if err != nil {
		return nil, err
	}
	// CREATE
	serverObj.baseListenServer = base
	return serverObj, nil
}

func (t *TcpListenServer) Start() {
	go func() {
		for {
			select {
			case <-t.Ctx().Done():
				t.CloseCallBack()
				return
			default:
				// 监听链接
				fd, err := t.GetListen().Accept()
				if err != nil {
					mlog.Warn(NetTypeTcp, "acceptTCP", err.Error())
				} else {
					linkObj := NewTcpLink(t.Ctx(), t.NetType(), fd, t)
					t.AddLink(linkObj)
					linkObj.Start()
				}
			}
		}
	}()
}
