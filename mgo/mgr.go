package mgo

import (
	"sync"
)

type mgr struct {
	sync.RWMutex

	useCount int
	pool     []*MgoConn
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
	index := m.useCount % m.Count()
	return m.pool[index]
}

func newMgr() *mgr {
	mgr := new(mgr)
	mgr.pool = make([]*MgoConn, 0)
	return mgr
}
