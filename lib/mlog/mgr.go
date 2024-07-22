package mlog

import (
	"context"
	"log"
	"strings"
	"sync"
)

const (
	DefaultLoggerName = "server"
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
	if strings.ContainsAny(key, "\\/:*?\"<>|") {
		// todo warn 不得包含字符()
		key = DefaultLoggerName
	}
	key = strings.TrimSpace(key)
	if key == "" {
		key = DefaultLoggerName
	} else {
		if strings.ContainsAny(key, "\\/:*?\"<>|") {
			// todo warn 不得包含字符()
			key = DefaultLoggerName
		}
	}
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
	m.getLoger(DefaultLoggerName).Output(Info, "", msg)
}

func (m *mgr) Debug(mod, label, msg string) {
	m.Debugf(mod, label, msg)
}

func (m *mgr) Debugf(mod, label, format string, datas ...Data) {
	m.getLoger(mod).Outputf(Debug, label, format, datas...)
}

func (m *mgr) Info(mod, label, msg string) {
	m.Infof(mod, label, msg)
}

func (m *mgr) Infof(mod, label, format string, datas ...Data) {
	m.getLoger(mod).Outputf(Info, label, format, datas...)
}

func (m *mgr) Warn(mod, label, msg string) {
	m.Warnf(mod, label, msg)
}

func (m *mgr) Warnf(mod, label, format string, datas ...Data) {
	m.getLoger(mod).Outputf(Warn, label, format, datas...)
}

func (m *mgr) Error(mod, label, msg string) {
	m.Errorf(mod, label, msg)
}

func (m *mgr) Errorf(mod, label, format string, datas ...Data) {
	m.getLoger(mod).Outputf(Error, label, format, datas...)
}

func (m *mgr) Fatal(mod, label, msg string) {
	m.Errorf(mod, label, msg)
}

func (m *mgr) Fatalf(mod, label, format string, datas ...Data) {
	m.getLoger(mod).Outputf(Fatal, label, format, datas...)
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
