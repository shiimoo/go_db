package mlog

import (
	"context"
	"fmt"
	"log"
)

// ***** 日志记录器 *****

// 构建 日志记录器
func newLogger(ctx context.Context) *logger {
	obj := new(logger)
	obj.ctx, obj.cancel = context.WithCancel(ctx)
	obj.outChan = make(chan *Log, 1000)
	// todo 日志记录器初始化
	return obj
}

type logger struct {
	ctx    context.Context
	cancel context.CancelFunc
	// opt // 可选项
	outChan chan *Log
	outFunc func(*Log) error

	isOpen bool // 开启状态
}

func (l *logger) Output(lv logLv, labels []string, msg string) {
	l.Outputf(lv, labels, msg)

}

func (l *logger) Outputf(lv logLv, labels []string, format string, datas ...Data) {
	msg := newLog(lv)
	msg.AddLabel(labels...)
	msg.SetFormat(format)
	msg.AddData(datas...)

	l.output(msg)
}

func (l *logger) output(msg *Log) {
	if !l.isOpen {
		log.Println(ErrLoggerIsClose)
		return
	}
	l.outChan <- msg
}

func (l *logger) closeCallBack() {
	l.isOpen = false // 设置关闭状态
	close(l.outChan) // 关闭channel
	for msg := range l.outChan {
		l._output(msg)
	}
}

func (l *logger) Start() {
	l.isOpen = true // 设置关闭状态
	go l._start()
}

func (l *logger) _start() {
	for {
		select {
		case <-l.ctx.Done():
			l.closeCallBack()
			return
		case msg := <-l.outChan:
			l._output(msg)
		}
	}
}

func (l *logger) _output(msg *Log) {
	// opt 输出选项 todo
	if l.outFunc != nil {
		l.outFunc(msg)
	} else {
		fmt.Println(msg)
	}
}
