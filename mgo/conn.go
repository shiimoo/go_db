package mgo

import (
	"context"
	"fmt"

	"github.com/shiimoo/godb/dberr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewConn(parent context.Context, url string) (*MgoConn, error) {
	cOpts := options.Client().ApplyURI(url)
	// ctx, cancel TODO:关闭监听待实现
	ctx, _ := context.WithCancel(parent)
	// Connect to MongoDB
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
	conn.client = client
	return conn, nil
}

type MgoConn struct {
	ctx    context.Context
	client *mongo.Client
}

// 判定数据库database中是否存在集合collection
func (c *MgoConn) hasCollection(database, collection string) bool {
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

// Indexs mongo 索引结构重写(复刻结构primitive.E) {{字段名:1/-1}, {关键字:值}}
type Indexs bson.D

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
