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
	// todo 分配规则优化
	m.useCount += 1
	index := m.useCount % m.Count()
	return m.pool[index]
}

// ***** 外部业务操作(安全沙盒) *****

func (m *mgr) newOp(cmd, database, collection string) *op {
	opObj := newOp(cmd)
	opObj.append(database, collection)
	return opObj
}

// HasCollection 判定数据库database中是否存在集合collection
func (m *mgr) HasCollection(database, collection string) bool {
	opObj := m.newOp(cmdHasCollection, database, collection)
	m.getConn().doOp(opObj)
	result := opObj.getResult()
	// todo 数据转型
	return result.results[0].(bool)
}

// CreateIndex 建立索引
func (m *mgr) CreateIndex(database, collection string, indexs Indexs) {
	opObj := m.newOp(cmdCreateIndex, database, collection)
	opObj.append(indexs)
	m.getConn().doOp(opObj)
	_ = opObj.getResult()
	// todo 数据转型
}

// InsertN 批量插入数据
func (m *mgr) InsertN(database, collection string, datas []any) error {
	opObj := m.newOp(cmdInsertN, database, collection)
	opObj.append(datas)
	m.getConn().doOp(opObj)
	result := opObj.getResult()
	// todo 数据转型
	return result.err
}

// InsertOne 插入单挑数据
func (m *mgr) InsertOne(database, collection string, data any) error {
	opObj := m.newOp(cmdInsertOne, database, collection)
	opObj.append(data)
	m.getConn().doOp(opObj)
	result := opObj.getResult()
	// todo 数据转型
	return result.err
}

// Find 加载数据 filter一般是bson.D
func (m *mgr) Find(database, collection string, filter any, num int64) ([][]byte, error) {

	opObj := m.newOp(cmdFind, database, collection)
	opObj.append(filter, num)
	m.getConn().doOp(opObj)
	result := opObj.getResult()
	// todo 数据转型
	return result.results[0].([][]byte), result.err
}

// FindAll 查找全部数据
func (m *mgr) FindAll(database, collection string) ([][]byte, error) {

	opObj := m.newOp(cmdFindAll, database, collection)
	m.getConn().doOp(opObj)
	result := opObj.getResult()
	// todo 数据转型
	return result.results[0].([][]byte), result.err
}

// FindOne 查找单个数据
func (m *mgr) FindOne(database, collection string, filter any) ([]byte, error) {
	opObj := m.newOp(cmdFindOne, database, collection)
	opObj.append(filter)
	m.getConn().doOp(opObj)
	result := opObj.getResult()
	// todo 数据转型
	return (result.results[0]).([]byte), result.err
}

// FindByObjId 根据mongo生成的ObjectId进行查找,等同于FindOne
func (m *mgr) FindByObjId(database, collection, oId string) ([]byte, error) {
	opObj := m.newOp(cmdFindByObjId, database, collection)
	opObj.append(oId)
	m.getConn().doOp(opObj)
	result := opObj.getResult()
	// todo 数据转型
	return result.results[0].([]byte), result.err
}

// Delete 删除数据
func (m *mgr) Delete(database, collection string, filter any) (int, error) {

	opObj := m.newOp(cmdDelete, database, collection)
	opObj.append(filter)
	m.getConn().doOp(opObj)
	result := opObj.getResult()
	// todo 数据转型
	return result.results[0].(int), result.err
}

// DeleteAll 删除全部数据(清空数据)
func (m *mgr) DeleteAll(database, collection string, filter any) (int, error) {
	opObj := m.newOp(cmdDeleteAll, database, collection)
	opObj.append(filter)
	m.getConn().doOp(opObj)
	result := opObj.getResult()
	// todo 数据转型
	return result.results[0].(int), result.err
}

// DeleteOne 删除数据(单个)
func (m *mgr) DeleteOne(database, collection string, filter any) (int, error) {
	opObj := m.newOp(cmdDeleteOne, database, collection)
	opObj.append(filter)
	m.getConn().doOp(opObj)
	result := opObj.getResult()
	// todo 数据转型
	return result.results[0].(int), result.err
}

// DeleteByObjId 根据mongo生成的ObjectId进行查找,等同于DeleteOne
func (m *mgr) DeleteByObjId(database, collection, oId string) (int, error) {
	opObj := m.newOp(cmdDeleteByObjId, database, collection)
	opObj.append(oId)
	m.getConn().doOp(opObj)
	result := opObj.getResult()
	// todo 数据转型
	return result.results[0].(int), result.err
}

// Update 更新数据
func (m *mgr) Update(database, collection string, filter, data any) error {
	opObj := m.newOp(cmdUpdate, database, collection)
	opObj.append(filter, data)
	m.getConn().doOp(opObj)
	result := opObj.getResult()
	// todo 数据转型
	return result.err
}

// UpdateOne 更新数据(单个)
func (m *mgr) UpdateOne(database, collection string, filter, data any) error {
	opObj := m.newOp(cmdUpdateOne, database, collection)
	opObj.append(filter, data)
	m.getConn().doOp(opObj)
	result := opObj.getResult()
	// todo 数据转型
	return result.err
}

// UpdateByObjId 根据mongo生成的ObjectId进行更新
func (m *mgr) UpdateByObjId(database, collection, oId string, data any) error {
	opObj := m.newOp(cmdUpdateByObjId, database, collection)
	opObj.append(oId, data)
	m.getConn().doOp(opObj)
	result := opObj.getResult()
	// todo 数据转型
	return result.err
}

// ReplaceOne 整体替换
func (m *mgr) ReplaceOne(database, collection string, filter, replacement any) error {
	opObj := m.newOp(cmdReplaceOne, database, collection)
	opObj.append(filter, replacement)
	m.getConn().doOp(opObj)
	result := opObj.getResult()
	// todo 数据转型
	return result.err
}

// ReplaceOne 整体替换
func (m *mgr) ReplaceByObjId(database, collection string, oId string, replacement any) error {
	opObj := m.newOp(cmdReplaceByObjId, database, collection)
	opObj.append(oId, replacement)
	m.getConn().doOp(opObj)
	result := opObj.getResult()
	// todo 数据转型
	return result.err
}

// ***** Service API *****

func (m *mgr) Start() {
	go m.waitClose()
	for _, c := range m.pool {
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
