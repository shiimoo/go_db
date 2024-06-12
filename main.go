package main

import (
	"context"
	"log"

	"github.com/shiimoo/godb/mgo"
)

var (
	database       = "testdb"
	testCollection = "testSet"
)

func main() {
	ctx := context.Background()
	if err := mgo.Connect(ctx, mgo.NewConnCfg("127.0.0.1", 27017)); err != nil {
		log.Fatalln(err)
	}
	log.Println("数据库连接成功!")

	if !mgo.HasCollection(database, testCollection) {
		log.Printf("判定数据库[%s]中的集合[%s]不存在!", database, testCollection)
		mgo.CreateIndex(database, testCollection, mgo.MgoIndex{
			{"id", 1},
		})
	}
}
