package mlog

import (
	"context"
	"log"
	"os"
)

// ***** 日志记录器 *****

// 构建 日志记录器
func newLogger(parent context.Context, key string) *logger {
	obj := new(logger)
	obj.key = key
	obj.ctx, obj.cancel = context.WithCancel(parent)
	obj.outChan = make(chan *Log, 1000)
	return obj
}

type logger struct {
	key    string             // 唯一标识(id)
	ctx    context.Context    // 上下文
	cancel context.CancelFunc // 关闭方法

	// opt // 可选项
	outChan chan *Log
	outFunc func(*Log) error

	isOpen bool // 开启状态
}

// Output 标准输出
func (l *logger) Output(lv logLv, label string, msg string) {
	l.Outputf(lv, label, msg)
}

// Outputf 格式化标准输出
func (l *logger) Outputf(lv logLv, label string, format string, datas ...Data) {
	msg := newLog(lv)
	msg.SetLabel(label)
	msg.SetFormat(format)
	msg.AddData(datas...)

	l.output(msg)
}

// SetOutFunc 设置个性化输出方法
func (l *logger) SetOutFunc(handler func(*Log) error) {
	l.outFunc = handler
}

func (l *logger) output(msg *Log) {
	if !l.isOpen {
		log.Println(l.isOpen, ErrLoggerIsClose)
		return
	}
	l.outChan <- msg
}

// ***** ServerAPI

func (l *logger) Close() {
	l.cancel()
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
			if err := l._output(msg); err != nil {
				log.Println(err)
			}
		}
	}
}

func (l *logger) _output(msg *Log) error {
	var err error
	if l.outFunc == nil {
		err = defaultOutFunc(msg)
	} else {
		err = l.outFunc(msg)
	}
	if msg.lv == LvFatal {
		os.Exit(1)
	}
	return err
}
