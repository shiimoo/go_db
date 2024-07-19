package mlog

import (
	"fmt"
	"os"
)

// 默认输出方法
func defaultOutFunc(msg *Log) error {
	// 输出信息：[时间]: 日志等级|日志数据 \n
	outStr := fmt.Sprintf("[%s] %s| %s\n",
		msg.Time().Format("2006-01-02 15:04:05"),
		msg.lv.String(),
		msg.String(),
	)
	os.Stdout.WriteString(outStr)
	return nil
}
