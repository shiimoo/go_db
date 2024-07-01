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

	// 查询测试(过滤)
	// dataList, err := dbMgr.Find(database, testCollection, bson.M{"name": "testname"}, 2)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// 查询测试(all)
	// dataList, err := dbMgr.FindAll(database, testCollection)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// 查询测试(单)

	// data, err := dbMgr.FindOne(database, testCollection, bson.M{"id": 1})
	// data, err := dbMgr.FindOne(database, testCollection, bson.M{"id": 1})
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// log.Println("数据库查询成功", bson.Raw(data).String())

	// 删除测试
	// delCount, err := dbMgr.DeleteOne(database, testCollection, nil) // bson.M{"id": 1})
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// log.Println("数据库删除成功", delCount)

	// 更新测试
	// err = dbMgr.UpdateOne(database, testCollection, nil, map[string]any{"name": "shimo111"})
	// err = dbMgr.UpdateByObjId(database, testCollection, "666beb15926cb7bca675d6f0", map[string]any{"name": "shimo111"})
	// err = dbMgr.Update(database, testCollection, bson.M{"name": "testname"}, map[string]any{"other_name": "testname----"})
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// log.Println("数据库更新成功")

}
