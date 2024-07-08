package main

import (
	"context"
	"log"
	"time"

	"github.com/shiimoo/godb/mgo"
)

func main() {
	ctx := context.Background()
	dbMgr, err := mgo.GetMgr(ctx, "default")
	if err != nil {
		log.Fatalln(err)
	}
	dbMgr.Seturl("127.0.0.1", 27017)
	num, err := dbMgr.Connect(10)
	if err != nil {
		log.Fatalln(err)
	}
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("数据库连接成功!", num)

	dbMgr.Start()
	time.AfterFunc(10*time.Second, func() {
		dbMgr.Close()
	})
	for {
	}
}
