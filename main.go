package main

import (
	"context"
	"log"

	"github.com/shiimoo/godb/mgo"
)

var ()

func main() {
	mgoText()

	// fmt.Println()

}

func mgoText() {
	database := "testdb"
	testCollection := "testSet"

	ctx := context.Background()
	dbMgr, err := mgo.GetMgr(ctx, "default")
	if err != nil {
		log.Fatalln(err)
	}
	dbMgr.Seturl("127.0.0.1", 27017)
	num, err := dbMgr.Connect(10)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("数据库连接成功!", num)
	// 查询 & 建立索引测试
	if !dbMgr.HasCollection(database, testCollection) {
		dbMgr.CreateIndex(database, testCollection, mgo.Indexs{
			{Key: "id", Value: 1},
		})
		log.Printf("%s.%s不存在, 重建索引!\n", database, testCollection)
	}

	// 插入数据测试
	testData := map[string]any{"id": 1, "name": "testname"}
	if err := dbMgr.InsertOne(database, testCollection, testData); err != nil {
		log.Fatalln(err)
	}
}
