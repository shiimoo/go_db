package mlog

import (
	"context"
	"log"
	"sync"
)

// ***** 生成管理器 *****

type mgr struct {
	key string // 唯一标识(id)
	sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc

	set map[string]*logger
}

// ***** 基础直接操作 *****
// 添加日志生成器
// 获取日志生成器

// ***** 外部业务操作(安全沙盒) *****

// ***** Service API *****
func (m *mgr) Start() {
	go m.waitClose()
	for _, c := range m.set {
		c.Start()
	}
}

func (m *mgr) waitClose() {
	<-m.ctx.Done()
	m._close()
}

func (m *mgr) Close() {
	m.cancel()
}

func (m *mgr) _close() {
	log.Println("mgr Done")
}
