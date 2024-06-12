package mgo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConnCfg struct {
	// 数据库访问地址
	host string
	// 数据库端口
	port int
}

func (cfg ConnCfg) GenUrl() string {
	return fmt.Sprintf("mongodb://%s:%d", cfg.host, cfg.port)
}

func NewConnCfg(host string, port int) ConnCfg {
	if host == "" {
		host = "localhost"
	}
	if port <= 0 {
		port = 27017
	}
	return ConnCfg{
		host: host,
		port: port,
	}
}

func Connect(parent context.Context, cfg ConnCfg) error {
	clientOptions := options.Client().ApplyURI(cfg.GenUrl())
	// ctx, cancel TODO:关闭监听待实现
	ctx, _ := context.WithCancel(parent)
	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}
	conn := new(MgoConn)
	conn.ctx = ctx
	conn.client = client
	getMgr().AddConn(conn)
	return nil
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

// MgoIndex mongo 索引 {{字段名:1/-1}, {关键字:值}}
type MgoIndex bson.D // todo 自定义的索引结构可以优化刻度性

// 创建索引
func (c *MgoConn) CreateIndex(database, collection string, index MgoIndex) {
	set := c.client.Database(database).Collection(collection)
	fmt.Println(set.Indexes().CreateOne(c.ctx, mongo.IndexModel{
		Keys: index,
	}))
}
