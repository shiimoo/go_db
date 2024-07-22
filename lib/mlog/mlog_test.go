package mlog

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestLogLv(t *testing.T) {
	log.Println(LvDebug.String(), 1)
	log.Println(LvInfo.String(), 1)
	log.Println(LvWarn.String(), 1)
	log.Println(LvError.String(), 1)
	log.Println(LvFatal.String(), 1)
}

func TestLogMsg(t *testing.T) {

	logger := newLogger(context.Background(), "default")
	// logger.
	logger.Start()
	logger.Output(LvInfo, "test", "log Test")
	logger.Outputf(LvInfo, "test", "log Test %s ", Data{"time", time.Now()})
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
