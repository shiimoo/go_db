package dberr

// ErrMgo%s
var (

	// 链接管理器key异常
	ErrMgoMgrKey = TempErr("mgoMgr key value Err : %s")
	// 请求创建的链接数量异常
	ErrMgoConnNum = TempErr("mgoMgr create conn num Err : %s")
	// 创建的链接异常
	ErrMgoConnectErr = TempErr("mgoMgr create conn Err : %s")
	// mongo 插入数据错误
	ErrMgoInsertErr = TempErr("mgoMgr insert data to database[%s], collection[%s], err: %s")
)
