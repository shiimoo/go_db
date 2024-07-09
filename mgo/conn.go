package mgo

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/shiimoo/godb/dberr"
	"github.com/shiimoo/godb/lib/savectrl"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	cmdHasCollection = "hasCollection" // 判定数据库database中是否存在集合collectio
	cmdCreateIndex   = "createIndex"   // 创建索引

	cmdInsertN   = "insertN"   // 批量插入数据
	cmdInsertOne = "insertOne" // 插入单挑数据

	cmdFind        = "find"        // 加载数据
	cmdFindAll     = "findAll"     // 加载全部数据
	cmdFindOne     = "findOne"     // 加载单个数据
	cmdFindByObjId = "findByObjId" // 根据mongo生成的ObjectId进行查找,等同于findOne

	cmdDelete        = "delete"        // 删除数据
	cmdDeleteAll     = "deleteAll"     // 全部删除(清空)
	cmdDeleteOne     = "deleteOne"     // 删除1个
	cmdDeleteByObjId = "deleteByObjId" // 根据mongo生成的ObjectId进行查找,等同于deleteOne

	cmdUpdate        = "update"        // 批量更新
	cmdUpdateOne     = "updateOne"     // 单一更新
	cmdUpdateByObjId = "updateByObjId" // 根据mongo生成的ObjectId进行更新,等同于updateOne

	cmdReplaceOne     = "replaceOne"     // 整个文档内容替换(除了ObjectId)
	cmdReplaceByObjId = "replaceByObjId" // 根据mongo生成的ObjectId进行更新,等同于replaceOne
)

type opFunc func(otherParams ...any) *opResult

func NewConn(parent context.Context, url string) (*MgoConn, error) {
	cOpts := options.Client().ApplyURI(url)
	ctx, cancel := context.WithCancel(parent)
	client, err := mongo.Connect(ctx, cOpts)
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	conn := new(MgoConn)
	conn.ctx = ctx
	conn.cancel = cancel
	conn.client = client
	conn.opChan = make(chan *op, 1000) // 长度为1000的处理队列
	// 注册处理方法
	conn.opFuncs = map[string]opFunc{
		cmdHasCollection:  conn.hasCollection,
		cmdCreateIndex:    conn.createIndex,
		cmdInsertN:        conn.insertN,
		cmdInsertOne:      conn.insertOne,
		cmdFind:           conn.find,
		cmdFindAll:        conn.findAll,
		cmdFindOne:        conn.findOne,
		cmdFindByObjId:    conn.findByObjId,
		cmdDelete:         conn.delete,
		cmdDeleteAll:      conn.deleteAll,
		cmdDeleteOne:      conn.deleteOne,
		cmdDeleteByObjId:  conn.deleteByObjId,
		cmdUpdate:         conn.update,
		cmdUpdateOne:      conn.updateOne,
		cmdUpdateByObjId:  conn.updateByObjId,
		cmdReplaceOne:     conn.replaceOne,
		cmdReplaceByObjId: conn.replaceByObjId,
	}
	return conn, nil
}

type MgoConn struct {
	ctx    context.Context
	cancel context.CancelFunc
	client *mongo.Client

	opFuncs map[string]opFunc
	opChan  chan *op // 操作参数
}

// ***** 数据操作 ****

// 判定数据库database中是否存在集合collection
func (c *MgoConn) _hasCollection(database, collection string, _ ...any) bool {
	list, err := c.client.Database(database).ListCollectionNames(c.ctx, bson.M{})
	if err != nil {
		return false
	}
	for _, cName := range list {
		if cName == collection {
			return true
		}
	}
	return false
}

// ***** 数据业务操作 ****

// 数据库名转换
func (c *MgoConn) toDatabase(data any) string {
	var err error
	database, ok := data.(string)
	if !ok {
		err = dberr.NewErr(
			ErrToDatabaseErr,
			fmt.Sprintf("need type string; but data[%v] type is %s", data, reflect.TypeOf(data)),
		)

	} else if len(strings.TrimSpace(database)) == 0 {
		err = dberr.NewErr(
			ErrToDatabaseErr,
			fmt.Sprintf("data[%s] length is zero", data),
		)
	}
	if err != nil {
		panic(err)
	}
	return database
}

// 集合名转换
func (c *MgoConn) toCollection(data any) string {
	var err error
	collection, ok := data.(string)
	if !ok {
		err = dberr.NewErr(
			ErrToCollectionErr,
			fmt.Sprintf("need type string; but data[%v] type is %s", data, reflect.TypeOf(data)),
		)

	} else if len(strings.TrimSpace(collection)) == 0 {
		err = dberr.NewErr(
			ErrToCollectionErr,
			fmt.Sprintf("data[%s] length is zero", data),
		)
	}
	if err != nil {
		panic(err)
	}
	return collection
}

// 参数解析
func (c *MgoConn) parseParams(params ...any) (database string, collection string, args []any) {
	var err error
	if params == nil {
		err = dberr.NewErr(ErrParamsErr, "params is nil")
	} else if len(params) < 2 {
		err = dberr.NewErr(ErrParamsErr, fmt.Sprintf("params[%v] length < 2", params))
	}
	if err != nil {
		panic(err)
	}
	return c.toDatabase(params[0]), c.toCollection(params[1]), params[2:]
}

// 判定数据库database中是否存在集合collection
func (c *MgoConn) hasCollection(params ...any) *opResult {
	database, collection, params := c.parseParams(params...)
	result := newOpResult()
	result.addResult(c._hasCollection(database, collection, params...))
	return result
}

// 创建索引
func (c *MgoConn) createIndex(params ...any) *opResult {
	database, collection, params := c.parseParams(params...)
	result := newOpResult()

	indexs := params[0]
	set := c.client.Database(database).Collection(collection)
	// todo 去除无用打印
	fmt.Println(set.Indexes().CreateOne(c.ctx, mongo.IndexModel{
		Keys: indexs,
	}))

	return result
}

// 批量插入数据
func (c *MgoConn) insertN(params ...any) *opResult {
	database, collection, params := c.parseParams(params...)
	result := newOpResult()

	datas := params[0].([]any)
	set := c.client.Database(database).Collection(collection)
	_, err := set.InsertMany(c.ctx, datas)
	if err != nil {
		result.err = dberr.NewErr(ErrMgoInsertErr, err, database, collection, datas)
	}

	return result
}

// 插入单挑数据
func (c *MgoConn) insertOne(params ...any) *opResult {
	database, collection, params := c.parseParams(params...)
	result := newOpResult()

	data := params[0]
	set := c.client.Database(database).Collection(collection)
	_, err := set.InsertOne(c.ctx, data)
	if err != nil {
		result.err = dberr.NewErr(ErrMgoInsertErr, err, database, collection, data)
	}

	return result
}

// 加载数据
func (c *MgoConn) find(params ...any) *opResult {
	database, collection, params := c.parseParams(params...)
	result := newOpResult()

	var filter any
	if len(params) > 0 {
		filter = params[0]
	}
	var num int64
	if len(params) > 1 {
		num = params[1].(int64)
	}
	if filter == nil {
		filter = bson.D{}
	}
	opt := options.Find()
	if num > 0 {
		opt.SetLimit(num)
	}
	set := c.client.Database(database).Collection(collection)
	cur, err := set.Find(c.ctx, filter, opt)
	if err != nil {
		result.err = dberr.NewErr(ErrMgoFindErr, err, database, collection, filter, num)
	} else {
		datas := make([][]byte, 0)
		for cur.Next(c.ctx) {
			datas = append(datas, []byte(cur.Current))
		}
		result.addResult(datas)
	}

	return result
}

// 加载全部数据
func (c *MgoConn) findAll(params ...any) *opResult {
	database, collection, _ := c.parseParams(params...)
	return c.find(database, collection, bson.D{}, int64(0))
}

// 查找单个数据
func (c *MgoConn) findOne(params ...any) *opResult {
	database, collection, params := c.parseParams(params...)
	result := newOpResult()

	var filter any
	if len(params) > 0 {
		filter = params[0]
	}
	if filter == nil {
		filter = bson.D{}
	}
	set := c.client.Database(database).Collection(collection)
	bs, err := set.FindOne(c.ctx, filter).Raw()
	if err != nil {
		result.err = dberr.NewErr(ErrMgoFindErr, err, database, collection, filter, 1)
	} else {
		result.addResult([]byte(bs))
	}
	return result
}

// FindByObjId 根据mongo生成的ObjectId进行查找,等同于FindOne
func (c *MgoConn) findByObjId(params ...any) *opResult {
	database, collection, params := c.parseParams(params...)
	result := newOpResult()

	oId := params[0].(string)
	_id, err := primitive.ObjectIDFromHex(oId)
	if err != nil {
		result.err = dberr.NewErr(ErrMgoObjectErr, oId, err)
		return result
	}

	filter := bson.D{
		{Key: "_id", Value: _id},
	}
	return c.findOne(database, collection, filter)
}

// 删除
func (c *MgoConn) delete(params ...any) *opResult {
	database, collection, params := c.parseParams(params...)
	result := newOpResult()

	var filter any
	if len(params) > 0 {
		filter = params[0]
	}
	if filter == nil {
		filter = bson.D{}
	}
	set := c.client.Database(database).Collection(collection)
	delResult, err := set.DeleteMany(c.ctx, filter)
	if err != nil {
		result.err = dberr.NewErr(ErrMgoDeleteErr, err, database, collection, filter)
	} else {
		result.addResult(int(delResult.DeletedCount))
	}
	return result
}

// DeletaAll 全部删除(清空)
func (c *MgoConn) deleteAll(params ...any) *opResult {
	database, collection, _ := c.parseParams(params...)
	return c.delete(database, collection)
}

func (c *MgoConn) deleteOne(params ...any) *opResult {
	database, collection, params := c.parseParams(params...)
	result := newOpResult()

	var filter any
	if len(params) > 0 {
		filter = params[0]
	}
	if filter == nil {
		filter = bson.D{}
	}
	set := c.client.Database(database).Collection(collection)
	delRes, err := set.DeleteOne(c.ctx, filter)
	if err != nil {
		result.err = dberr.NewErr(ErrMgoDeleteErr, err, database, collection, filter)
	} else {
		result.addResult(int(delRes.DeletedCount))
	}
	return result
}

// DeleteByObjId 根据mongo生成的ObjectId进行查找,等同于DeleteOne
func (c *MgoConn) deleteByObjId(params ...any) *opResult {

	database, collection, params := c.parseParams(params...)
	result := newOpResult()

	oId := params[0].(string)
	_id, err := primitive.ObjectIDFromHex(oId)
	if err != nil {
		result.err = dberr.NewErr(ErrMgoObjectErr, oId, err)
		return result
	}

	filter := bson.D{
		{Key: "_id", Value: _id},
	}
	return c.deleteOne(database, collection, filter)
}

// UpdateN 批量更新
func (c *MgoConn) update(params ...any) *opResult {
	database, collection, params := c.parseParams(params...)
	result := newOpResult()

	var filter any
	if len(params) > 0 {
		filter = params[0]
	}
	if filter == nil {
		filter = bson.D{}
	}
	var data any
	if len(params) > 1 {
		data = params[1]
	}
	set := c.client.Database(database).Collection(collection)
	if _, err := set.UpdateMany(c.ctx, filter, bson.D{{Key: "$set", Value: data}}); err != nil {
		result.err = dberr.NewErr(ErrMgoUpdateErr, err, database, collection, filter, data)
	}
	return result
}

// 单一更新
func (c *MgoConn) updateOne(params ...any) *opResult {
	database, collection, params := c.parseParams(params...)
	result := newOpResult()

	var filter any
	if len(params) > 0 {
		filter = params[0]
	}
	if filter == nil {
		filter = bson.D{}
	}
	var data any
	if len(params) > 1 {
		data = params[1]
	}
	set := c.client.Database(database).Collection(collection)
	if _, err := set.UpdateOne(c.ctx, filter, bson.D{{Key: "$set", Value: data}}); err != nil {
		result.err = dberr.NewErr(ErrMgoUpdateErr, err, database, collection, filter, data)
	}
	return result
}

// 根据mongo生成的ObjectId进行更新,等同于UpdateOne
func (c *MgoConn) updateByObjId(params ...any) *opResult {
	database, collection, params := c.parseParams(params...)
	result := newOpResult()

	oId := params[0].(string)
	_id, err := primitive.ObjectIDFromHex(oId)
	if err != nil {
		result.err = dberr.NewErr(ErrMgoObjectErr, oId, err)
		return result
	}

	filter := bson.D{
		{Key: "_id", Value: _id},
	}
	return c.updateOne(database, collection, filter, params[1])
}

// 整个文档内容替换(除了ObjectId)
func (c *MgoConn) replaceOne(params ...any) *opResult {
	database, collection, params := c.parseParams(params...)
	result := newOpResult()

	var filter any
	if len(params) > 0 {
		filter = params[0]
	}
	if filter == nil {
		filter = bson.D{}
	}
	var replacement any
	if len(params) > 1 {
		replacement = params[1]
	}
	set := c.client.Database(database).Collection(collection)

	if _, err := set.ReplaceOne(c.ctx, filter, replacement); err != nil {
		result.err = dberr.NewErr(ErrMgoReplaceErr, database, err, collection, filter, replacement)
	}
	return result
}

// 根据mongo生成的ObjectId进行更新,等同于ReplaceOne
func (c *MgoConn) replaceByObjId(params ...any) *opResult {
	database, collection, params := c.parseParams(params...)
	result := newOpResult()

	oId := params[0].(string)
	_id, err := primitive.ObjectIDFromHex(oId)
	if err != nil {
		result.err = dberr.NewErr(ErrMgoObjectErr, oId, err)
		return result
	}

	filter := bson.D{
		{Key: "_id", Value: _id},
	}
	return c.replaceOne(database, collection, filter, params[1])
}

// ***** ServerAPI

func (c *MgoConn) Close() {
	c.cancel()
}

func (c *MgoConn) closeCallBack() {
	// todo 链接关闭回调
}

func (c *MgoConn) Start() {
	go c._start()
}

func (c *MgoConn) _start() {
	for {
		select {
		case <-c.ctx.Done():
			c.closeCallBack()
			return
		case opObj := <-c.opChan:
			if err := savectrl.SaveBox(func() error {
				handler, found := c.opFuncs[opObj.cmd]
				if found {
					opObj.resultAccept <- handler(opObj.args...)
				} else {
					return dberr.NewErr(
						ErrToOpErr,
						fmt.Sprintf("op[%s] not found", opObj.cmd),
					)
				}
				return nil
			}); err != nil {
				// todo 错误日志打印
				log.Println(err)
			}
		}
	}
}

func (c *MgoConn) doOp(opObj *op) {
	c.opChan <- opObj
}
