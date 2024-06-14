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

	// 插入数据测试 单条数据
	// testData := map[string]any{"id": 1, "name": "testname"}
	// if err := dbMgr.InsertOne(database, testCollection, testData); err != nil {
	// 	log.Fatalln(err)
	// }

	// 插入数据测试 多条数据
	// testDatas := make([]any, 0)
	// ti := time.Now().Unix()
	// for i := 1; i < 10; i++ {
	// 	testDatas = append(testDatas, map[string]any{"id": ti*1000 + int64(i), "name": "testname"})
	// }
	// if err := dbMgr.InsertN(database, testCollection, testDatas); err != nil {
	// 	log.Fatalln(err)
	// }

	// 查询测试(all)
	if err := dbMgr.FindAll(database, testCollection); err != nil {
		log.Fatalln(1, err)
	}
	// 查询测试(单/过滤)
}
