package mlog

import "github.com/shiimoo/godb/lib/base/errors"

var (
	// mongo链接关闭
	ErrLoggerIsClose = errors.TempErr("logger close")
)
