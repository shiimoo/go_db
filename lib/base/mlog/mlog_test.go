package mlog

import (
	"log"
	"testing"
)

func TestLogLv(t *testing.T) {
	log.Println(Debug.String(), 1)
	log.Println(Info.String(), 1)
	log.Println(Warn.String(), 1)
	log.Println(Error.String(), 1)
	log.Println(Fatal.String(), 1)
}

func TestLogMsg(t *testing.T) {
	msg := newLog(Info)
	msg.format = "hello %s"
	msg.AddData("test", "测试数据")
	log.Println(msg)
}
