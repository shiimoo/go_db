package mgo

import "github.com/shiimoo/godb/dberr"

// ErrMgo%s
var (
	// mongo链接关闭
	ErrConnIsClose = dberr.TempErr("mgoMgr conn close")

	/***** 参数转型 *****/

	// 操作码 错误
	ErrToOpErr = dberr.TempErr("mgoMgr toOp value err : %s")
	// 库名 错误
	ErrToDatabaseErr = dberr.TempErr("mgoMgr toDatabase value err : %s")
	// 集合名 错误
	ErrToCollectionErr = dberr.TempErr("mgoMgr toCollection value err : %s")
	// 原始参数解析 错误
	ErrParamsErr = dberr.TempErr("mgoMgr params err : %s")
	// 操作码专用的参数解析 错误
	ErrOpParmsErr = dberr.TempErr("mgoMgr op params parse err : %s; details: op[%s], params[%v]")
	// 操作码结果解析 错误
	ErrOpResultPraseErr = dberr.TempErr("mgoMgr op results parse err : %s; details: op[%s], results[%v]")

	/***** 业务操作 *****/

	// 链接管理器key异常
	ErrMgoMgrKey = dberr.TempErr("mgoMgr key value err : %s")
	// 请求创建的链接数量异常
	ErrMgoConnNum = dberr.TempErr("mgoMgr create conn num err : %s")
	// 创建的链接异常
	ErrMgoConnectErr = dberr.TempErr("mgoMgr create conn err : %s")
	// 创建索引错误
	ErrMgoCreateIndexErr = dberr.TempErr("mgoMgr createIndex value err : %s; details: database[%s], collection[%s], indexs[%v]")
	// mongo 插入数据错误
	ErrMgoInsertErr = dberr.TempErr("mgoMgr insert err: %s; details: database[%s], collection[%s], data[%v]")
	// mongo 查询错误
	ErrMgoFindErr = dberr.TempErr("mgoMgr find err: %s; details: database[%s], collection[%s], filter[%v], num[%d]")
	// mongo 更新错误
	ErrMgoUpdateErr = dberr.TempErr("mgoMgr update err: %s; details: database[%s], collection[%s], filter[%v], update[%v]")
	// mongo 删除错误
	ErrMgoDeleteErr = dberr.TempErr("mgoMgr insert err: %s; details: database[%s], collection[%s], filter[%v]")
	// mongo 替换错误
	ErrMgoReplaceErr = dberr.TempErr("mgoMgr replace err: %s; details: database[%s], collection[%s], filter[%v], replacement[%v]")
	// mongo ObjectId错误
	ErrMgoObjectErr = dberr.TempErr("ObjectId Hex String[%s] Invalid, err: %s")
)
