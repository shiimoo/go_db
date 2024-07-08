package mgo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/shiimoo/godb/dberr"
)

type mgr struct {
	sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc

	key string // 唯一标识(id)
	url string // mongo链接

	// 使用计数 每次获取链接自动+1
	useCount int
	// 链接池列表
	pool []*MgoConn
}

// ***** 基础直接操作 *****

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
		return 0, dberr.NewErr(ErrMgoConnNum, "num must > 0")
	}
	for i := 0; i < num; i++ {
		conn, err := NewConn(m.ctx, m.url)
		if err != nil {
			return i, dberr.NewErr(ErrMgoConnectErr, err)
		}
		m.addConn(conn)
	}
	return m.Count(), nil
}

func (m *mgr) Count() int {
	return len(m.pool)
}

func (m *mgr) addConn(conn *MgoConn) {
	m.Lock()
	defer m.Unlock()
	m.pool = append(m.pool, conn)
}

func (m *mgr) getConn() *MgoConn {
	m.RLock()
	defer m.RUnlock()
	m.useCount += 1
	index := m.useCount % m.Count()
	return m.pool[index]
}

// ***** 外部业务操作(安全沙盒) *****

// HasCollection 判定数据库database中是否存在集合collection
func (m *mgr) HasCollection(database, collection string) bool {
	return m.getConn().hasCollection(database, collection)
}

// CreateIndex 建立索引
func (m *mgr) CreateIndex(database, collection string, indexs Indexs) {
	m.getConn().CreateIndex(database, collection, indexs)
}

// InsertOne 插入单挑数据
func (m *mgr) InsertOne(database, collection string, data any) error {
	return m.getConn().InsertOne(database, collection, data)
}

// InsertN 批量插入数据
func (m *mgr) InsertN(database, collection string, datas []any) error {
	return m.getConn().InsertN(database, collection, datas)
}

// Find 加载数据 filter一般是bson.D
func (m *mgr) Find(database, collection string, filter any, num int64) ([][]byte, error) {
	return m.getConn().Find(database, collection, filter, num)
}

// FindAll 查找全部数据
func (m *mgr) FindAll(database, collection string) ([][]byte, error) {
	return m.getConn().FindAll(database, collection)
}

// FindOne 查找单个数据
func (m *mgr) FindOne(database, collection string, filter any) ([]byte, error) {
	return m.getConn().FindOne(database, collection, filter)
}

// FindByObjId 根据mongo生成的ObjectId进行查找,等同于FindOne
func (m *mgr) FindByObjId(database, collection, oId string) ([]byte, error) {
	return m.getConn().FindByObjId(database, collection, oId)
}

// Delete 删除数据
func (m *mgr) Delete(database, collection string, filter any) (int, error) {
	return m.getConn().Delete(database, collection, filter)
}

// DeleteAll 删除全部数据(清空数据)
func (m *mgr) DeleteAll(database, collection string, filter any) (int, error) {
	return m.getConn().DeleteAll(database, collection, filter)
}

// DeleteOne 删除数据(单个)
func (m *mgr) DeleteOne(database, collection string, filter any) (int, error) {
	return m.getConn().DeleteOne(database, collection, filter)
}

// FindByObjId 根据mongo生成的ObjectId进行查找,等同于FindOne
func (m *mgr) DeleteByObjId(database, collection, oId string) (int, error) {
	return m.getConn().DeleteByObjId(database, collection, oId)
}

// 增查 ：改
// UpdateOne 更新数据(单个)
func (m *mgr) UpdateOne(database, collection string, filter, update any) error {
	return m.getConn().UpdateOne(database, collection, filter, update)
}

// UpdateByObjId 根据mongo生成的ObjectId进行更新
func (m *mgr) UpdateByObjId(database, collection, oId string, data any) error {
	return m.getConn().UpdateByObjId(database, collection, oId, data)
}

// UpdateOne 更新数据(单个)
func (m *mgr) Update(database, collection string, filter, update any) error {
	return m.getConn().Update(database, collection, filter, update)
}

// ReplaceOne 整体替换
func (m *mgr) ReplaceOne(database, collection string, filter, replacement any) error {
	return m.getConn().ReplaceOne(database, collection, filter, replacement)
}

// ***** Service API *****

func (m *mgr) Start() {
	for {

	}
	go m.waitClose()
	// 消息队列
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

// 池管理

// 工厂方法, key值在外层进行检查和校准
func newMgr(parent context.Context, key string) *mgr {
	ctx, cancel := context.WithCancel(parent)
	mgr := new(mgr)
	mgr.ctx = ctx
	mgr.cancel = cancel
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
		return nil, dberr.NewErr(ErrMgoMgrKey, "key is \"\"")
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
