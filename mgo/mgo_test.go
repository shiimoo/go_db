package mgo

import (
	"context"
	"fmt"
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

func _getMgr(t *testing.T) {
	var err error
	ctx := context.Background()
	dbMgr, err = GetMgr(ctx, "default")
	if err != nil {
		t.Error(err)
	}
	dbMgr.Seturl("127.0.0.1", 27017)
	num, err := dbMgr.Connect(10)
	if err != nil {
		t.Error(err)
	}
	dbMgr.Start()
	log.Println("数据库连接成功!", num)
}

// TestCreateIndex 创建索引
func TestCreateIndex(t *testing.T) {
	_getMgr(t)
	// 对于已存在的集合不予创建索引,视实际业务需求场景确定
	ok, err := dbMgr.HasCollection(database, testCollection)
	if err != nil {
		t.Errorf("HasCollection Err: %s", err)
	}
	if !ok {
		if err := dbMgr.CreateIndex(database, testCollection, Indexs{
			{Key: "id", Value: 1},
		}); err != nil {
			t.Errorf("CreateIndex Err: %s", err)
		}
		t.Logf("%s.%s不存在, 重建索引!\n", database, testCollection)
	}
}

// TestInsertN 插入数据测试:多条数据
func TestInsertN(t *testing.T) {
	_getMgr(t)

	// 插入数据测试 多条数据
	testDatas := make([]any, 0)
	ti := time.Now().Unix()
	for i := 1; i < 10; i++ {
		id := ti*1000 + int64(i)
		testDatas = append(testDatas, map[string]any{"id": id, "name": "testname_" + fmt.Sprint(id)})
	}
	if err := dbMgr.InsertN(database, testCollection, testDatas); err != nil {
		t.Error(err)
	}
	t.Log("插入数据成功")
}

// TestInsertOne 插入数据测试:单条数据
func TestInsertOne(t *testing.T) {
	_getMgr(t)

	testData := map[string]any{"id": 1, "name": "testname"}
	if err := dbMgr.InsertOne(database, testCollection, testData); err != nil {
		t.Error(err)
	}
	t.Log("插入数据成功")
}

// Find 查询测试(带筛选)
func TestFind(t *testing.T) {
	_getMgr(t)

	dataList, err := dbMgr.Find(database, testCollection, bson.M{"name": "testname"}, 2)
	if err != nil {
		t.Error(err)
	}

	t.Log("数据库查询成功", len(dataList))
	for index, data := range dataList {
		t.Log(index, bson.Raw(data).String())
	}
}

// FindAll 查询全部
func TestFindAll(t *testing.T) {
	_getMgr(t)

	dataList, err := dbMgr.FindAll(database, testCollection)
	if err != nil {
		t.Error(err)
	}

	t.Log("数据库All查询成功", len(dataList))
	for index, data := range dataList {
		t.Log(index, bson.Raw(data).String())
	}
}

// TestFindOne 查询测试(单)
func TestFindOne(t *testing.T) {
	_getMgr(t)
	data, err := dbMgr.FindOne(database, testCollection, bson.M{"name": "testname"})
	if err != nil {
		t.Error(err)
	}
	t.Log("数据库查询成功", bson.Raw(data).String())
}

// TestFindByObjId 更新测试#单条数据使用ObjectId
func TestFindByObjId(t *testing.T) {
	_getMgr(t)

	data, err := dbMgr.FindByObjId(database, testCollection, "668e9c41237c4c425c66c4d8")
	if err != nil {
		t.Error(err)
	}
	t.Log("数据库查询成功", bson.Raw(data).String())
}

// TestDelete 删除测试#指定范围
func TestDelete(t *testing.T) {
	_getMgr(t)

	delCount, err := dbMgr.Delete(database, testCollection, bson.M{"name": "testname"})
	if err != nil {
		t.Error(err)
	}
	t.Log("数据库删除成功", delCount)
}

// TestDeleteOne 删除测试#单个范围
func TestDeleteOne(t *testing.T) {
	_getMgr(t)

	delCount, err := dbMgr.DeleteOne(database, testCollection, bson.M{"name": "testname_1720622145004"})
	if err != nil {
		t.Error(err)
	}
	t.Log("数据库删除成功", delCount)
}

// TestDeleteAll 删除测试#范围内全部删除
func TestDeleteAll(t *testing.T) {
	_getMgr(t)

	delCount, err := dbMgr.DeleteAll(database, testCollection, bson.M{"name": "testnameN"})
	if err != nil {
		t.Error(err)
	}
	t.Log("数据库删除成功", delCount)
}

// TestDeleteByObjId 更新测试#单条数据使用ObjectId
func TestDeleteByObjId(t *testing.T) {
	_getMgr(t)

	delCount, err := dbMgr.DeleteByObjId(database, testCollection, "668e9c41237c4c425c66c4d3")
	if err != nil {
		t.Error(err)
	}
	t.Log("数据库删除成功", delCount)
}

// TestUpdate 更新测试
func TestUpdate(t *testing.T) {
	_getMgr(t)

	if err := dbMgr.Update(
		database,
		testCollection,
		bson.M{"name": "testnameN"},
		map[string]any{"other_name_update": "TestUpdated"}); err != nil {
		t.Error(err)
	}
	t.Log("数据库更新成功")
}

// TestUpdateOne 更新测试#单条数据
func TestUpdateOne(t *testing.T) {
	_getMgr(t)

	if err := dbMgr.UpdateOne(
		database,
		testCollection,
		bson.M{"id": 1720622145005},
		map[string]any{"other_name": "testname----"}); err != nil {
		t.Error(err)
	}
	t.Log("数据库更新成功")
}

// TestUpdateByObjId 更新测试#单条数据使用ObjectId
func TestUpdateByObjId(t *testing.T) {
	_getMgr(t)

	if err := dbMgr.UpdateByObjId(
		database,
		testCollection,
		"668e9c41237c4c425c66c4d5",
		map[string]any{"other_name": "shimo111"}); err != nil {
		t.Error(err)
	}
	t.Log("数据库更新成功")
}

// TestReplaceOne 文档替换
func TestReplaceOne(t *testing.T) {
	_getMgr(t)

	err := dbMgr.ReplaceOne(
		database,
		testCollection,
		bson.M{"other_name": "shimo----TestReplaceOne"},
		map[string]any{"other_name": "shimo----TestReplaceOn1111e"},
	)
	if err != nil {
		t.Error(err)
	}
	t.Log("数据库替换成功")
}

// TestReplaceByObjId 文档替换
func TestReplaceByObjId(t *testing.T) {
	_getMgr(t)

	err := dbMgr.ReplaceByObjId(
		database,
		testCollection,
		"668e9c41237c4c425c66c4d6",
		map[string]any{
			"id":         20240710000,
			"other_name": "shimo----TestReplaceByObjId",
		},
	)
	if err != nil {
		t.Error(err)
	}
	t.Log("数据库替换成功")
}
