package mgo

import "github.com/shiimoo/godb/lib/base/errors"

// ErrMgo%s
var (
	// mongo链接关闭
	ErrConnIsClose = errors.TempErr("mgoMgr conn close")

	/***** 参数转型 *****/

	// 操作码 错误
	ErrToOpErr = errors.TempErr("mgoMgr toOp value err : %s")
	// 库名 错误
	ErrToDatabaseErr = errors.TempErr("mgoMgr toDatabase value err : %s")
	// 集合名 错误
	ErrToCollectionErr = errors.TempErr("mgoMgr toCollection value err : %s")
	// 原始参数解析 错误
	ErrParamsErr = errors.TempErr("mgoMgr params err : %s")
	// 操作码专用的参数解析 错误
	ErrOpParmsErr = errors.TempErr("mgoMgr op params parse err : %s; details: op[%s], params[%v]")
	// 操作码结果解析 错误
	ErrOpResultPraseErr = errors.TempErr("mgoMgr op results parse err : %s; details: op[%s], results[%v]")

	/***** 业务操作 *****/

	// 链接管理器key异常
	ErrMgoMgrKey = errors.TempErr("mgoMgr key value err : %s")
	// 请求创建的链接数量异常
	ErrMgoConnNum = errors.TempErr("mgoMgr create conn num err : %s")
	// 创建的链接异常
	ErrMgoConnectErr = errors.TempErr("mgoMgr create conn err : %s")
	// 创建索引错误
	ErrMgoCreateIndexErr = errors.TempErr("mgoMgr createIndex value err : %s; details: database[%s], collection[%s], indexs[%v]")
	// mongo 插入数据错误
	ErrMgoInsertErr = errors.TempErr("mgoMgr insert err: %s; details: database[%s], collection[%s], data[%v]")
	// mongo 查询错误
	ErrMgoFindErr = errors.TempErr("mgoMgr find err: %s; details: database[%s], collection[%s], filter[%v], num[%d]")
	// mongo 更新错误
	ErrMgoUpdateErr = errors.TempErr("mgoMgr update err: %s; details: database[%s], collection[%s], filter[%v], update[%v]")
	// mongo 删除错误
	ErrMgoDeleteErr = errors.TempErr("mgoMgr insert err: %s; details: database[%s], collection[%s], filter[%v]")
	// mongo 替换错误
	ErrMgoReplaceErr = errors.TempErr("mgoMgr replace err: %s; details: database[%s], collection[%s], filter[%v], replacement[%v]")
	// mongo ObjectId错误
	ErrMgoObjectErr = errors.TempErr("ObjectId Hex String[%s] Invalid, err: %s")
)
