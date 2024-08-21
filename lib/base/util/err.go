package util

import "github.com/shiimoo/godb/lib/base/errors"

var (
	// 分包数据异常
	ErrPackNumError = errors.TempErr("Pack Err: total[%d] num[%d]")
	// 包体数据异常
	ErrPackSizeError = errors.TempErr("Pack size Err: total[%d] has[%d]")
)
