package mgo

import (
	"sync"
)

type mgr struct {
	sync.RWMutex

	// 使用计数 每次获取链接自动+1
	useCount int
	// 链接池列表
	pool []*MgoConn
}

func (m *mgr) Count() int {
	return len(m.pool)
}

func (m *mgr) AddConn(conn *MgoConn) {
	m.Lock()
	defer m.Unlock()
	m.pool = append(m.pool, conn)
}

func (m *mgr) GetConn() *MgoConn {
	m.RLock()
	defer m.RUnlock()
	m.useCount += 1
	index := m.useCount % m.Count()
	return m.pool[index]
}

// HasCollection 判定数据库database中是否存在集合collection
func (m *mgr) HasCollection(database, collection string) bool {
	return m.GetConn().hasCollection(database, collection)
}

// CreateIndex 建立索引
func (m *mgr) CreateIndex(database, collection string, index MgoIndex) {
	m.GetConn().CreateIndex(database, collection, index)
}

// ---------------- TODO: 工厂方法迁移

// 工厂方法
func newMgr() *mgr {
	mgr := new(mgr)
	mgr.pool = make([]*MgoConn, 0)
	return mgr
}

// 单例
var soleMgr *mgr = newMgr()

func getMgr() *mgr {
	if soleMgr == nil {
		soleMgr = newMgr()
	}
	return soleMgr
}
