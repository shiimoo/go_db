package mlog

import (
	"context"
	"log"
	"sync"
)

// ***** 生成管理器 *****

type mgr struct {
	ctx    context.Context
	cancel context.CancelFunc

	pool sync.Map // 管理器安全池
}

// ***** 基础直接操作 *****

// 添加日志生成器
func (m *mgr) addLoger(l *logger) {
	l.Start()
	m.pool.Store(l.key, l)
}

// 获取日志生成器
func (m *mgr) getLoger(key string) *logger {
	var l *logger
	val, ok := m.pool.Load(key)
	if !ok {
		l = newLogger(m.ctx, key)
		m.addLoger(l)
	} else {
		l = val.(*logger)
	}
	return l
}

// ***** 外部业务操作(安全沙盒) *****
func (m *mgr) SetOutFunc(key string, handler func(l *Log) error) {
	m.getLoger(key).SetOutFunc(handler)
}

// Println 默认输出
func (m *mgr) Println(msg string) {
	m.getLoger("server").Output(Info, nil, msg)
}

// ***** Service API *****
func (m *mgr) Start() {
	go m.waitClose()
	m.pool.Range(func(key, value any) bool {
		l := value.(*logger)
		l.Start()
		return true
	})
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
	m.pool.Range(func(key, value any) bool {
		l := value.(*logger)
		l.Close()
		return true
	})
}

// ***** mgr 工厂方法(全局单例) *****

var uniqueMgr *mgr

func MgrInit(parent context.Context) {
	uniqueMgr = new(mgr)
	uniqueMgr.ctx, uniqueMgr.cancel = context.WithCancel(parent)
	uniqueMgr.Start()
}

func GetMgr() *mgr {
	return uniqueMgr
}
