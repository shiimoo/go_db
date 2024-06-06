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
	// 指定数据库名
	database string
}

func (m *mgr) SetDatabase(database string) {
	m.database = database
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

func newMgr() *mgr {
	mgr := new(mgr)
	mgr.pool = make([]*MgoConn, 0)
	return mgr
}

var soleMgr *mgr = newMgr()

func getMgr() *mgr {
	if soleMgr == nil {
		soleMgr = newMgr()
	}
	return soleMgr
}
