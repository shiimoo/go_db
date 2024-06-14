package mgo

import (
	"context"
	"fmt"

	"github.com/shiimoo/godb/dberr"
	"go.mongodb.org/mongo-driver/bson"
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
		return dberr.NewErr(dberr.ErrMgoInsertErr, database, collection, err)
	}
	return nil
}

// InsertN 批量插入数据
