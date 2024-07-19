package mlog

import (
	"fmt"
	"time"
)

// ***** 日志消息体 *****

func newLog(lv logLv) *Log {
	msg := new(Log)
	msg.lv = lv
	msg.createTime = time.Now()
	return msg
}

type Data struct {
	Value any
	Key   string
}

type Log struct {
	lv         logLv     // 日志级别
	createTime time.Time // 日志创建时间
	labels     []string  // 日志标签列表, 一个日志可能有多个标签, 标签顺序即优先级

	format string // 日志模版
	datas  []Data // 日志数据
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
func (l *Log) AddData(list ...Data) {
	if l.datas == nil {
		l.datas = make([]Data, 0)
	}

	l.datas = append(l.datas, list...)
}

// String 日志输出
func (l *Log) String() string {
	if l.datas == nil || len(l.datas) == 0 {
		return l.format
	}

	values := make([]any, 0)
	for _, data := range l.datas {
		values = append(values, data.Value)
	}
	return fmt.Sprintf(l.format, values...)
}
