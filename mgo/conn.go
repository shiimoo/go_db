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
	hasCollection = "hasCollection" // 是否拥有集合
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
		hasCollection: conn.hasCollection,
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

// ***** 数据业务操作 ****

// 操作方法索引转换
func (c *MgoConn) toOp(data any) string {
	var err error
	op, ok := data.(string)
	if !ok {
		err = dberr.NewErr(
			ErrToOpErr,
			fmt.Sprintf("need type string; but op[%v] type is %s", data, reflect.TypeOf(data)),
		)
	} else if len(strings.TrimSpace(op)) == 0 {
		err = dberr.NewErr(
			ErrToOpErr,
			fmt.Sprintf("op[%s] length is zero", data),
		)
	} else if _, found := c.opFuncs[op]; !found {
		err = dberr.NewErr(
			ErrToOpErr,
			fmt.Sprintf("op[%s] not found", data),
		)
	}
	if err != nil {
		panic(err)
	}
	return op
}

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
	return c.toDatabase(params[0]), c.toDatabase(params[1]), params[2:]
}

// 判定数据库database中是否存在集合collection
func (c *MgoConn) hasCollection(otherParams ...any) *opResult {
	database, collection, otherParams := c.parseParams(otherParams...)
	result := newOpResult()
	result.addResult(c._hasCollection(database, collection, otherParams...))
	return result
}

// 判定数据库database中是否存在集合collection
func (c *MgoConn) _hasCollection(database, collection string, params ...any) bool {
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

// 创建索引
func (c *MgoConn) CreateIndex(database, collection string, indexs Indexs) {
	set := c.client.Database(database).Collection(collection)
	fmt.Println(set.Indexes().CreateOne(c.ctx, mongo.IndexModel{
		Keys: indexs,
	}))
}

// InsertOne 插入单挑数据
func (c *MgoConn) InsertOne(database, collection string, data any) error {
	set := c.client.Database(database).Collection(collection)
	_, err := set.InsertOne(c.ctx, data)
	if err != nil {
		return dberr.NewErr(ErrMgoInsertErr, err, database, collection, data)
	}
	return nil
}

// InsertN 批量插入数据
func (c *MgoConn) InsertN(database, collection string, datas []any) error {
	set := c.client.Database(database).Collection(collection)
	_, err := set.InsertMany(c.ctx, datas)
	if err != nil {
		return dberr.NewErr(ErrMgoInsertErr, err, database, collection, datas)
	}
	return nil
}

// Find 加载数据 filter一般是bson.D
func (c *MgoConn) Find(database, collection string, filter any, num int64) ([][]byte, error) {
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
		return nil, dberr.NewErr(ErrMgoFindErr, err, database, collection, filter, num)
	}

	result := make([][]byte, 0)
	for cur.Next(c.ctx) {
		result = append(result, cur.Current)
	}
	return result, nil
}

// FindAll 全部加载
func (c *MgoConn) FindAll(database, collection string) ([][]byte, error) {
	return c.Find(database, collection, bson.D{}, 0)
}

// FindOne 查找单个数据
func (c *MgoConn) FindOne(database, collection string, filter any) ([]byte, error) {
	if filter == nil {
		filter = bson.D{}
	}
	set := c.client.Database(database).Collection(collection)
	bs, err := set.FindOne(c.ctx, filter).Raw()
	if err != nil {
		return nil, dberr.NewErr(ErrMgoFindErr, err, database, collection, filter, 1)
	}
	return bs, nil
}

// FindByObjId 根据mongo生成的ObjectId进行查找,等同于FindOne
func (c *MgoConn) FindByObjId(database, collection, oId string) ([]byte, error) {
	_id, err := primitive.ObjectIDFromHex(oId)
	if err != nil {
		return nil, dberr.NewErr(ErrMgoObjectErr, oId, err)
	}

	filter := bson.D{
		{Key: "_id", Value: _id},
	}
	return c.FindOne(database, collection, filter)
}

// Delete 删除
func (c *MgoConn) Delete(database, collection string, filter any) (int, error) {
	if filter == nil {
		filter = bson.D{}
	}
	set := c.client.Database(database).Collection(collection)
	result, err := set.DeleteMany(c.ctx, filter)
	if err != nil {
		return 0, dberr.NewErr(ErrMgoDeleteErr, err, database, collection, filter)
	}
	return int(result.DeletedCount), nil
}

// DeletaAll 全部删除(清空)
func (c *MgoConn) DeleteAll(database, collection string, filter any) (int, error) {
	return c.Delete(database, collection, nil)
}

func (c *MgoConn) DeleteOne(database, collection string, filter any) (int, error) {
	if filter == nil {
		filter = bson.D{}
	}
	set := c.client.Database(database).Collection(collection)
	result, err := set.DeleteOne(c.ctx, filter)
	if err != nil {
		return 0, dberr.NewErr(ErrMgoDeleteErr, err, database, collection, filter)
	}
	return int(result.DeletedCount), nil
}

// DeleteByObjId 根据mongo生成的ObjectId进行查找,等同于DeleteOne
func (c *MgoConn) DeleteByObjId(database, collection, oId string) (int, error) {
	_id, err := primitive.ObjectIDFromHex(oId)
	if err != nil {
		return 0, dberr.NewErr(ErrMgoObjectErr, oId, err)
	}

	filter := bson.D{
		{Key: "_id", Value: _id},
	}
	return c.DeleteOne(database, collection, filter)
}

// UpdateOne 单一更新
func (c *MgoConn) UpdateOne(database, collection string, filter, data any) error {
	if filter == nil {
		filter = bson.D{}
	}
	set := c.client.Database(database).Collection(collection)
	if _, err := set.UpdateOne(c.ctx, filter, bson.D{{Key: "$set", Value: data}}); err != nil {
		return dberr.NewErr(ErrMgoUpdateErr, err, database, collection, filter, data)
	}
	return nil
}

// UpdateByObjId 根据mongo生成的ObjectId进行更新,等同于UpdateOne
func (c *MgoConn) UpdateByObjId(database, collection, oId string, data any) error {
	_id, err := primitive.ObjectIDFromHex(oId)
	if err != nil {
		return dberr.NewErr(ErrMgoObjectErr, oId, err)
	}

	filter := bson.D{
		{Key: "_id", Value: _id},
	}
	return c.UpdateOne(database, collection, filter, data)
}

// UpdateN 批量更新
func (c *MgoConn) Update(database, collection string, filter, data any) error {
	if filter == nil {
		filter = bson.D{}
	}
	set := c.client.Database(database).Collection(collection)
	if _, err := set.UpdateMany(c.ctx, filter, bson.D{{Key: "$set", Value: data}}); err != nil {
		return dberr.NewErr(ErrMgoUpdateErr, err, database, collection, filter, data)
	}
	return nil
}

// ReplaceOne 整个文档内容替换(除了ObjectId)
func (c *MgoConn) ReplaceOne(database, collection string, filter, replacement any) error {
	if filter == nil {
		filter = bson.D{}
	}
	set := c.client.Database(database).Collection(collection)

	if _, err := set.ReplaceOne(c.ctx, filter, replacement); err != nil {
		return dberr.NewErr(ErrMgoReplaceErr, database, err, collection, filter, replacement)
	}
	return nil
}

// ReplaceByObjId 根据mongo生成的ObjectId进行更新,等同于ReplaceOne
func (c *MgoConn) ReplaceByObjId(database, collection string, oId string, replacement any) error {
	_id, err := primitive.ObjectIDFromHex(oId)
	if err != nil {
		return dberr.NewErr(ErrMgoObjectErr, oId, err)
	}

	filter := bson.D{
		{Key: "_id", Value: _id},
	}
	return c.ReplaceOne(database, collection, filter, replacement)
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
			if err := savectrl.SaveBox(func() {
				// todo cmd 检查?
				handler := c.opFuncs[opObj.cmd]
				opObj.resultAccept <- handler(opObj.args...)
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
