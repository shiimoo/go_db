package savectrl

import "github.com/shiimoo/godb/lib/base/errors"

var (
	// 转型获取数据库名错误
	ErrSaveBoxPanic = errors.TempErr("SaveBox Err: %s")
)
