package network

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/shiimoo/godb/lib/base/errors"
	"github.com/shiimoo/godb/lib/base/snowflake"
	"github.com/xtaci/kcp-go"
)

type ListenServer interface {
	Ctx() context.Context    // 获取ctx
	GetListen() net.Listener // 获取底层监听器
	NetType() string         // 获取网络类型(tcp等)
	CloseCallBack()          // 关闭回调

	/* 管理 */
	AddLink(linkObj Link)                   // 添加链接
	GetLink(id uint) Link                   // 获取链接
	DelLink(linkObj Link, brokenType int)   // 删除链接
	CloseLink(linkObj Link, brokenType int) // 关闭链接:先关闭后删除
	CloseLinkByID(id uint, brokenType int)  // 关闭链接(id索引)：先关闭后删除
	Dispatch(id uint, bs []byte)            // 数据派发
	LinkCount() int                         // 当前拥有的链接数量
	SendData(id uint, data []byte)          // 发送数据

	/* 服务 */
	Start()
	Close()
}

// 网络监听服务基类
type baseListenServer struct {
	ctx    context.Context    // 上下文
	cancel context.CancelFunc // 关闭方法

	netType string       // 网络类型
	id      uint         // 服务id
	address string       // 监听地址
	_listen net.Listener // 监听器

	mu    sync.RWMutex  // links锁
	links map[uint]Link // 链接池
}

func newBaseListenServer(parent context.Context, netType, address string) (*baseListenServer, error) {
	// @param address "0.0.0.0:8080"

	serverObj := new(baseListenServer)

	var listener net.Listener
	var err error
	if netType == NetTypeTcp {
		listener, err = net.Listen(netType, address)
	} else if netType == NetTypeKcp {
		listener, err = kcp.ListenWithOptions(address, nil, 10, 3)
	}
	if err != nil {
		return nil, errors.NewErr(ErrCreateListenError, netType, address, err)
	}
	serverObj._listen = listener

	// CREATE
	serverObj.ctx, serverObj.cancel = context.WithCancel(parent)
	serverObj.netType = strings.TrimSpace(netType)
	serverObj.id = snowflake.GenUint()
	serverObj.address = address
	serverObj.links = make(map[uint]Link)
	return serverObj, nil
}

func (b *baseListenServer) Ctx() context.Context {
	return b.ctx
}
func (b *baseListenServer) GetListen() net.Listener {
	return b._listen
}

func (b *baseListenServer) NetType() string {
	return b.netType
}

/* 链接管理 */

func (b *baseListenServer) AddLink(linkObj Link) {
	b.mu.Lock()
	b.links[linkObj.ID()] = linkObj
	b.mu.Unlock()
	fmt.Println("添加链接", linkObj.ID(), b.LinkCount())
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

func (b *baseListenServer) DelLink(linkObj Link, brokenType int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.delLink(linkObj)
	fmt.Println(errors.NewErr(ErrLinkDisconnect, b.NetType(), brokenType))
}

func (b *baseListenServer) delLink(linkObj Link) {
	_, found := b.links[linkObj.ID()]
	if !found {
		return
	}
	delete(b.links, linkObj.ID())
	fmt.Println("链接关闭", linkObj.ID(), b.LinkCount())
}

func (b *baseListenServer) CloseLink(linkObj Link, brokenType int) {
	linkObj.Close(brokenType)
}

func (b *baseListenServer) CloseLinkByID(id uint, brokenType int) {
	linkObj := b.GetLink(id)
	if linkObj == nil {
		return
	}
	linkObj.Close(brokenType)
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
		b.CloseLink(linkObj, DisConnectTypeBroken)
	}
}

// Dispatch 数据派发: 链接获取到的数据进行派发
func (b *baseListenServer) Dispatch(id uint, bs []byte) {
	fmt.Println("todo 接受到的数据处理", id, len(bs), bs)
	b.SendData(id, bs)
}

func (b *baseListenServer) CloseCallBack() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, linkObj := range b.links {
		b.delLink(linkObj)
	}
}

/* todo service interface */

func (b *baseListenServer) Start() {
	panic("please rewrite in a subclass method[Start]")
}

func (b *baseListenServer) Close() {
	b.cancel()
}
