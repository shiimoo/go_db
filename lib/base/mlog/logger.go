package mlog

// ***** 日志记录器 *****

// 构建 日志记录器
func newLogger() *logger {
	obj := new(logger)
	// todo 日志记录器初始化
	return obj
}

type logger struct {
	count uint64 // 日志计数器
	size  uint64 // 当前已记录的日志数据量(字节数量),数据两激增时方便截断

	// 默认输出 io.writer
	// 其他输出 io.writer
}
