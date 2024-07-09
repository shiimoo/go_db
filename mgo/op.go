package mgo

import (
	"fmt"
	"reflect"

	"github.com/shiimoo/godb/dberr"
	"go.mongodb.org/mongo-driver/bson"
)

// ***** 操作码-枚举 *****

const (
	cmdHasCollection = "hasCollection" // 判定数据库database中是否存在集合collectio
	cmdCreateIndex   = "createIndex"   // 创建索引

	cmdInsertN   = "insertN"   // 批量插入数据
	cmdInsertOne = "insertOne" // 插入单挑数据

	cmdFind        = "find"        // 加载数据
	cmdFindAll     = "findAll"     // 加载全部数据
	cmdFindOne     = "findOne"     // 加载单个数据
	cmdFindByObjId = "findByObjId" // 根据mongo生成的ObjectId进行查找,等同于findOne

	cmdDelete        = "delete"        // 删除数据
	cmdDeleteAll     = "deleteAll"     // 全部删除(清空)
	cmdDeleteOne     = "deleteOne"     // 删除1个
	cmdDeleteByObjId = "deleteByObjId" // 根据mongo生成的ObjectId进行查找,等同于deleteOne

	cmdUpdate        = "update"        // 批量更新
	cmdUpdateOne     = "updateOne"     // 单一更新
	cmdUpdateByObjId = "updateByObjId" // 根据mongo生成的ObjectId进行更新,等同于updateOne

	cmdReplaceOne     = "replaceOne"     // 整个文档内容替换(除了ObjectId)
	cmdReplaceByObjId = "replaceByObjId" // 根据mongo生成的ObjectId进行更新,等同于replaceOne
)

// ***** 操作码-统一调度方法定义 *****

type opFunc func(otherParams ...any) *opResult

// ***** 操作码-定制输入参数检查并解析(抛错误) *****

// cmdHasCollection: [] // no params

// cmdCreateIndex: [Indexs, ] // 1 params
func parseCreateIndexParams(params ...any) Indexs {
	if len(params) < 1 {
		panic(dberr.NewErr(
			ErrOpParmsErr,
			"params length < 1",
			cmdCreateIndex,
			params,
		))
	}
	oriIndexs := params[0]
	indexs, ok := oriIndexs.(Indexs)
	if !ok {
		panic(dberr.NewErr(
			ErrOpParmsErr,
			fmt.Sprintf("params[0] type must be mgo.Indexs, but is %s", reflect.TypeOf(oriIndexs)),
			cmdCreateIndex,
			params,
		))
	}
	return indexs
}

// cmdInsertN: [[]any, ] // 1 params
func parseInsertNParams(params ...any) []any {
	if len(params) < 1 {
		panic(dberr.NewErr(
			ErrOpParmsErr,
			"params length < 1",
			cmdInsertN,
			params,
		))
	}
	oriDatas := params[0]
	datas, ok := oriDatas.([]any)
	if !ok {
		panic(dberr.NewErr(
			ErrOpParmsErr,
			fmt.Sprintf("params[0] type must be []any, but is %s", reflect.TypeOf(oriDatas)),
			cmdInsertN,
			params,
		))
	}
	return datas
}

// cmdInsertOne: [any, ] // 1 params
func parseInsertOneParams(params ...any) any {
	if len(params) < 1 {
		panic(dberr.NewErr(
			ErrOpParmsErr,
			"params length < 1",
			cmdInsertOne,
			params,
		))
	}
	return params[0]
}

// cmdFind: [bson.D, int64, ] // 2 params
func parseFindParams(params ...any) (bson.D, int64) {
	var filterOri, numOri any
	if len(params) > 0 {
		filterOri = params[0]
	}
	if len(params) > 1 {
		numOri = params[1]
	}

	var filter bson.D
	if filterOri == nil {
		filter = bson.D{}
	} else {
		var ok bool
		filter, ok = filterOri.(bson.D)
		if !ok {
			panic(dberr.NewErr(
				ErrOpParmsErr,
				fmt.Sprintf("filter type must be bson.D, but is %s", reflect.TypeOf(filterOri)),
				cmdFind,
				params,
			))
		}
	}

	var num int64
	if numOri == nil {
		num = 0
	} else {
		var ok bool
		num, ok = numOri.(int64)
		if !ok {
			panic(dberr.NewErr(
				ErrOpParmsErr,
				fmt.Sprintf("num type must be int64, but is %s", reflect.TypeOf(numOri)),
				cmdFind,
				params,
			))
		}
	}
	return filter, num
}

// cmdFindAll: [] // no params

// cmdFindOne: [bson.D, ] // 1 params
func parseFindOneParams(params ...any) bson.D {
	var filterOri any
	if len(params) > 0 {
		filterOri = params[0]
	}

	var filter bson.D
	if filterOri == nil {
		filter = bson.D{}
	} else {
		var ok bool
		filter, ok = filterOri.(bson.D)
		if !ok {
			panic(dberr.NewErr(
				ErrOpParmsErr,
				fmt.Sprintf("filter type must be bson.D, but is %s", reflect.TypeOf(filterOri)),
				cmdFind,
				params,
			))
		}
	}

	return filter
}

// cmdFindByObjId: [string, ] // 1 params
func parseFindByObjIdParams(params ...any) string {
	if len(params) < 1 {
		panic(dberr.NewErr(
			ErrOpParmsErr,
			"params length < 1",
			cmdFindByObjId,
			params,
		))
	}

	oriOId := params[0]
	oId, ok := oriOId.(string)
	if !ok {
		panic(dberr.NewErr(
			ErrOpParmsErr,
			fmt.Sprintf("params[0] type must be string, but is %s", reflect.TypeOf(oriOId)),
			cmdFindByObjId,
			params,
		))
	}

	return oId
}

// cmdDelete: [bson.D, ] // 1 params
func parseDeleteParams(params ...any) bson.D {
	var filterOri any
	if len(params) > 0 {
		filterOri = params[0]
	}

	var filter bson.D
	if filterOri == nil {
		filter = bson.D{}
	} else {
		var ok bool
		filter, ok = filterOri.(bson.D)
		if !ok {
			panic(dberr.NewErr(
				ErrOpParmsErr,
				fmt.Sprintf("filter type must be bson.D, but is %s", reflect.TypeOf(filterOri)),
				cmdDelete,
				params,
			))
		}
	}

	return filter
}

// cmdDeleteAll: [] // no params

// cmdDeleteOne: [bson.D, ] // 1 params
func parseDeleteOneParams(params ...any) bson.D {
	var filterOri any
	if len(params) > 0 {
		filterOri = params[0]
	}

	var filter bson.D
	if filterOri == nil {
		filter = bson.D{}
	} else {
		var ok bool
		filter, ok = filterOri.(bson.D)
		if !ok {
			panic(dberr.NewErr(
				ErrOpParmsErr,
				fmt.Sprintf("filter type must be bson.D, but is %s", reflect.TypeOf(filterOri)),
				cmdDeleteOne,
				params,
			))
		}
	}

	return filter
}

// cmdDeleteByObjId: [string, ] // 1 params
func parseDeleteByObjIdParams(params ...any) string {
	if len(params) < 1 {
		panic(dberr.NewErr(
			ErrOpParmsErr,
			"params length < 1",
			cmdDeleteByObjId,
			params,
		))
	}

	oriOId := params[0]
	oId, ok := oriOId.(string)
	if !ok {
		panic(dberr.NewErr(
			ErrOpParmsErr,
			fmt.Sprintf("params[0] type must be string, but is %s", reflect.TypeOf(oriOId)),
			cmdDeleteByObjId,
			params,
		))
	}
	return oId
}

// cmdUpdate: [bson.D, any, ] // 2 params
func parseUpdateParams(params ...any) (bson.D, any) {
	if len(params) < 2 {
		panic(dberr.NewErr(
			ErrOpParmsErr,
			"params length < 2",
			cmdUpdate,
			params,
		))
	}
	filterOri := params[0]
	var filter bson.D
	if filterOri == nil {
		filter = bson.D{}
	} else {
		var ok bool
		filter, ok = filterOri.(bson.D)
		if !ok {
			panic(dberr.NewErr(
				ErrOpParmsErr,
				fmt.Sprintf("filter type must be bson.D, but is %s", reflect.TypeOf(filterOri)),
				cmdUpdate,
				params,
			))
		}
	}
	data := params[0]
	if data == nil {
		panic(dberr.NewErr(
			ErrOpParmsErr,
			"update data is nil",
			cmdUpdate,
			params,
		))
	}

	return filter, data
}

// cmdUpdateOne: [bson.D, any, ] // 2 params
func parseUpdateOneParams(params ...any) (bson.D, any) {
	if len(params) < 2 {
		panic(dberr.NewErr(
			ErrOpParmsErr,
			"params length < 2",
			cmdUpdateOne,
			params,
		))
	}
	filterOri := params[0]
	var filter bson.D
	if filterOri == nil {
		filter = bson.D{}
	} else {
		var ok bool
		filter, ok = filterOri.(bson.D)
		if !ok {
			panic(dberr.NewErr(
				ErrOpParmsErr,
				fmt.Sprintf("filter type must be bson.D, but is %s", reflect.TypeOf(filterOri)),
				cmdUpdateOne,
				params,
			))
		}
	}
	data := params[0]
	if data == nil {
		panic(dberr.NewErr(
			ErrOpParmsErr,
			"update data is nil",
			cmdUpdate,
			params,
		))
	}

	return filter, data
}

// cmdUpdateByObjId: [string, ] // 1 params
func parseUpdateByObjIdParams(params ...any) string {
	if len(params) < 1 {
		panic(dberr.NewErr(
			ErrOpParmsErr,
			"params length < 1",
			cmdUpdateByObjId,
			params,
		))
	}

	oriOId := params[0]
	oId, ok := oriOId.(string)
	if !ok {
		panic(dberr.NewErr(
			ErrOpParmsErr,
			fmt.Sprintf("params[0] type must be string, but is %s", reflect.TypeOf(oriOId)),
			cmdUpdateByObjId,
			params,
		))
	}
	return oId
}

// cmdReplaceOne: [bson.D, any, ] // 2 params
func parseReplaceOneParams(params ...any) (bson.D, any) {
	if len(params) < 2 {
		panic(dberr.NewErr(
			ErrOpParmsErr,
			"params length < 2",
			cmdReplaceOne,
			params,
		))
	}
	filterOri := params[0]
	var filter bson.D
	if filterOri == nil {
		filter = bson.D{}
	} else {
		var ok bool
		filter, ok = filterOri.(bson.D)
		if !ok {
			panic(dberr.NewErr(
				ErrOpParmsErr,
				fmt.Sprintf("filter type must be bson.D, but is %s", reflect.TypeOf(filterOri)),
				cmdReplaceOne,
				params,
			))
		}
	}
	replacement := params[0]
	if replacement == nil {
		panic(dberr.NewErr(
			ErrOpParmsErr,
			"replacement is nil",
			cmdReplaceOne,
			params,
		))
	}

	return filter, replacement
}

// cmdReplaceByObjId: [bson.D, any, ] // 2 params
func parseReplaceByObjIdParams(params ...any) string {
	if len(params) < 1 {
		panic(dberr.NewErr(
			ErrOpParmsErr,
			"params length < 1",
			cmdReplaceByObjId,
			params,
		))
	}

	oriOId := params[0]
	oId, ok := oriOId.(string)
	if !ok {
		panic(dberr.NewErr(
			ErrOpParmsErr,
			fmt.Sprintf("params[0] type must be string, but is %s", reflect.TypeOf(oriOId)),
			cmdReplaceByObjId,
			params,
		))
	}
	return oId
}

// ***** 操作码-结果参数解析(返回错误) *****

// cmdHasCollection: [bool, ] // 1 params
func parseResultHasCollection(results ...any) (bool, error) {
	if len(results) < 1 {
		return false, dberr.NewErr(
			ErrOpResultPraseErr,
			"results length < 1",
			cmdHasCollection,
			results,
		)
	}
	oriData := results[0]
	data, ok := oriData.(bool)
	if !ok {
		return false, dberr.NewErr(
			ErrOpResultPraseErr,
			fmt.Sprintf("result[0] type must be bool, but is %s", reflect.TypeOf(oriData)),
			cmdHasCollection,
			results,
		)
	}
	return data, nil
}

// cmdCreateIndex: [] // no params

// cmdInsertN: [] // no params

// cmdInsertOne: [] // no params

// cmdFind: [[][]byte, ] // 1 params
func parseResultFind(results ...any) ([][]byte, error) {
	if len(results) < 1 {
		return nil, dberr.NewErr(
			ErrOpResultPraseErr,
			"results length < 1",
			cmdFind,
			results,
		)
	}
	oriDatas := results[0]
	datas, ok := oriDatas.([][]byte)
	if !ok {
		return nil, dberr.NewErr(
			ErrOpResultPraseErr,
			fmt.Sprintf("result[0] type must be [][]byte, but is %s", reflect.TypeOf(oriDatas)),
			cmdFind,
			results,
		)
	}
	return datas, nil
}

// cmdFindAll: [[][]byte, ] // 1 params
func parseResultFindAll(results ...any) ([][]byte, error) {
	if len(results) < 1 {
		return nil, dberr.NewErr(
			ErrOpResultPraseErr,
			"results length < 1",
			cmdFindAll,
			results,
		)
	}
	oriDatas := results[0]
	datas, ok := oriDatas.([][]byte)
	if !ok {
		return nil, dberr.NewErr(
			ErrOpResultPraseErr,
			fmt.Sprintf("result[0] type must be [][]byte, but is %s", reflect.TypeOf(oriDatas)),
			cmdFindAll,
			results,
		)
	}
	return datas, nil
}

// cmdFindOne: [[]byte, ] // 1 params
func parseResultFindOne(results ...any) ([]byte, error) {
	if len(results) < 1 {
		return nil, dberr.NewErr(
			ErrOpResultPraseErr,
			"results length < 1",
			cmdFindOne,
			results,
		)
	}
	oriData := results[0]
	datas, ok := oriData.([]byte)
	if !ok {
		return nil, dberr.NewErr(
			ErrOpResultPraseErr,
			fmt.Sprintf("result[0] type must be []byte, but is %s", reflect.TypeOf(oriData)),
			cmdFindOne,
			results,
		)
	}
	return datas, nil
}

// cmdFindByObjId: [[]byte, ] // 1 params
func parseResultFindByObjId(results ...any) ([]byte, error) {
	if len(results) < 1 {
		return nil, dberr.NewErr(
			ErrOpResultPraseErr,
			"results length < 1",
			cmdFindByObjId,
			results,
		)
	}
	oriData := results[0]
	datas, ok := oriData.([]byte)
	if !ok {
		return nil, dberr.NewErr(
			ErrOpResultPraseErr,
			fmt.Sprintf("result[0] type must be []byte, but is %s", reflect.TypeOf(oriData)),
			cmdFindByObjId,
			results,
		)
	}
	return datas, nil
}

// cmdDelete: [int, ] // 1 params
func parseResultDelete(results ...any) (int, error) {
	if len(results) < 1 {
		return 0, dberr.NewErr(
			ErrOpResultPraseErr,
			"results length < 1",
			cmdDelete,
			results,
		)
	}
	oriNum := results[0]
	num, ok := oriNum.(int)
	if !ok {
		return 0, dberr.NewErr(
			ErrOpResultPraseErr,
			fmt.Sprintf("result[0] type must be int, but is %s", reflect.TypeOf(oriNum)),
			cmdDelete,
			results,
		)
	}
	return num, nil
}

// cmdDeleteAll: [int, ] // 1 params
func parseResultDeleteAll(results ...any) (int, error) {
	if len(results) < 1 {
		return 0, dberr.NewErr(
			ErrOpResultPraseErr,
			"results length < 1",
			cmdDeleteAll,
			results,
		)
	}
	oriNum := results[0]
	num, ok := oriNum.(int)
	if !ok {
		return 0, dberr.NewErr(
			ErrOpResultPraseErr,
			fmt.Sprintf("result[0] type must be int, but is %s", reflect.TypeOf(oriNum)),
			cmdDeleteAll,
			results,
		)
	}
	return num, nil
}

// cmdDeleteOne: [int, ] // 1 params
func parseResultDeleteOne(results ...any) (int, error) {
	if len(results) < 1 {
		return 0, dberr.NewErr(
			ErrOpResultPraseErr,
			"results length < 1",
			cmdDeleteOne,
			results,
		)
	}
	oriNum := results[0]
	num, ok := oriNum.(int)
	if !ok {
		return 0, dberr.NewErr(
			ErrOpResultPraseErr,
			fmt.Sprintf("result[0] type must be int, but is %s", reflect.TypeOf(oriNum)),
			cmdDeleteOne,
			results,
		)
	}
	return num, nil
}

// cmdDeleteByObjId: [int, ] // 1 params
func parseResultDeleteByObjId(results ...any) (int, error) {
	if len(results) < 1 {
		return 0, dberr.NewErr(
			ErrOpResultPraseErr,
			"results length < 1",
			cmdDeleteByObjId,
			results,
		)
	}
	oriNum := results[0]
	num, ok := oriNum.(int)
	if !ok {
		return 0, dberr.NewErr(
			ErrOpResultPraseErr,
			fmt.Sprintf("result[0] type must be int, but is %s", reflect.TypeOf(oriNum)),
			cmdDeleteByObjId,
			results,
		)
	}
	return num, nil
}

// cmdUpdate: [] // no params

// cmdUpdateOne: [] // no params

// cmdUpdateByObjId: [] // no params

// cmdReplaceOne: [] // no params

// cmdReplaceByObjId: [] // no params
