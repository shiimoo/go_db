package mgo

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/shiimoo/godb/dberr"
)

type mgr struct {
	sync.RWMutex
	ctx context.Context

	key string // 唯一标识(id)
	url string // mongo链接

	// 使用计数 每次获取链接自动+1
	useCount int
	// 链接池列表
	pool []*MgoConn
}

// Seturl 设置mongo链接参数
func (m *mgr) Seturl(host string, port int) {
	if host == "" {
		host = "localhost"
	}
	if port <= 0 {
		port = 27017
	}
	// TODO:用户密码等认证
	m.url = fmt.Sprintf("mongodb://%s:%d", host, port)
}

// Connect 创建链接, num为创建链接的数量;
// int error为nil时返回管理器的总conn计数；error~=nil时,返回成功创建的数量
// 创建失败的参数
func (m *mgr) Connect(num int) (int, error) {
	if num <= 0 {
		return 0, dberr.NewErr(dberr.ErrMgoConnNum, "num must > 0")
	}
	for i := 0; i < num; i++ {
		conn, err := NewConn(m.ctx, m.url)
		if err != nil {
			return i, dberr.NewErr(dberr.ErrMgoConnectErr, err)
		}
		m.AddConn(conn)
	}
	return m.Count(), nil
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
func (m *mgr) CreateIndex(database, collection string, indexs Indexs) {
	m.GetConn().CreateIndex(database, collection, indexs)
}

// InsertOne 插入单挑数据
func (m *mgr) InsertOne(database, collection string, data any) error {
	return m.GetConn().InsertOne(database, collection, data)
}

// ----------------

// 工厂方法, key值在外层进行检查和校准
func newMgr(parent context.Context, key string) *mgr {
	// ctx, cancel TODO:关闭监听待实现
	ctx, _ := context.WithCancel(parent)
	mgr := new(mgr)
	mgr.ctx = ctx
	mgr.key = key
	mgr.pool = make([]*MgoConn, 0)
	return mgr
}

// mgr池
var mp struct {
	sync.Map // 管理器安全池
}

// GetMgr 获取指定Mgr，无则创建;
// 其中parent上下文只有在创建的时会使用, 仅获取时无意义
func GetMgr(parent context.Context, key string) (*mgr, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, dberr.NewErr(dberr.ErrMgoMgrKey, "key is \"\"")
	}
	var kMgr *mgr
	v, ok := mp.Load(key)
	if !ok {
		kMgr = newMgr(parent, key) // 创建
		mp.Store(key, kMgr)        // 添加集合
	} else {
		kMgr = v.(*mgr) // 转换
	}
	return kMgr, nil
}
