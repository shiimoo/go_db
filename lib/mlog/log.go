package mlog

import (
	"fmt"
	"strings"
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
	Key   string
	Value any
}

func (d *Data) String() string {
	if strings.TrimSpace(d.Key) == "" {
		return fmt.Sprintf("%v", d.Value)
	} else {
		return fmt.Sprintf("%s=%v", d.Key, d.Value)
	}
}

type Log struct {
	lv         logLv     // 日志级别
	createTime time.Time // 日志创建时间
	label      string    // 日志标签

	format string // 日志模版
	datas  []Data // 日志数据
}

// Time 获取时间
func (l *Log) Time() time.Time {
	return l.createTime
}

// SetLabel 设置标签
func (l *Log) SetLabel(label string) {
	l.label = label
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
		values = append(values, data.String())
	}
	return fmt.Sprintf(l.format, values...)
}
