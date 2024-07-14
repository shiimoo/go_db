package mlog

import (
	"fmt"
	"time"
)

// ***** 日志消息体 *****

func newLog(lv logLv) *Log {
	msg := new(Log)
	msg.lv = lv
	msg.RefreshCreateTime()
	return msg
}

type Log struct {
	lv         logLv     // 日志级别
	createTime time.Time // 日志创建时间
	labels     []string  // 日志标签列表, 一个日志可能有多个标签, 标签顺序即优先级

	format   string   // 日志模版
	data     []any    // 关键数据
	dataDesc []string // 关键数据的描述
}

// RefreshTime 刷新最新的日志时间
func (l *Log) RefreshCreateTime() {
	l.createTime = time.Now()
}

// Time 获取时间
func (l *Log) Time() time.Time {
	return l.createTime
}

// AddLabel 添加标签
func (l *Log) AddLabel(labels ...string) {
	if l.labels == nil {
		l.labels = make([]string, 0)
	}
	l.labels = append(l.labels, labels...)
}

// SetFormat 设置模版
func (l *Log) SetFormat(format string) {
	l.format = format
}

// AddData 添加数据
func (l *Log) AddData(data any, desc string) {
	if l.data == nil {
		l.data = make([]any, 0)
	}
	if l.dataDesc == nil {
		l.dataDesc = make([]string, 0)
	}

	l.data = append(l.data, data)
	l.dataDesc = append(l.dataDesc, desc)
}

// String 日志输出
func (l *Log) String() string {
	if l.data == nil || len(l.data) == 0 {
		return l.format
	}
	return fmt.Sprintf(l.format, l.data...)
}
