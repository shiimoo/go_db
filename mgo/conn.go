package mgo

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/shiimoo/godb/lib/base/errors"
	"github.com/shiimoo/godb/lib/savectrl"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewConn(parent context.Context, url string) (*MgoConn, error) {
	cOpts := options.Client().ApplyURI(url)
	conn := new(MgoConn)
	conn.ctx, conn.cancel = context.WithCancel(parent)
	client, err := mongo.Connect(conn.ctx, cOpts)
	if err != nil {
		return nil, err
	}
	err = client.Ping(conn.ctx, nil)
	if err != nil {
		return nil, err
	}
	conn.client = client
	conn.isOpen = false                // 默认关闭
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

	isOpen  bool              // 开启状态
	opFuncs map[string]opFunc // 指令集
	opChan  chan *op          // 操作参数 NewConn创建时会初始化1000的缓存
}

// ***** 数据业务操作 ****

// 数据库名转换
func (c *MgoConn) toDatabase(data any) string {
	var err error
	database, ok := data.(string)
	if !ok {
		err = errors.NewErr(
			ErrToDatabaseErr,
			fmt.Sprintf("need type string; but data[%v] type is %s", data, reflect.TypeOf(data)),
		)

	} else if len(strings.TrimSpace(database)) == 0 {
		err = errors.NewErr(
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
		err = errors.NewErr(
			ErrToCollectionErr,
			fmt.Sprintf("need type string; but data[%v] type is %s", data, reflect.TypeOf(data)),
		)

	} else if len(strings.TrimSpace(collection)) == 0 {
		err = errors.NewErr(
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
		err = errors.NewErr(ErrParamsErr, "params is nil")
	} else if len(params) < 2 {
		err = errors.NewErr(ErrParamsErr, fmt.Sprintf("params[%v] length < 2", params))
	}
	if err != nil {
		panic(err)
	}
	return c.toDatabase(params[0]), c.toCollection(params[1]), params[2:]
}

// 判定数据库database中是否存在集合collection
func (c *MgoConn) hasCollection(params ...any) *opResult {
	database, collection, _ := c.parseParams(params...)
	result := newOpResult()

	hasFlag := false
	list, err := c.client.Database(database).ListCollectionNames(c.ctx, bson.M{})
	if err == nil {
		for _, cName := range list {
			if cName == collection {
				hasFlag = true
			}
		}
		result.addResult(hasFlag)
	}
	return result
}

// 创建索引
func (c *MgoConn) createIndex(params ...any) *opResult {
	database, collection, params := c.parseParams(params...)
	result := newOpResult()

	indexs := parseCreateIndexParams(params...)
	set := c.client.Database(database).Collection(collection)
	if _, err := set.Indexes().CreateOne(c.ctx, mongo.IndexModel{
		Keys: indexs,
	}); err != nil {
		result.err = errors.NewErr(ErrMgoCreateIndexErr, err, database, collection, indexs)
	}
	return result
}

// 批量插入数据
func (c *MgoConn) insertN(params ...any) *opResult {
	database, collection, params := c.parseParams(params...)
	result := newOpResult()

	datas := parseInsertNParams(params...)
	set := c.client.Database(database).Collection(collection)
	_, err := set.InsertMany(c.ctx, datas)
	if err != nil {
		result.err = errors.NewErr(ErrMgoInsertErr, err, database, collection, datas)
	}

	return result
}

// 插入单挑数据
func (c *MgoConn) insertOne(params ...any) *opResult {
	database, collection, params := c.parseParams(params...)
	result := newOpResult()

	data := parseInsertOneParams(params...)
	set := c.client.Database(database).Collection(collection)
	_, err := set.InsertOne(c.ctx, data)
	if err != nil {
		result.err = errors.NewErr(ErrMgoInsertErr, err, database, collection, data)
	}

	return result
}

// 加载数据
func (c *MgoConn) find(params ...any) *opResult {
	database, collection, params := c.parseParams(params...)
	result := newOpResult()

	filter, num := parseFindParams(params...)
	opt := options.Find()
	if num > 0 {
		opt.SetLimit(num)
	}
	set := c.client.Database(database).Collection(collection)
	cur, err := set.Find(c.ctx, filter, opt)
	if err != nil {
		result.err = errors.NewErr(ErrMgoFindErr, err, database, collection, filter, num)
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

	filter := parseFindOneParams(params...)
	set := c.client.Database(database).Collection(collection)
	bs, err := set.FindOne(c.ctx, filter).Raw()
	if err != nil {
		result.err = errors.NewErr(ErrMgoFindErr, err, database, collection, filter, 1)
	} else {
		result.addResult([]byte(bs))
	}
	return result
}

// 根据mongo生成的ObjectId进行查找,等同于 findOne
func (c *MgoConn) findByObjId(params ...any) *opResult {
	database, collection, params := c.parseParams(params...)
	result := newOpResult()

	oId := parseFindByObjIdParams(params...)
	_id, err := primitive.ObjectIDFromHex(oId)
	if err != nil {
		result.err = errors.NewErr(ErrMgoObjectErr, oId, err)
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

	filter := parseDeleteParams(params...)
	set := c.client.Database(database).Collection(collection)
	delResult, err := set.DeleteMany(c.ctx, filter)
	if err != nil {
		result.err = errors.NewErr(ErrMgoDeleteErr, err, database, collection, filter)
	} else {
		result.addResult(int(delResult.DeletedCount))
	}
	return result
}

// 全部删除(清空数据)
func (c *MgoConn) deleteAll(params ...any) *opResult {
	database, collection, _ := c.parseParams(params...)
	return c.delete(database, collection)
}

// 单个删除
func (c *MgoConn) deleteOne(params ...any) *opResult {
	database, collection, params := c.parseParams(params...)
	result := newOpResult()

	filter := parseDeleteOneParams(params...)
	set := c.client.Database(database).Collection(collection)
	delRes, err := set.DeleteOne(c.ctx, filter)
	if err != nil {
		result.err = errors.NewErr(ErrMgoDeleteErr, err, database, collection, filter)
	} else {
		result.addResult(int(delRes.DeletedCount))
	}
	return result
}

// 根据mongo生成的ObjectId进行查找,等同于 deleteOne
func (c *MgoConn) deleteByObjId(params ...any) *opResult {
	database, collection, params := c.parseParams(params...)
	result := newOpResult()

	oId := parseDeleteByObjIdParams(params...)
	_id, err := primitive.ObjectIDFromHex(oId)
	if err != nil {
		result.err = errors.NewErr(ErrMgoObjectErr, oId, err)
		return result
	}

	filter := bson.D{
		{Key: "_id", Value: _id},
	}
	return c.deleteOne(database, collection, filter)
}

// 批量更新(选定范围内)
func (c *MgoConn) update(params ...any) *opResult {
	database, collection, params := c.parseParams(params...)
	result := newOpResult()

	filter, data := parseUpdateParams(params...)
	set := c.client.Database(database).Collection(collection)
	if _, err := set.UpdateMany(c.ctx, filter, bson.D{{Key: "$set", Value: data}}); err != nil {
		result.err = errors.NewErr(ErrMgoUpdateErr, err, database, collection, filter, data)
	}
	return result
}

// 单一更新
func (c *MgoConn) updateOne(params ...any) *opResult {
	database, collection, params := c.parseParams(params...)
	result := newOpResult()

	filter, data := parseUpdateOneParams(params...)
	set := c.client.Database(database).Collection(collection)
	if _, err := set.UpdateOne(c.ctx, filter, bson.D{{Key: "$set", Value: data}}); err != nil {
		result.err = errors.NewErr(ErrMgoUpdateErr, err, database, collection, filter, data)
	}
	return result
}

// 根据mongo生成的ObjectId进行更新,等同于 updateOne
func (c *MgoConn) updateByObjId(params ...any) *opResult {
	database, collection, params := c.parseParams(params...)
	result := newOpResult()

	oId := parseUpdateByObjIdParams(params...)
	_id, err := primitive.ObjectIDFromHex(oId)
	if err != nil {
		result.err = errors.NewErr(ErrMgoObjectErr, oId, err)
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

	filter, replacement := parseReplaceOneParams(params...)
	set := c.client.Database(database).Collection(collection)
	if _, err := set.ReplaceOne(c.ctx, filter, replacement); err != nil {
		result.err = errors.NewErr(ErrMgoReplaceErr, database, err, collection, filter, replacement)
	}
	return result
}

// 根据mongo生成的ObjectId进行更新,等同于 replaceOne
func (c *MgoConn) replaceByObjId(params ...any) *opResult {
	database, collection, params := c.parseParams(params...)
	result := newOpResult()

	oId := parseReplaceByObjIdParams(params...)
	_id, err := primitive.ObjectIDFromHex(oId)
	if err != nil {
		result.err = errors.NewErr(ErrMgoObjectErr, oId, err)
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
	c.isOpen = false // 设置关闭状态
	close(c.opChan)  // 关闭channel
	// 循环检出channel, 处理剩余指令
	for opObj := range c.opChan {
		if err := c.doOpHandler(opObj); err != nil {
			// todo 错误日志打印# 日志模块完善
			log.Println(err)
		}
	}
}

func (c *MgoConn) Start() {
	c.isOpen = true
	go c._start()
}

func (c *MgoConn) _start() {
	for {
		select {
		case <-c.ctx.Done():
			c.closeCallBack()
			return
		case opObj := <-c.opChan:
			if err := c.doOpHandler(opObj); err != nil {
				// todo 错误日志打印# 日志模块完善
				log.Println(err)
			}
		}
	}
}

// 执行指令
func (c *MgoConn) doOpHandler(opObj *op) error {
	return savectrl.SaveBox(func() error {
		handler, found := c.opFuncs[opObj.cmd]
		if found {
			opObj.resultAccept <- handler(opObj.args...)
		} else {
			return errors.NewErr(
				ErrToOpErr,
				fmt.Sprintf("op[%s] not found", opObj.cmd),
			)
		}
		return nil
	})
}

func (c *MgoConn) sendOp(opObj *op) error {
	if !c.isOpen {
		return ErrConnIsClose
	}
	c.opChan <- opObj
	return nil
}
