package savectrl

import "github.com/shiimoo/godb/dberr"

var (
	// 转型获取数据库名错误
	ErrSaveBoxPanic = dberr.TempErr("SaveBox Err: %s")
)
