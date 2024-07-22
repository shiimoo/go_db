package mlog

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestLogLv(t *testing.T) {
	log.Println(Debug.String(), 1)
	log.Println(Info.String(), 1)
	log.Println(Warn.String(), 1)
	log.Println(Error.String(), 1)
	log.Println(Fatal.String(), 1)
}

func TestLogMsg(t *testing.T) {

	logger := newLogger(context.Background(), "default")
	// logger.
	logger.Start()
	logger.Output(Info, "test", "log Test")
	logger.Outputf(Info, "test", "log Test %s ", Data{"time", time.Now()})
	time.Sleep(10 * time.Second)
}

func TestLogMgr(t *testing.T) {
	root := context.Background()
	MgrInit(root)

	GetMgr().Info("player", "bus", "info text")
	GetMgr().Infof("player", "bus", "info texts %s %s",
		Data{"time", "tyime"},
		Data{"name", "shimo"},
	)

	time.Sleep(1 * time.Second)
}
