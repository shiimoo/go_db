package mlog

import "github.com/shiimoo/godb/lib/base/errors"

var (
	// 日志生成器已关闭
	ErrLoggerIsClose = errors.TempErr("logger close")
	// 生成式方法为空
	ErrLoggerOutFuncIsNil = errors.TempErr("logger outfunc is nil")
)
