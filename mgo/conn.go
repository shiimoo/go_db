package mgo

import (
	"context"
	"fmt"

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

type MgoConn struct {
	client *mongo.Client
}

func Connect(ctx context.Context, cfg ConnCfg) error {
	clientOptions := options.Client().ApplyURI(cfg.GenUrl())
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
	conn.client = client
	connMgr.AddConn(conn)
	return nil
}
