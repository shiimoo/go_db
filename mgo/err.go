package mgo

import "github.com/shiimoo/godb/dberr"

// ErrMgo%s
var (
	/***** 参数转型 *****/

	// 转型获取数据库名错误
	ErrToOpErr         = dberr.TempErr("mgoMgr toOp value Err : %s")
	ErrToDatabaseErr   = dberr.TempErr("mgoMgr toDatabase value Err : %s")
	ErrToCollectionErr = dberr.TempErr("mgoMgr toCollection value Err : %s")
	ErrParamsErr       = dberr.TempErr("mgoMgr params Err : %s")

	/***** 业务操作 *****/

	// 链接管理器key异常
	ErrMgoMgrKey = dberr.TempErr("mgoMgr key value Err : %s")
	// 请求创建的链接数量异常
	ErrMgoConnNum = dberr.TempErr("mgoMgr create conn num Err : %s")
	// 创建的链接异常
	ErrMgoConnectErr = dberr.TempErr("mgoMgr create conn Err : %s")
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
