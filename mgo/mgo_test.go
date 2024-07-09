package mgo

import (
	"context"
	"log"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	database       = "testdb"
	testCollection = "testSet"
)

var dbMgr *mgr

func _getMgr() {
	var err error
	ctx := context.Background()
	dbMgr, err = GetMgr(ctx, "default")
	if err != nil {
		log.Fatalln(err)
	}
	dbMgr.Seturl("127.0.0.1", 27017)
	num, err := dbMgr.Connect(10)
	if err != nil {
		log.Fatalln(err)
	}
	dbMgr.Start()
	log.Println("数据库连接成功!", num)
}

// TestCreateIndex 创建索引
func TestCreateIndex(t *testing.T) {
	_getMgr()
	// 对于已存在的集合不予创建索引,视实际业务需求场景确定
	if !dbMgr.HasCollection(database, testCollection) {
		dbMgr.CreateIndex(database, testCollection, Indexs{
			{Key: "id", Value: 1},
		})
		log.Printf("%s.%s不存在, 重建索引!\n", database, testCollection)
	}
}

// TestInsertOne 插入数据测试:单条数据
func TestInsertOne(t *testing.T) {
	_getMgr()

	testData := map[string]any{"id": 1, "name": "testname"}
	if err := dbMgr.InsertOne(database, testCollection, testData); err != nil {
		log.Fatalln(err)
	}
	log.Println("插入数据成功")
}

// TestInsertN 插入数据测试:多条数据
func TestInsertN(t *testing.T) {
	_getMgr()

	// 插入数据测试 多条数据
	testDatas := make([]any, 0)
	ti := time.Now().Unix()
	for i := 1; i < 10; i++ {
		testDatas = append(testDatas, map[string]any{"id": ti*1000 + int64(i), "name": "testname"})
	}
	if err := dbMgr.InsertN(database, testCollection, testDatas); err != nil {
		log.Fatalln(err)
	}
	log.Println("插入数据成功")
}

// Find 查询测试(带筛选)
func TestFind(t *testing.T) {
	_getMgr()

	dataList, err := dbMgr.Find(database, testCollection, bson.M{"name": "testname"}, 2)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("数据库查询成功", len(dataList))
	for index, data := range dataList {
		log.Println(index, bson.Raw(data).String())
	}
}

// FindAll 查询全部
func TestFindAll(t *testing.T) {
	_getMgr()

	dataList, err := dbMgr.FindAll(database, testCollection)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("数据库All查询成功", len(dataList))
	for index, data := range dataList {
		log.Println(index, bson.Raw(data).String())
	}
}

// TestFindOne 查询测试(单)
func TestFindOne(t *testing.T) {
	_getMgr()
	data, err := dbMgr.FindOne(database, testCollection, bson.M{"name": "testname"})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("数据库查询成功", bson.Raw(data).String())
}

// TestFindByObjId 更新测试#单条数据使用ObjectId
func TestFindByObjId(t *testing.T) {
	_getMgr()

	data, err := dbMgr.FindByObjId(database, testCollection, "666beb15926cb7bca675d6f0")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("数据库查询成功", bson.Raw(data).String())
}

// TestDelete 删除测试#指定范围
func TestDelete(t *testing.T) {
	_getMgr()

	delCount, err := dbMgr.Delete(database, testCollection, bson.M{"name": "testname"})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("数据库删除成功", delCount)
}

// TestDeleteOne 删除测试#单个范围
func TestDeleteOne(t *testing.T) {
	_getMgr()

	delCount, err := dbMgr.DeleteOne(database, testCollection, bson.M{"name": "testname"})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("数据库删除成功", delCount)
}

// TestDeleteAll 删除测试#范围内全部删除
func TestDeleteAll(t *testing.T) {
	_getMgr()

	delCount, err := dbMgr.DeleteAll(database, testCollection, bson.M{"name": "testname"})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("数据库删除成功", delCount)
}

// TestDeleteByObjId 更新测试#单条数据使用ObjectId
func TestDeleteByObjId(t *testing.T) {
	_getMgr()

	delCount, err := dbMgr.DeleteByObjId(database, testCollection, "666beb15926cb7bca675d6f0")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("数据库删除成功", delCount)
}

// TestUpdateOne 更新测试#单条数据
func TestUpdateOne(t *testing.T) {
	_getMgr()

	if err := dbMgr.Update(
		database,
		testCollection,
		bson.M{"name": "testname"},
		map[string]any{"other_name": "testname----"}); err != nil {
		log.Fatalln(err)
	}
	log.Println("数据库更新成功")
}

// TestUpdateByObjId 更新测试#单条数据使用ObjectId
func TestUpdateByObjId(t *testing.T) {
	_getMgr()

	if err := dbMgr.UpdateByObjId(
		database,
		testCollection,
		"666beb15926cb7bca675d6f0",
		map[string]any{"other_name": "shimo111"}); err != nil {
		log.Fatalln(err)
	}
	log.Println("数据库更新成功")
}

// TestUpdateByObjId 更新测试#单条数据
func TestUpdated(t *testing.T) {
	_getMgr()

	if err := dbMgr.Update(
		database,
		testCollection,
		bson.M{"name": "testname"},
		map[string]any{"other_name_update": "TestUpdated"}); err != nil {
		log.Fatalln(err)
	}
	log.Println("数据库更新成功")
}

// TestReplaceOne 文档替换
func TestReplaceOne(t *testing.T) {
	_getMgr()

	err := dbMgr.ReplaceOne(database, testCollection, bson.M{"name": "shimo"}, map[string]any{"other_name": "shimo----shimo"})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("数据库替换成功")
}

// TestReplaceByObjId 文档替换
func TestReplaceByObjId(t *testing.T) {
	_getMgr()

	err := dbMgr.ReplaceOne(database, testCollection, "666beb15926cb7bca675d6f0", map[string]any{"other_name": "shimo----shimo"})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("数据库替换成功")
}
