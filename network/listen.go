package network

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/shiimoo/godb/lib/base/errors"
	"github.com/shiimoo/godb/lib/base/snowflake"
)

type Listen interface {
	NetType() string               // 网络类型
	Dispatch(bs []byte)            // 数据派发
	LinkCount() int                // 当前拥有的链接数量
	RemoveLink(id uint)            // 移除链接
	SendData(id uint, data []byte) // 发送数据
}

// 网络监听服务基类
type baseListenServer struct {
	ctx    context.Context    // 上下文
	cancel context.CancelFunc // 关闭方法

	id      uint         // 服务id
	_listen net.Listener // 监听器

	mu    sync.RWMutex  // links锁
	links map[uint]Link // 链接池
}

func newBaseListenServer(parent context.Context, address string) (*baseListenServer, error) {
	// @param address "0.0.0.0:8080"

	serverObj := new(baseListenServer)
	netType := serverObj.NetType()
	listenAddr, err := net.ResolveTCPAddr(netType, address)
	if err != nil {
		return nil, errors.NewErr(ErrListenAddrError, netType, address, err)
	}

	listener, err := net.ListenTCP(netType, listenAddr)
	if err != nil {
		return nil, errors.NewErr(ErrCreateListenError, netType, err)
	}
	// CREATE
	serverObj.ctx, serverObj.cancel = context.WithCancel(parent)
	serverObj._listen = listener
	serverObj.id = snowflake.GenUint()
	serverObj.links = make(map[uint]Link)
	return serverObj, nil
}

func (b *baseListenServer) NetType() string {
	return "" // todo 子类重写实现
}

func (b *baseListenServer) AddLink(linkObj Link) {
	b.mu.Lock()
	b.links[linkObj.ID()] = linkObj
	b.mu.Unlock()
}

func (b *baseListenServer) GetLink(id uint) Link {
	b.mu.Lock()
	defer b.mu.Unlock()
	linkObj, found := b.links[id]
	if found {
		return linkObj
	}
	return nil
}

func (b *baseListenServer) DelLink(linkObj Link) {
	b.mu.Lock()
	defer b.mu.Unlock()
	_, found := b.links[linkObj.ID()]
	if !found {
		return
	}
	delete(b.links, linkObj.ID())
}

func (b *baseListenServer) CloseLink(linkObj Link) {
	// todo 关闭原因
	linkObj.Close()
	b.DelLink(linkObj)
}

func (b *baseListenServer) CloseLinkByID(id uint) {
	linkObj := b.GetLink(id)
	if linkObj == nil {
		return
	}
	linkObj.Close()
	b.DelLink(linkObj)
}

func (b *baseListenServer) LinkCount() int {
	return len(b.links)
}

func (b *baseListenServer) SendData(id uint, data []byte) {
	if len(data) == 0 {
		return // 发送数据为空
	}
	linkObj := b.GetLink(id)
	if linkObj == nil {
		return // 链接不存在
	}
	if _, err := linkObj.Write(data); err != nil {
		b.CloseLink(linkObj)
	}
}

// Dispatch 数据派发: 链接获取到的数据进行派发
func (b *baseListenServer) Dispatch(id uint, bs []byte) {
	fmt.Println("todo 接受到的数据处理", id, len(bs), bs)
}
